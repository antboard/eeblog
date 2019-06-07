package main

import (
	"github.com/antboard/eeblog/handle"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Static("/", "./web")
	r.GET("/ping", handle.Text)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
