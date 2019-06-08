package model

import (
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	Engine = GetDBEngine()
	os.Exit(m.Run())
}

func TestAddExt(t *testing.T) {
	_, err := Engine.Query(`create extension "uuid-ossp"`)
	if err != nil {
		// 已经存在也不需要打印
		if strings.Index(err.Error(), "already exists") > 0 {
			return
		}
		t.Error(err)
	} else {
		// 成功不需要打印
	}

}

func TestUUID(t *testing.T) {
	m, err := Engine.Query("select uuid_generate_v4();")
	if err != nil {
		t.Error(err)
		return
	}
	uuid := string(m[0]["uuid_generate_v4"])
	t.Log(uuid) // 无错误此函数不输出
	// t.Error()   // 解除此行注释,上一行log会输出
}
