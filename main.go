package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.Default()
	r.Static("/", "./web")

	// r.GET("/ping", handle.Text)

	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
