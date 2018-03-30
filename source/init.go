package source

import (
	"fmt"
	"github.com/9chain/nbcapid/config"
	"github.com/9chain/nbcapid/primitives"
	"github.com/9chain/nbcapid/sdkclient"
	log "github.com/cihub/seelog"
	"sync"
	"time"
	"github.com/9chain/nbcapid/apikey"
	"errors"
)

type SourceRecord struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SourceBatchRecord struct {
	Channel string         `json:"channel"`
	Records []SourceRecord `json:"records"`
}

var (
	wsReady       = false
	maxQueueCount = 0
	queueChan     chan map[string]interface{}

	nextJsID func() string
	nextTxID func() string

	flightRpcMap  = make(map[string]map[string]interface{})
	flightTxMap   = make(map[string]map[string]interface{})
	flightMapLock sync.Mutex
)

func EnqueueBatch(batch *SourceBatchRecord) bool {
	sendToChannel := func(userChannel string, records []SourceRecord) {
		// map channel & key
		realChannel, _ := apikey.MasterChannel(batch.Channel)
		for i, _ := range records {
			records[i].Key = fmt.Sprintf("%s_%s", userChannel, records[i].Key)
		}

		jsid, txid := nextJsID(), nextTxID()
		batchRecord := SourceBatchRecord{Channel:realChannel, Records:records}
		p := struct {
			Rid string `json:"rid"`
			SourceBatchRecord
		}{txid, batchRecord}
		req := primitives.NewJSON2Request("create", jsid, p)
		bs, _ := req.JSONByte()

		flightMapLock.Lock()
		defer flightMapLock.Unlock()
		flightRpcMap[jsid] = map[string]interface{}{"type": "create", "rid": txid, "jsid": jsid, "batch_record": &batchRecord, "active": time.Now()}

		queueChan <- map[string]interface{}{"bytes":bs, "jsid":jsid, "type":"create"}
	}

	go func() {
		maxCount := config.Cfg.Source.MaxRecordsPerTx
		records := batch.Records

		var left []SourceRecord
		var slice []SourceRecord
		for {
			if len(records) <= maxCount {
			sendToChannel(batch.Channel, records)
				break
			}

			slice, left = records[:maxCount], records[maxCount:]
			records = left

			sendToChannel(batch.Channel,slice)
		}
	}()
	return true
}

func writeMessageCb(typ, jsid string, err error) {
	switch typ {
	case "create":
		flightMapLock.Lock()
		defer flightMapLock.Unlock()

		if err != nil {
			delete(flightRpcMap, jsid)
			return
		}

		m, _ := flightRpcMap[jsid]
		rid, _ := m["rid"]
		batchRecord, _ := m["batch_record"]

		flightTxMap[rid.(string)] = map[string]interface{}{"type": "create", "rid": rid, "batch_record": batchRecord, "active": time.Now()}

		return
	case "transactions":
		flightMapLock.Lock()
		defer flightMapLock.Unlock()

		if err != nil {
			delete(flightRpcMap, jsid)
			return
		}

		return
	case "state":
		flightMapLock.Lock()
		defer flightMapLock.Unlock()

		if err != nil {
			delete(flightRpcMap, jsid)
			return
		}

		return
	}
}

func QueryTransactions(channel, key string) (interface{}, error) {
	realChannel, _ := apikey.MasterChannel(channel)
	realKey := fmt.Sprintf("%s_%s", channel, key)

	jsid := nextJsID()

	p := map[string]string {"channel": realChannel, "key": realKey}
	req := primitives.NewJSON2Request("transactions", jsid, p)
	bs, _ := req.JSONByte()

	waitChan := make(chan interface{})
	flightMapLock.Lock()
	flightRpcMap[jsid] = map[string]interface{}{"wait_chan": waitChan, "jsid": jsid, "bytes": bs, "type":"state"}
	flightMapLock.Unlock()

	queueChan <- map[string]interface{}{"bytes":bs, "jsid":jsid, "type":"state"}

	res, ok := <-waitChan
	close(waitChan)

	flightMapLock.Lock()
	delete(flightRpcMap, jsid)
	flightMapLock.Unlock()

	if !ok {
		fmt.Println("wait chan error", res, ok)
		return nil, errors.New("query transactions fail")
	}

	j, ok := res.(*primitives.JSON2Response)
	if !ok {
		panic("xxx")
	}
	if j.Error != nil {
		return nil, j.Error
	}

	return j.Result, nil
}

