package api

import (
	"github.com/9chain/nbcapid/apikey"
	"github.com/9chain/nbcapid/primitives"
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

type JSONError = primitives.JSONError
type JSON2Request = primitives.JSON2Request
type JSON2Response = primitives.JSON2Response

var (
	handlers = make(map[string]func(ctx *gin.Context, params interface{}) (interface{}, *JSONError))
)

func Init(r *gin.RouterGroup) {
	// 标准json2rpc处理
	r.POST("v1", func(ctx *gin.Context) {
		j, err := parseJSON2Request(ctx)
		if err != nil {
			handleV1Error(ctx, nil, primitives.NewInvalidRequestError())
			return
		}

		jsonResp, jsonError := handleV1Request(ctx, j)
		if jsonError != nil {
			handleV1Error(ctx, j, jsonError)
			return
		}

		ctx.JSON(200, jsonResp)
	})
	log.Info("init api ok")
}

func checkApiKey(ctx *gin.Context) bool {
	k := ctx.GetHeader("X-Api-Key")
	if len(k) == 0 {
		return false
	}

	if apikey.CheckApiKey(k) {
		ctx.Set("apiKey", k)
		return true
	}

	return false
}

func handleV1Request(ctx *gin.Context, j *JSON2Request) (*JSON2Response, *JSONError) {
	var resp interface{}
	var jsonError *JSONError

	if !checkApiKey(ctx) {
		log.Debug("invalid apikey")
		return nil, primitives.NewCustomInternalError("invalid apikey")
	}

	//　查找、调用　处理函数
	if f, ok := handlers[j.Method]; ok {
		resp, jsonError = f(ctx, j.Params)
	} else {
		jsonError = primitives.NewMethodNotFoundError()
	}

	if jsonError != nil {
		return nil, jsonError
	}

	jsonResp := primitives.NewJSON2Response()
	jsonResp.ID = j.ID
	jsonResp.Result = resp

	return jsonResp, nil
}
