package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"github.com/9chain/nbcapid/primitives"
)

// 解析json2rpc参数
func parseJSON2Request(ctx *gin.Context) (*JSON2Request, error) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return nil, err
	}

	j, err := primitives.ParseJSON2Request(body)
	if err != nil {
		return nil, err
	}

	return j, nil
}

func MapToObject(source interface{}, dst interface{}) error {
	b, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}

func handleV1Error(ctx *gin.Context, j *JSON2Request, err *JSONError) {
	resp := primitives.NewJSON2Response()
	if j != nil {
		resp.ID = j.ID
	} else {
		resp.ID = nil
	}
	resp.Error = err

	ctx.JSON(200, resp)
}
