package model

import (
	"os"
	"strings"
	"testing"

	uuid "github.com/satori/go.uuid"

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

// 插入一条数据
func TestModelUUID(t *testing.T) {
	mt := &Blog{}
	err := Engine.Sync2(mt)
	if err != nil {
		t.Error(err)
	}
	mp, err := Engine.Query(`select * from blog where title='这是一个题目, 最长64字节'`)
	if err != nil {
		t.Error(err)
		return
	}
	if len(mp) > 0 {
		return
	}
	// 插入一个数据
	uuidex, _ := uuid.NewV4()
	mt.ID = uuidex.String()
	mt.Title = "这是一个题目, 最长64字节"
	mt.Summary = "这是一个测试描述, 最长120字节, 这些限制来源于微信订阅号,我考虑能够和微信订阅号同步"
	mt.Text = `这种描述对资深工程师没有什么压力,但是,博客的展现效果就会非常棒. 你可以想象一个模拟工程师的博客会是一个什么样子? 为了写个博客,还有费劲的去示波器截图,做实验,一篇文章成本上千元.写完后,很多人竟然看不懂,过段时间自己也看不懂了.这是个多么悲剧的事情.当然,我做这个博客成本也很高,但是极客,总要玩有意思的事情.
	如果你是一个想转互联网的电工,欢迎你阅读此源代码,我保证了源代码是教科书级的优秀. 如果有些代码难免有坑, 那我就标记出来. 新手必读这篇文章是写给新来者的, 如果你看不懂某些东西,或者没法一次性部署成功, 那请告诉我.这篇文章也是写给开发者的,如果你的代码需要很复杂的过程才能完成开始开发,那我想这一定不是我这个工程的原则.
	我假定没一个开发者可能不太熟悉这里用到的知识,也许你在其他领域是专家. 所以, 每一个技术都有一系列的入门文章. postgresql-极简入门 本项目在用的pgsql之前, 就写了这篇文章. 保证你能直接用起来.也许知识量很少,却能让每一个开发者有成就感.`
	n, err := Engine.Insert(mt)
	if (n != 1) && (err != nil) {
		t.Error(err, n)
	}
}
