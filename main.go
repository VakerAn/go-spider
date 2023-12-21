package main

import (
	"context"
	"go-spider/config"
	"go-spider/package/bootstrap"
	"go-spider/package/server"
	"go-spider/package/shutdown"
	"log"
	"time"
)

func main() {
	bootstrap.InitBootstrap()
	srv := server.Server
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("startup service failed, err:%v\n\n" + err.Error())
		}
	}()
	shutdown.NewHook().Close(
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), config.ConfData.Server.ServerTimeout*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				log.Fatal("server shutdown err", err)
			}
		},
	)
}
