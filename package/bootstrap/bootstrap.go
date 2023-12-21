package bootstrap

import (
	"go-spider/config"
	"go-spider/package/pkghttp"
	"go-spider/package/server"
	"go-spider/package/zaplog"
	"go-spider/src/router"
)

func InitBootstrap() {
	config.InitConfig()
	zaplog.InitLogger()
	router.InitRouter()
	server.InitServer()
	pkghttp.InitHttp()
}
