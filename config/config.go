package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	log "github.com/cihub/seelog"
	"os"
	"path"
	"reflect"
)

type Config struct {
	SDKClient struct {
		WSUrl              string
		ApiKey             string
		RetryConnectSecond int
	}
	Source struct {
		MaxQueueCount   int
		MaxRecordsPerTx int
	}
}

var (
	Cfg    Config
	cfgDir = "./cfg"
)

const defaultConfig = `
[SDKClient]
WSUrl = "ws://localhost:8082/v1/ws"
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
	dir := os.Getenv("NBCAPID_CFG_DIR")
	if len(dir) > 0 {
		cfgDir = dir
	}

	fmt.Println("cfgDir is", cfgDir)

	if _, err := os.Stat(cfgDir); err != nil {
		panic("not found cfg dir: " + cfgDir)
	}

	initSeelog()

	log.Infof("cfg dir is %s", cfgDir)

	defer printCfg("final", &Cfg)

	var cfg Config
	if _, err := toml.Decode(defaultConfig, &cfg); err != nil {
		panic(err)
	}
	//printCfg("default", &cfg)

	cfgFilePath := path.Join(cfgDir, "nbcapid.toml")
	if _, err := os.Stat(cfgFilePath); err != nil {
		return
	}

	var newCfg Config
	if _, err := toml.DecodeFile(cfgFilePath, &newCfg); err != nil {
		panic(err)
	}

	printCfg(cfgFilePath, &newCfg)

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
