package api

import (
	"encoding/base64"
	"fmt"
	"github.com/9chain/nbcapid/apikey"
	"github.com/9chain/nbcapid/primitives"
	"github.com/9chain/nbcapid/source"
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"sync"
	"time"
)

type ridChan struct {
	done chan interface{}
}

var (
	nextTxID       func() string
	ridChanMap     = make(map[string]*ridChan)
	ridChanMapLock sync.Mutex
)

func init() {
	// 注册标准json2rpc处理函数
	handlers["source-insert-batch"] = sourceInsertBatch
	handlers["source-state"] = sourceState
	handlers["source-transactions"] = sourceTransactions
	handlers["source-transaction"] = sourceTransaction

	nextTxID = primitives.NewGenerator("tx")

	source.SourceEmitter.On("create", func(id string, r *primitives.JSON2Response) {
		ridChanMapLock.Lock()
		defer ridChanMapLock.Unlock()

		chanInfo, ok := ridChanMap[id]
		if ok {
			chanInfo.done <- r
		}
	})
}

/*
{
  "records": [
    {
      "value": "d961824324fb98f7b1b2b8e7278d960f3712df7d8411fd36527a29fb5e64914b",
      "key": "d961824324fb98f7b1b2b8e7278d960f1"
    },
    {
      "value": "d961824324fb98f7b1b2b8e7278d960f3712df7d8411fd36527a29fb5e64914b",
      "key": "d961824324fb98f7b1b2b8e7278d960f2"
    },
    {
      "value": "d961824324fb98f7b1b2b8e7278d960f3712df7d8411fd36527a29fb5e64914b",
      "key": "d961824324fb98f7b1b2b8e7278d960f3"
    },
    {
      "value": "d961824324fb98f7b1b2b8e7278d960f3712df7d8411fd36527a29fb5e64914b",
      "key": "d961824324fb98f7b1b2b8e7278d960f4"
    }
  ],
  "channel": "my_channel"
}
*/
func sourceInsertBatch(ctx *gin.Context, params interface{}) (interface{}, *JSONError) {
	var p source.SourceBatchRecord
	log.Trace(params)
	if err := MapToObject(params, &p); err != nil {
		log.Warnf("invalid param %s", err.Error())
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	ak, _ := ctx.Get("apiKey")
	if err := apikey.CheckChannel(ak.(string), p.Channel); err != nil {
		log.Warnf("invalid channel %+v %+v", ak, p.Channel)
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	total := len(p.Records)
	if total == 0 {
		return gin.H{"message": "success"}, nil
	}

	if total > 1000 {
		log.Warnf("too many records in batch %d", total)
		return nil, primitives.NewCustomInternalError("too many records")
	}

	rid := nextTxID()
	if !source.EnqueueBatch(&p, rid) {
		log.Errorf("enqueue fail")
		return nil, primitives.NewCustomInternalError("server busy")
	}

	// TODO
	done := make(chan interface{})
	ridChanMapLock.Lock()
	ridChanMap[rid] = &ridChan{done: done}
	ridChanMapLock.Unlock()

	var res *primitives.JSON2Response
	select {
	case r := <-done:
		res = r.(*primitives.JSON2Response)
		break
	case <-time.After(time.Second * 10):
		fmt.Println("===== timeout")
		break
	}

	ridChanMapLock.Lock()
	delete(ridChanMap, rid)
	close(done)
	ridChanMapLock.Unlock()

	if nil == res {
		return nil, primitives.NewCustomInternalError("timeout to get txid")
	}

	if res.Error != nil {
		return nil, res.Error
	}

	mm, ok := res.Result.(map[string]interface{})
	if !ok {
		return nil, primitives.NewCustomInternalError("invalid result")
	}

	return gin.H{"message": "success", "tx_id": mm["txId"]}, nil
}

/*
{
  "key": "d961824324fb98f7b1b2b8e7278d960f4",
  "channel": "my_channel"
}

*/
func sourceState(ctx *gin.Context, params interface{}) (interface{}, *JSONError) {
	var p struct {
		Channel string `json:"channel"`
		Key     string `json:"key"`
	}

	if err := MapToObject(params, &p); err != nil {
		log.Warnf("invalid param %+v", params)
		return nil, primitives.NewInvalidParamsError()
	}

	if len(p.Key) == 0 || len(p.Channel) == 0 {
		log.Warnf("invalid param %+v", params)
		return nil, primitives.NewInvalidParamsError()
	}

	ak, _ := ctx.Get("apiKey")
	if err := apikey.CheckChannel(ak.(string), p.Channel); err != nil {
		log.Warnf("invalid channel %+v %+v", ak, p.Channel)
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	res, err := source.QueryState(p.Channel, p.Key)
	if err != nil {
		if _, ok := err.(*primitives.JSONError); ok {
			return nil, err.(*primitives.JSONError)
		}
		log.Errorf("QueryState fail %s", err.Error())
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	var result struct {
		State string `json:"state"`
	}

	if err := MapToObject(res, &result); err != nil {
		log.Errorf("Parse Result fail %s", err.Error())
		log.Flush()
		panic(err)
	}

	bs, err := base64.StdEncoding.DecodeString(result.State)
	result.State = string(bs)

	return result, nil
}

/*
{
  "key": "d961824324fb98f7b1b2b8e7278d960f4",
  "channel": "my_channel"
}
*/
func sourceTransactions(ctx *gin.Context, params interface{}) (interface{}, *JSONError) {
	var p struct {
		Channel string `json:"channel"`
		Key     string `json:"key"`
	}

	if err := MapToObject(params, &p); err != nil {
		log.Warnf("invalid param %+v", params)
		return nil, primitives.NewInvalidParamsError()
	}

	if len(p.Key) == 0 || len(p.Channel) == 0 {
		log.Warnf("invalid param %+v", params)
		return nil, primitives.NewInvalidParamsError()
	}

	ak, _ := ctx.Get("apiKey")
	if err := apikey.CheckChannel(ak.(string), p.Channel); err != nil {
		log.Warnf("invalid channel %+v %+v", ak, p.Channel)
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	res, err := source.QueryTransactions(p.Channel, p.Key)
	if err != nil {
		if _, ok := err.(*primitives.JSONError); ok {
			return nil, err.(*primitives.JSONError)
		}
		log.Errorf("QueryState fail %s", err.Error())
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	var results []struct {
		Timestamp struct {
			Nanos   int32 `json:"nanos"`
			Seconds int32 `json:"seconds"`
		} `json:"timestamp"`
		TxId  string `json:"tx_id"`
		Value string `json:"value"`
	}

	if err := MapToObject(res, &results); err != nil {
		log.Errorf("Parse Result fail %s", err.Error())
		log.Flush()
		panic(err)
	}

	for i, r := range results {
		if s, err := base64.StdEncoding.DecodeString(r.Value); err != nil {
			log.Errorf("Base64 Decode fail %s", err.Error())
			return nil, primitives.NewCustomInternalError(err.Error())
		} else {
			results[i].Value = string(s)
		}
	}
	_ = fmt.Println
	return results, nil
}

/*
[
    {
      "timestamp": {
        "nanos": 4000000,
        "seconds": 1522393304
      },
      "tx_id": "26c952bb482534a2ed9ac337c74a9d7a22d9f121e5123b67fcd36ddf5dfd1edb",
      "value": "ZDk2MTgyNDMyNGZiOThmN2IxYjJiOGU3Mjc4ZDk2MGYzNzEyZGY3ZDg0MTFmZDM2NTI3YTI5ZmI1ZTY0OTE0Yg=="
    }
]
*/

/*
{
  "key": "d961824324fb98f7b1b2b8e7278d960f4",
  "channel": "my_channel"
}
*/
func sourceTransaction(ctx *gin.Context, params interface{}) (interface{}, *JSONError) {
	var p struct {
		Channel string `json:"channel"`
		TxId    string `json:"tx_id"`
	}

	if err := MapToObject(params, &p); err != nil {
		log.Warnf("invalid param %+v", params)
		return nil, primitives.NewInvalidParamsError()
	}

	if len(p.TxId) == 0 || len(p.Channel) == 0 {
		log.Warnf("invalid param %+v", params)
		return nil, primitives.NewInvalidParamsError()
	}

	ak, _ := ctx.Get("apiKey")
	if err := apikey.CheckChannel(ak.(string), p.Channel); err != nil {
		log.Warnf("invalid channel %+v %+v", ak, p.Channel)
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	res, err := source.QueryTransaction(p.Channel, p.TxId)
	if err != nil {
		if _, ok := err.(*primitives.JSONError); ok {
			return nil, err.(*primitives.JSONError)
		}
		log.Errorf("QueryState fail %s", err.Error())
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	mm, ok := res.(map[string]interface{})
	if !ok {
		return nil, primitives.NewCustomInternalError("invalid response")
	}

	mm["channel_id"] = p.Channel

	return res, nil
	//
	//var results []struct {
	//	Timestamp struct {
	//		Nanos   int32 `json:"nanos"`
	//		Seconds int32 `json:"seconds"`
	//	} `json:"timestamp"`
	//	TxId  string `json:"tx_id"`
	//	Value string `json:"value"`
	//}
	//
	//if err := MapToObject(res, &results); err != nil {
	//	log.Errorf("Parse Result fail %s", err.Error())
	//	log.Flush()
	//	panic(err)
	//}
	//
	//for i, r := range results {
	//	if s, err := base64.StdEncoding.DecodeString(r.Value); err != nil {
	//		log.Errorf("Base64 Decode fail %s", err.Error())
	//		return nil, primitives.NewCustomInternalError(err.Error())
	//	} else {
	//		results[i].Value = string(s)
	//	}
	//}
	//_ = fmt.Println
	//return results, nil
}
