package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-spider/config"
	"net/http"
)

var Router *gin.Engine

func InitRouter() {
	gin.SetMode(config.ConfData.Server.RunMode)
	r := gin.Default()
	conf := cors.DefaultConfig()
	conf.AllowHeaders = []string{"*"}
	conf.AllowAllOrigins = true
	r.Use(cors.New(conf))

	r.GET("/health_check/up", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})
	v1 := r.Group("api")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "success")
		})
	}
	Router = r
}
