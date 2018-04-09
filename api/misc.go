package api

import (
	"encoding/json"
	"github.com/9chain/nbcapid/primitives"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"fmt"
	"crypto/md5"
	"encoding/hex"
)

func checkSign(body []byte, secretKey string) {
	var p struct {
		Data string		`json:"data"`
		Nonce string	`json:"nonce"`
		Timestamp string	`json:"timestamp"`
		Alg string			`json:"alg"`
		Sign string `json:"sign"`
	}

	if err := json.Unmarshal(body, &p); err != nil {
		fmt.Println("xxx", err)
		return
	}

	s := p.Alg + p.Data + p.Nonce + p.Timestamp + secretKey

	if len(p.Alg) > 0 && p.Alg != "md5" {
		fmt.Println("not support yet")
		return
	}

	h := md5.New()
	h.Write([]byte([]byte(s))) // 需要加密的字符串为 123456
	sum := h.Sum(nil)
	hexStr := hex.EncodeToString(sum)
	fmt.Println(111, s)
	fmt.Println(222, hexStr, p.Sign, hexStr == p.Sign)



	fmt.Printf("yyy === %+v\n", p)
}

// 解析json2rpc参数
func parseJSON2Request(ctx *gin.Context) (*JSON2Request, error) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return nil, err
	}

	checkSign(body, "xx")

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
	//fmt.Println(string(b))
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
