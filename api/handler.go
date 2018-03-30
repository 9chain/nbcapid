package api

import (
	"github.com/9chain/nbcapid/apikey"
	"github.com/9chain/nbcapid/primitives"
	"github.com/9chain/nbcapid/source"
	"github.com/gin-gonic/gin"
	"fmt"
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

	return res, nil
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
	_ = fmt.Println
	return res, nil
}
