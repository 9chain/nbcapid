package api

import (
	"github.com/9chain/nbcapid/primitives"
	"github.com/9chain/nbcapid/source"
	"github.com/gin-gonic/gin"
	"github.com/9chain/nbcapid/apikey"
)

func init() {
	// 注册标准json2rpc处理函数
	handlers["source-insert-batch"] = sourceInsertBatch
	handlers["source-state"] = sourceState
	handlers["source-transactions"] = sourceTransactions
}

/*
	{
        "channel":"my_channel",
        "data":[{
            "key":"k1",
            "value":"hash1"
        },{
            "key":"k2",
            "value":"hash2"
        }]
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

	p.Channel, _ = apikey.MasterChannel(p.Channel)

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
chain:chainname,
data:[{key1:"xxx1"},{key2:"xxx1"},{key3:"xxx1"}]
}

*/
func sourceState(ctx *gin.Context, params interface{}) (interface{}, *JSONError) {
	// TODO check channel
	return nil, primitives.NewCustomInternalError("not implement")
}

func sourceTransactions(ctx *gin.Context, params interface{}) (interface{}, *JSONError) {
	// TODO check channel
	return nil, primitives.NewCustomInternalError("not implement")
}
