package api

import (
	"encoding/json"
	"github.com/9chain/nbcapid/primitives"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"fmt"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"
	"math/rand"
)

type SignedData struct {
	Data string		`json:"data"`
	Nonce string	`json:"nonce"`
	Timestamp *string	`json:"timestamp,omitempty"`
	Deadline *string			`json:"deadline,omitempty"`
	Alg *string			`json:"alg,omitempty"`
	Sign string `json:"sign"`
}

func randInt() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(10000000)
}

func signData(body interface{}, secretKey string) *SignedData {
	bs, _ := json.Marshal(body)
	now := time.Now()

	timestamp:=fmt.Sprintf("%d", now.Unix())
	p := SignedData{
		Data: string(bs),
		Nonce: fmt.Sprintf("%d", randInt()),
		Timestamp:&timestamp,
	}

	hh := md5.New()
	hh.Write(bs)
	hh.Write([]byte(p.Nonce))
	hh.Write([]byte(secretKey))
	hh.Write([]byte(*p.Timestamp))
	p.Sign = hex.EncodeToString(hh.Sum(nil))
	return &p
}

func checkSign(body []byte, secretKey string) (*SignedData, error) {
	var p SignedData

	if err := json.Unmarshal(body, &p); err != nil {
		fmt.Println("xxx", err)
		return nil, err
	}

	hh := md5.New()

	if p.Alg != nil {
		hh.Write([]byte(*p.Alg))
	}

	hh.Write([]byte(p.Data))
	if p.Deadline != nil {
		hh.Write([]byte(*p.Deadline))
	}
	hh.Write([]byte(p.Nonce))
	hh.Write([]byte(secretKey))
	if p.Timestamp != nil {
		hh.Write([]byte(*p.Timestamp))
	}

	if p.Alg != nil && (*p.Alg != "md5" && *p.Alg != "") {
		fmt.Println("not support yet")
		return nil, errors.New("not support yet " + *p.Alg)
	}

	hexStr := hex.EncodeToString(hh.Sum(nil))
	if hexStr == p.Sign {
		return &p, nil
	}

	return nil, errors.New("invalid sign")
}

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
