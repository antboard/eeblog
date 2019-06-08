package model

import (
	"fmt"

	"code.aliyun.com/icecut/ilibs/icelog"
	"github.com/go-xorm/xorm"
)

var (
	//TODO: 这些值后面改成从环境变量中获取
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "example"
	dbName   = "eeblog"
	schema   = "public"
)

// Engine db
var Engine *xorm.Engine

// GetDBEngine 获取...
func GetDBEngine() *xorm.Engine {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbName)
	if password != "" {
		psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	}
	//格式
	engine, err := xorm.NewEngine("postgres", psqlInfo)
	if err != nil {
		icelog.Fatal(err)
		return nil
	}
	// engine.ShowSQL(true) //菜鸟必备
	engine.SetSchema(schema)
	err = engine.Ping()
	if err != nil {
		icelog.Fatal(err)
		return nil
	}
	icelog.Info("connect postgresql success")
	return engine
}
