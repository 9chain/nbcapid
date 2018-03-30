package source

import (
	"fmt"
	"github.com/9chain/nbcapid/config"
	"github.com/9chain/nbcapid/primitives"
	"github.com/9chain/nbcapid/sdkclient"
	log "github.com/cihub/seelog"
	"time"
	"sync"
)

type SourceRecord struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SourceBatchRecord struct {
	Channel string `json:"channel"`
	Records []SourceRecord
}

const (
	maxTransactionRecods = 3
)

var (
	wsReady       = false
	maxQueueCount = 0
	queueChan     chan *SourceBatchRecord
	flightMap     = make(map[string]interface{})
	flightMapLock sync.Mutex
)

func EnqueueBatch(batch *SourceBatchRecord) bool {
	go func() {
		maxCount := maxTransactionRecods
		records := batch.Records
		var left []SourceRecord
		var slice []SourceRecord
		for {
			if len(records) <= maxCount {
				queueChan <- &SourceBatchRecord{Channel: batch.Channel, Records: records}
				break
			}
			slice, left = records[:maxCount], records[maxCount:]
			records = left
			queueChan <- &SourceBatchRecord{Channel: batch.Channel, Records: slice}
		}
	}()
	return true
}

func Query(params interface{}) {
	log.Error(1)
}

func Init() {
	maxQueueCount = config.Cfg.Source.MaxQueueCount
	queueChan = make(chan *SourceBatchRecord, maxQueueCount)

	sdkclient.On("connect", func() {
		fmt.Println("connect ok")
		wsReady = true // sdksrvd已经准备好，开始发送数据
	})

	sdkclient.On("close", func() {
		fmt.Println("close")
		wsReady = false
	})

	sdkclient.On("message", func(message []byte) {
		j, err := primitives.ParseJSON2Request(message)
		if err != nil {
			return
		}

		fmt.Println("message", j)
	})

	//go func() {
	//	for i:=1; i < 5; i++{
	//		time.Sleep(time.Second)
	//		Enqueue(11111)
	//	}
	//	close(queueChan)
	//}()

	go func() {
		for {
			if wsReady {
				process()
			}
			time.Sleep(time.Second * 5)
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

	nextID := primitives.NewGenerator("rpc")
	nextTxID := primitives.NewGenerator("tx")
	for wsReady {
		select {
		case e, ok := <-queueChan:
			if !ok {
				panic("queueChan error. Closed?!")
			}

			jsid, txid := nextID(), nextTxID()
			p := struct {
				Rid string `json:"rid"`
				SourceBatchRecord
			}{txid, *e}
			req := primitives.NewJSON2Request("create", jsid, p)
			bs, _ := req.JSONByte()
			if err := sdkclient.WriteMessage(bs); err != nil {
				log.Errorf("writeMessage fail %s", err.Error())
				break
			}

			flightMapLock.Lock()
			flightMap[jsid] = p
			flightMapLock.Unlock()

			break
		case <-tick.C:
			break
		}
	}
}
