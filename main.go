package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"os"
	"github.com/9chain/nbcapid/api"
	"github.com/9chain/nbcapid/apikey"
)

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

	fmt.Println("use default seelog config")

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

func main() {
	if _, err := os.Stat("cfg"); err != nil {
		panic("not found dir cfg")
	}

	initSeelog()
	apikey.Init()

	r := gin.Default()
	r.Use(gin.Recovery())

	api.Init(r.Group("api"))

	r.Run() // listen and serve on 0.0.0.0:8080
}
