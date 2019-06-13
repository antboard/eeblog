package handle

import (
	"log"
	"net/http"

	"github.com/antboard/eeblog/model"
	"github.com/gin-gonic/gin"
)

// Tags 顶部标签
type Tags struct {
	Active bool
	Tag    string
	URL    string
}

// BlogSummary 博文摘要
type BlogSummary struct {
	Title   string
	Summary string
	URL     string
}

// MainPage 首页渲染数据
type MainPage struct {
	Title      string
	Project    string
	Tags       []*Tags
	Bigtitle   string
	Bitsummary string
	Blogs      []*BlogSummary
}

// Index ...
func Index(c *gin.Context) {
	vbs := model.GetResentBlog(0)
	mp := &MainPage{}
	mp.Title = "eeblog"
	mp.Project = "EEBLOG"
	mp.Tags = make([]*Tags, 0, 10)
	tg := new(Tags)
	tg.Active = true
	tg.Tag = "首页"
	tg.URL = "/"
	mp.Tags = append(mp.Tags, tg)
	mp.Bigtitle = "Hello EEBLOG"
	mp.Bitsummary = "这就是一个比较好的电子工程学博客"
	for _, v := range vbs {
		blog := new(BlogSummary)
		blog.Title = v.Title
		blog.Summary = v.Summary
		blog.URL = "/eeblog/" + v.ID
		mp.Blogs = append(mp.Blogs, blog)
	}

	c.HTML(http.StatusOK, "index.tmpl", mp)
	// c.HTML(http.StatusOK, "index.tmpl", gin.H{"title": "测试", "project": "EEBLOG", "tags": []gin.H{gin.H{"active": true, "tag": "tag", "URL": "/"}}})
	// c.JSON(http.StatusOK, vbs)
	// c.String(http.StatusOK, "...")
}

// Login ...
func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", nil)
}

// Plogin ...
func Plogin(c *gin.Context) {
	mp := make(map[string]string)
	err := c.Bind(mp)
	if err != nil {
		log.Println(err)
	}
	// post 不能重定向. 需要让前端重定向
	c.JSON(http.StatusOK, gin.H{"url": "/backend/"})
}
