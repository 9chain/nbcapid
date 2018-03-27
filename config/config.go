package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"reflect"
)

type Config struct {
	App struct {
		SeeLogXml string `toml:"seelogXml,omitempty"`
	}
	User struct {
		ConfigDir    string `toml:"configDir,omitempty"`
		MaxUserFiles int    `toml:"maxUserFiles,omitempty"`
	}
	Session struct {
		SessionDir string `toml:"sessionDir,omitempty"`
		MaxAgeMin  int    `toml:"maxAgeMin,omitempty"`
		SessionKey string `toml:"sessionKey,omitempty"`
	}
	SDKSrv struct {
		Username string `toml:"username,omitempty"`
		ApiKey   string `toml:"apiKey,omitempty"`
	} `toml:"sdksrv,omitempty"`
	SMTP struct {
		Host               string `toml:"host,omitempty"`
		ServerAddr         string `toml:"serverAddr,omitempty"`
		User               string `toml:"user,omitempty"`
		Password           string `toml:"password,omitempty"`
		Salt               string `toml:"salt,omitempty"`
		TimeoutMin         int    `toml:"timeoutMin,omitempty"`
		PageForgetPassord  string `toml:"pageForgetPassword,omitempty"`
		ConfirmUrl         string `toml:"confirmUrl,omitempty"`
		ActiveTitle        string `toml:"activeTitle,omitempty"`
		ResetPasswordTitle string `toml:"resetPasswordTitle,omitempty"`
	} `toml:"smtp,omitempty"`
}

var (
	Cfg Config
)

const (
	cfgFileName = "./config.toml"
)

const defaultConfig = `
[App]
seelogXml = "./seelog.xml"

[user]
configDir = "./userconfig"
maxUserFiles = 20

[session]
sessionDir = "/tmp/session_nbcsrv01"
maxageMin = 600
sessionKey = "session key"

[sdksrv]
username = "superuser"
apiKey = "api key"

[smtp]
host = "smtp.163.com"
serverAddr = "smtp.163.com:25"
user = "xxx@xx.com"
password = "xxxxxx"
salt = "salt"
timeoutMin = 120
pageForgetPassword = "/public/resetpassword.html"
confirmUrl = "http://localhost:8080/panel"
activeTitle = "注册确认邮件"
resetPasswordTitle = "重设密码确认邮件"
`

func printCfg(flag string, cfg *Config) {
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(cfg); err != nil {
		panic(err)
	}

	fmt.Printf("==========%s================\n%s\n", flag, buf.String())
}

func LoadConfig() {
	defer printCfg("final", &Cfg)

	var cfg Config
	if _, err := toml.Decode(defaultConfig, &cfg); err != nil {
		panic(err)
	}

	if _, err := os.Stat(cfgFileName); err != nil {
		return
	}

	var newCfg Config
	if _, err := toml.DecodeFile(cfgFileName, &newCfg); err != nil {
		panic(err)
	}

	printCfg(cfgFileName, &newCfg)

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
