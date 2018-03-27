package apikey

import (
	"github.com/BurntSushi/toml"
	"os"
)

type ApiKey struct {
	Username string
	ApiKey string
	Channel string
}

type UserConfig struct {
	UserKeys map[string]ApiKey
}

type ChannelInfo struct {
	MasterChannel string
}

type ChannelConfig struct {
	Channels map[string]ChannelInfo
}

var (
	apiKeyCfg = make(map[string]ApiKey)
	channelCfg ChannelConfig
)

func loadConfig(path string, cfg interface{}) error {
	if _, err := os.Stat(path); err != nil {
		return nil
	}

	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return err
	}

	return nil
}

func Init() {
	if err := loadConfig("cfg/channels.toml", &channelCfg); err != nil {
		panic(err)
	}

	var userConfig UserConfig
	if err := loadConfig("cfg/usercfg.toml", &userConfig); err != nil {
		panic(err)
	}

	for _, r := range userConfig.UserKeys {
		apiKeyCfg[r.ApiKey] = r
	}
}

func CheckApiKey(apiKey string) bool {
	_, ok := apiKeyCfg[apiKey]
	return ok
}