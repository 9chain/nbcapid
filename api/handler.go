package api

import (
	"encoding/base64"
	"fmt"
	"github.com/9chain/nbcapid/apikey"
	"github.com/9chain/nbcapid/primitives"
	"github.com/9chain/nbcapid/source"
	"github.com/gin-gonic/gin"
)

func init() {
	// 注册标准json2rpc处理函数
	handlers["source-insert-batch"] = sourceInsertBatch
	handlers["source-state"] = sourceState
	handlers["source-transactions"] = sourceTransactions
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
	if err := MapToObject(params, &p); err != nil {
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	ak, _ := ctx.Get("apiKey")
	if err := apikey.CheckChannel(ak.(string), p.Channel); err != nil {
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	total := len(p.Records)
	if total == 0 {
		return gin.H{"message": "success"}, nil
	}

	if total > 1000 {
		return nil, primitives.NewCustomInternalError("too many records")
	}

	if !source.EnqueueBatch(&p) {
		return nil, primitives.NewCustomInternalError("server busy")
	}

	return gin.H{"message": "success"}, nil
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
		return nil, primitives.NewInvalidParamsError()
	}

	if len(p.Key) == 0 || len(p.Channel) == 0 {
		return nil, primitives.NewInvalidParamsError()
	}

	ak, _ := ctx.Get("apiKey")
	if err := apikey.CheckChannel(ak.(string), p.Channel); err != nil {
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	res, err := source.QueryState(p.Channel, p.Key)
	if err != nil {
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	var result struct {
		State string `json:"state"`
	}

	if err := MapToObject(res, &result); err != nil {
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
		return nil, primitives.NewInvalidParamsError()
	}

	if len(p.Key) == 0 || len(p.Channel) == 0 {
		return nil, primitives.NewInvalidParamsError()
	}

	ak, _ := ctx.Get("apiKey")
	if err := apikey.CheckChannel(ak.(string), p.Channel); err != nil {
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	res, err := source.QueryTransactions(p.Channel, p.Key)
	if err != nil {
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
		panic(err)
	}

	for i, r := range results {
		if s, err := base64.StdEncoding.DecodeString(r.Value); err != nil {
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