func QueryState(channel, key string) (interface{}, error) {
	realChannel, _ := apikey.MasterChannel(channel)
	realKey := fmt.Sprintf("%s_%s", channel, key)

	jsid := nextJsID()

	p := map[string]string {"channel": realChannel, "key": realKey}
	req := primitives.NewJSON2Request("state", jsid, p)
	bs, _ := req.JSONByte()

	waitChan := make(chan interface{})
	flightMapLock.Lock()
	flightRpcMap[jsid] = map[string]interface{}{"wait_chan": waitChan, "jsid": jsid, "bytes": bs, "type":"state"}
	flightMapLock.Unlock()

	queueChan <- map[string]interface{}{"bytes":bs, "jsid":jsid, "type":"state"}

	res, ok := <-waitChan
	close(waitChan)

	flightMapLock.Lock()
	delete(flightRpcMap, jsid)
	flightMapLock.Unlock()

	if !ok {
		fmt.Println("wait chan error", res, ok)
		return nil, errors.New("query transactions fail")
	}

	j, ok := res.(*primitives.JSON2Response)
	if !ok {
		panic("xxx")
	}
	if j.Error != nil {
		return nil, j.Error
	}

	return j.Result, nil
}

func handleWSResponse(r *primitives.JSON2Response) {
	flightMapLock.Lock()
	defer flightMapLock.Unlock()
	id := r.ID.(string)
	flight, ok := flightRpcMap[id]
	if !ok {
		log.Errorf("missing json id %s", id)
		return
	}

	typ, ok := flight["type"].(string)
	if  !ok {
		panic("logical error")
	}

	delete(flightRpcMap, id)

	switch typ {
	case "transactions":
		waitChan, ok := flight["wait_chan"].(chan interface{})
		if !ok {
			panic("logical error")
		}
		waitChan <- r
		break
	case "create":
		rid, _ := flight["rid"].(string)
		_, ok = flightTxMap[rid]
		if !ok {
			log.Errorf("ERROR why missing rid %s", rid)
			return
		}

		delete(flightRpcMap, rid)
		log.Debugf("create rid %s ok", rid)
		break
	case "state":
		waitChan, ok := flight["wait_chan"].(chan interface{})
		if !ok {
			panic("logical error")
		}
		waitChan <- r
		break
	default:
		panic("no such type " + typ)
	}
}

func Init() {
	maxQueueCount = config.Cfg.Source.MaxQueueCount
	queueChan = make(chan map[string]interface{}, maxQueueCount)

	nextTxID = primitives.NewGenerator("tx")
	nextJsID = primitives.NewGenerator("rpc")

	sdkclient.On("connect", func() {
		fmt.Println("connect ok")
		wsReady = true // sdksrvd已经准备好，开始发送数据
	})

	sdkclient.On("close", func() {
		fmt.Println("close")
		wsReady = false
	})

	sdkclient.On("message", func(message []byte) {
		r, err := primitives.ParseJSON2Response(message)
		if err != nil {
			return
		}

		handleWSResponse(r)
	})

	go func() {
		for {
			if wsReady {
				process()
			}
			time.Sleep(time.Second * 1)
		}
	}()
}

/*
{
	"rid": "txid_727574_1",
	"channel": "master_1",
	"Records": [
	  {
		"key": "k1",
		"value": "hash1"
	  },
	  {
		"key": "k2",
		"value": "hash2"
	  },
	  {
		"key": "k3",
		"value": "hash2"
	  }
	]
},
*/
func process() {
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()

	nextTxID = primitives.NewGenerator("tx")
	nextJsID = primitives.NewGenerator("rpc")

	for wsReady {
		select {
		case im, ok := <-queueChan:
			if !ok {
				panic("queueChan error. Closed?!")
			}

			bs, ok := im["bytes"]
			if !ok {
				panic("not find bytes")
			}

			jsid, ok := im["jsid"]
			if !ok {
				panic("not find jsid")
			}

			typ, ok := im["type"]
			if !ok {
				panic("not find jsid")
			}

			err := sdkclient.WriteMessage(bs.([]byte))
			writeMessageCb(typ.(string), jsid.(string), err)
			fmt.Println("Sfasfdas", err)
			if  err != nil {
				log.Errorf("writeMessage fail %s", err.Error())
				break
			}

			break
		case <-tick.C:
			break
		}
	}
}

// 100.125 8082