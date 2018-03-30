package main

import (
	"github.com/9chain/nbcapid/api"
	"github.com/9chain/nbcapid/apikey"
	"github.com/gin-gonic/gin"
	"github.com/9chain/nbcapid/sdkclient"
	"github.com/9chain/nbcapid/config"
	"github.com/9chain/nbcapid/source"
)

func main() {
	config.Init()
	apikey.Init()

	r := gin.Default()
	r.Use(gin.Recovery())

	api.Init(r.Group("api"))

	sdkclient.Init()
	source.Init()

	r.Run() // listen and serve on 0.0.0.0:8080
}
