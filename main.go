// gojsmap project main.go
package main

import (
	"github.com/lw000/gocommon/app/gin"
	"github.com/lw000/gocommon/web/gin/middleware"
	"gojsmap/global"
	"gojsmap/routers"

	_ "github.com/icattlecoder/godaemon"
	log "github.com/sirupsen/logrus"
)

func initCommonServer() {
	var err error
	err = global.IpServer.LoadData("data/ip2region.db")
	if err != nil {
		log.Panic(err)
	}

	err = global.SourceMapManager.Parse(global.ProjectConfig.SourceMap)
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	if err := global.LoadGlobalConfig(); err != nil {
		log.Error(err)
		return
	}

	initCommonServer()

	if global.ProjectConfig.Port < 0 {
		global.ProjectConfig.Port = 9096
	}

	app := tygin.NewApplication(global.ProjectConfig.Debug)
	app.SetEnableTLS(global.ProjectConfig.TLS.Enable)
	if app.EnableTLS() {
		app.SetTlsFile(global.ProjectConfig.TLS.CertFile, global.ProjectConfig.TLS.KeyFile)
	}

	err := app.Run(global.ProjectConfig.Port, func(a *tygin.WebApplication) {
		a.Engine().Use(tymiddleware.CorsHandler(nil))

		routers.RegisterService(a.Engine())
	})
	log.Panic(err)
}
