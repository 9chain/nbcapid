package api

import (
	"github.com/gin-gonic/gin"
	"github.com/9chain/nbcapid/primitives"
)

func init() {
	// 注册标准json2rpc处理函数
	handlers["create-chain"] = handleCreateChain
	handlers["create-entry"] = handleCreateEntry
	handlers["entry-status"] = handleStatus
	handlers["entry"] = handleEntry
	handlers["validate"] = handleValidate
}

func handleCreateChain(ctx *gin.Context, params interface{}) (interface{}, *JSONError) {
	// TODO check channel
	return nil, primitives.NewCustomInternalError("not implement")
}

func handleCreateEntry(ctx *gin.Context, params interface{}) (interface{}, *JSONError) {
	// TODO check channel
	return nil, primitives.NewCustomInternalError("not implement")
}

func handleStatus(ctx *gin.Context, params interface{}) (interface{}, *JSONError) {
	// TODO check channel
	return nil, primitives.NewCustomInternalError("not implement")
}

func handleEntry(ctx *gin.Context, params interface{}) (interface{}, *JSONError) {
	// TODO check channel
	return nil, primitives.NewCustomInternalError("not implement")
}

func handleValidate(ctx *gin.Context, params interface{}) (interface{}, *JSONError) {
	// TODO check channel
	return nil, primitives.NewCustomInternalError("not implement")
}
