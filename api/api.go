package api

import (
	"github.com/9chain/nbcapid/apikey"
	"github.com/9chain/nbcapid/primitives"
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"fmt"
	"errors"
	"io/ioutil"
	"strconv"
	"time"
)

type JSONError = primitives.JSONError
type JSON2Request = primitives.JSON2Request
type JSON2Response = primitives.JSON2Response

var (
	handlers = make(map[string]func(ctx *gin.Context, params interface{}) (interface{}, *JSONError))
)

func Init(r *gin.RouterGroup) {
	// 标准json2rpc处理
	r.POST("v2", func(ctx *gin.Context) {
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

	// 标准json2rpc处理
	r.POST("v2.1", func(ctx *gin.Context) {
		// 1. check api-key, check sign
		k := ctx.GetHeader("X-Api-Key")
		if len(k) == 0 {
			ctx.AbortWithStatus(400)
			return
		}

		secretKey, ok := apikey.GetSecretKey(k)
		if !ok {
			ctx.AbortWithError(400, errors.New("invalid apiKey"))
			return
		}

		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.AbortWithError(400, errors.New("read body fail"))
			return
		}

		signedData, err := checkSign(body, secretKey)
		if err != nil {
			ctx.AbortWithError(400, err)
			return
		}

		// 2. check dupliated commit TODO
		// 3. check deadline
		if signedData.Deadline != nil {
			deadline, err := strconv.ParseInt(*signedData.Deadline, 10, 64)
			if err != nil {
				ctx.AbortWithError(400, err)
				return
			}
			now := time.Now()
			dead := time.Unix(deadline, 0)
			if now.After(dead) {
				ctx.AbortWithError(400, errors.New("request is timeout"))
				return
			}
		}

		// 3. parse json2rpc params
		j, err := primitives.ParseJSON2Request([]byte(signedData.Data))
		if err != nil {
			ctx.AbortWithError(400, err)
			return
		}

		jsResult, jsError := handleV1Request2(ctx, j)
		jsonResp := primitives.NewJSON2Response()
		if jsError != nil {
			jsonResp.Error = jsError
		} else {
			jsonResp.Result = jsResult
		}
		jsonResp.ID = j.ID

		data := signData(jsonResp, secretKey)

		ctx.JSON(200, data)
	})
	_ = fmt.Println
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


func handleV1Request2(ctx *gin.Context, j *JSON2Request) (interface{}, *JSONError) {
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

	return resp, nil
}
