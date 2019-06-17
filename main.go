package main

import (
	"github.com/antboard/eeblog/handle"
	"github.com/antboard/eeblog/model"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	model.Engine = model.GetDBEngine()
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/js/", "./web/js/")
	r.Static("/css/", "./web/css/")

	r.GET("/", handle.Index)
	r.GET("/myblog/", handle.Login)
	r.POST("/myblog/", handle.Plogin)
	r.GET("/backend/", handle.Backend)
	r.GET("/new/", handle.NewBlog)

	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
