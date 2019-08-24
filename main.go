package main

import (
	"log"

	"github.com/antboard/eeblog/handle"
	_ "github.com/antboard/eeblog/mdex"
	"github.com/antboard/eeblog/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {
	log.SetFlags(log.Llongfile)
	viper.BindEnv("name")
	viper.BindEnv("password")
	viper.BindEnv("dbhost")
	viper.BindEnv("dbuser")
	viper.BindEnv("dbpassword")
	log.SetFlags(log.Lshortfile)

	model.Engine = model.GetDBEngine()
	err := model.Engine.Sync2(new(model.Blog))
	if err != nil {
		log.Fatal(err)
		return
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	store := cookie.NewStore([]byte("EEBLOG"))
	r.Use(sessions.Sessions("eeblog", store))
	r.LoadHTMLGlob("templates/*")
	r.Static("/js/", "./web/js/")
	r.Static("/css/", "./web/css/")

	r.GET("/", handle.Index)
	r.GET("/myblog/", handle.Login)
	r.POST("/myblog/", handle.Plogin)
	r.GET("/backend/", handle.CheckLogin, handle.Backend)
	r.GET("/new/*id", handle.CheckLogin, handle.NewBlog)
	r.GET("/eeblog/:id", handle.Blog)                      // 查看博客
	r.GET("/draft/:id", handle.CheckLogin, handle.Draft)   // 状态草稿
	r.GET("/online/:id", handle.CheckLogin, handle.Online) // 状态上线
	r.GET("/edit/:id", handle.CheckLogin, handle.Edit)     // 编辑博文
	r.POST("/edit/:id", handle.CheckLogin, handle.PEdit)   // 编辑博文

	r.GET("/dumy", handle.DumyBlog)

	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
