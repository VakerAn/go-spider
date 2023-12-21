package server

import (
	"fmt"
	"go-spider/config"
	"go-spider/src/router"
	"net/http"
	"time"
)

var Server *http.Server

func InitServer() {
	readTimeout := config.ConfData.Server.ReadTimeout * time.Second
	writeTimeout := config.ConfData.Server.WriteTimeout * time.Second
	endPoint := fmt.Sprintf(":%d", config.ConfData.Server.HTTPPort)
	maxHeaderBytes := 1 << 20

	Server = &http.Server{
		Addr:           endPoint,
		Handler:        router.Router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}
}
