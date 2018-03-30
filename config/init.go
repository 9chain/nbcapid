package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"reflect"
	log "github.com/cihub/seelog"
)

type Config struct {
	SDKClient struct {
		WSUrl string
		ApiKey string
		RetryConnectSecond int
	}
	Source struct {
		MaxQueueCount int
		MaxRecordsPerTx int
	}
}

var (
	Cfg Config
)

const (
	cfgFileName = "./cfg/nbcapid.toml"
)

const defaultConfig = `
[SDKClient]
WSUrl = "ws://localhost:8080/v1/ws"
ApiKey = "test api key"
RetryConnectSecond = 3
[Source]
MaxQueueCount = 1000
MaxRecordsPerTx = 100
`

func printCfg(flag string, cfg *Config) {
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(cfg); err != nil {
		panic(err)
	}

	log.Infof("\n==========%s================\n%s\n", flag, buf.String())
}


func mustCfg() {
	if _, err := os.Stat("cfg"); err != nil {
		panic("not found cfg dir")
	}
}


func initSeelog() {
	cfgPath := "cfg/seelog.xml"
	if _, err := os.Stat(cfgPath); err == nil {
		logger, err := log.LoggerFromConfigAsFile(cfgPath)
		if err != nil {
			panic(err)
		}

		log.ReplaceLogger(logger)
		return
	}

	fmt.Println("---- use default seelog config")

	defaultConfig := `
<seelog>
    <outputs formatid="main">
        <console />
    </outputs>
    <formats>
        <format id="main" format="%l %Date %Time %File:%Line %Msg%n"/>
    </formats>
</seelog>`

	logger, err := log.LoggerFromConfigAsString(defaultConfig)
	if err != nil {
		panic(err)
	}

	log.ReplaceLogger(logger)
}

func Init() {
	mustCfg()
	initSeelog()

	defer printCfg("final", &Cfg)

	var cfg Config
	if _, err := toml.Decode(defaultConfig, &cfg); err != nil {
		panic(err)
	}
	//printCfg("default", &cfg)

	if _, err := os.Stat(cfgFileName); err != nil {
		return
	}

	var newCfg Config
	if _, err := toml.DecodeFile(cfgFileName, &newCfg); err != nil {
		panic(err)
	}

	printCfg(cfgFileName, &newCfg)
	log.Warn("00000000000000000000")
	o, n := toMap(cfg), toMap(newCfg)
	walk(o, n)

	bs, _ := json.Marshal(o)
	if err := json.Unmarshal(bs, &Cfg); err != nil {
		panic(err)
	}
}

func toMap(obj interface{}) map[string]interface{} {
	bs, _ := json.Marshal(obj)
	var res map[string]interface{}
	json.Unmarshal(bs, &res)
	return res
}

func walk(o map[string]interface{}, n map[string]interface{}) {
	for k, v := range o {
		if "map[string]interface {}" == reflect.TypeOf(v).String() {
			walk(v.(map[string]interface{}), n[k].(map[string]interface{}))
			continue
		}

		nv, ok := n[k]
		if !ok {
			continue
		}

		switch nv.(type) {
		case string:
			if nv != "" {
				fmt.Println("reset", k, v, nv)
				o[k] = nv
			}

			break
		case float64:
			if int(nv.(float64)) != 0 {
				fmt.Println("reset", k, v, nv)
				o[k] = nv
			}
			break
		default:
			msg := fmt.Sprintf("not support type: %s yet!!!!!!", reflect.TypeOf(v).String())
			panic(msg)
			break
		}
	}
}
