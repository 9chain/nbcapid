package apikey

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	log "github.com/cihub/seelog"
	"os"
)

type ApiKey struct {
	Username string
	Channel  string
	ApiKey   string
	SecretKey   string
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
	apiKeyCfg  = make(map[string]ApiKey)
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
		log.Error("loadConfig channel fail")
		log.Flush()
		panic(err)
	}

	var userConfig UserConfig
	if err := loadConfig("cfg/usercfg.toml", &userConfig); err != nil {
		log.Error("loadConfig user fail")
		log.Flush()
		panic(err)
	}

	for _, r := range userConfig.UserKeys {
		apiKeyCfg[r.ApiKey] = r
	}

	fmt.Println("===== apiKeyCfg", apiKeyCfg)
}

func CheckApiKey(apiKey string) bool {
	_, ok := apiKeyCfg[apiKey]
	return ok
}

func CheckChannel(apiKey, channel string) error {
	// channel是否存在
	_, ok := channelCfg.Channels[channel]
	if !ok {
		fmt.Println(apiKey, channel, channelCfg.Channels)
		return errors.New("invalid channel")
	}

	// apiKey是否有权限
	apiKeyInfo, ok := apiKeyCfg[apiKey]
	if !ok {
		return errors.New("invalid apiKey")
	}

	if apiKeyInfo.Channel == channel {
		return nil
	}

	return errors.New("no permission")
}

func MasterChannel(channel string) (string, error) {
	info, ok := channelCfg.Channels[channel]
	if !ok {
		return "", errors.New("invalid channel")
	}

	return info.MasterChannel, nil
}
