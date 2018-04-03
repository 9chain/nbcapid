package sdkclient

import (
	"errors"
	"github.com/9chain/nbcapid/config"
	"github.com/chuckpreslar/emission"
	log "github.com/cihub/seelog"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var (
	wsClient  *websocket.Conn
	wsEmitter = emission.NewEmitter()
)

func On(event string, listener interface{}) {
	wsEmitter.On(event, listener)
}

func Off(event string, listener interface{}) {
	wsEmitter.Off(event, listener)
}

type SDKParams struct {
	ID     string
	Method string
	Params struct {
		Chain  string
		Params interface{}
	}
}

// metohd, chain, params
func WriteMessage(bs []byte) error {
	if wsClient == nil {
		return errors.New("disconnected!")
	}

	err := wsClient.WriteMessage(websocket.TextMessage, bs)
	if err != nil {
		wsClient.Close()
		return err
	}

	return nil
}

func Init() {
	go func() {
		for {
			start()
			time.Sleep(time.Second * time.Duration(config.Cfg.SDKClient.RetryConnectSecond))
		}
	}()
	log.Info("init sdk client ok")
}

func start() {
	sdkCfg := config.Cfg.SDKClient
	wsUrl := sdkCfg.WSUrl
	c, _, err := websocket.DefaultDialer.Dial(wsUrl, http.Header{
		"X-Api-Key": []string{sdkCfg.ApiKey},
	})

	if err != nil {
		log.Errorf("connect to %s fail. %s\n", wsUrl, err.Error())
		return
	}

	log.Infof("connect %s ok\n", wsUrl)

	wsClient = c
	wsEmitter.Emit("connect")

	defer func() {
		wsClient.Close()
		wsClient = nil
		wsEmitter.Emit("close")
		log.Error("close sdk client")
	}()

	for {
		_, message, err := wsClient.ReadMessage()
		if err != nil { // 出错，停止
			log.Errorf("ReadMessage fail %s", err.Error())
			return
		}

		wsEmitter.Emit("message", message)
	}
}
