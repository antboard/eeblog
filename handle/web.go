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
	// c.HTML(http.StatusOK, "backend.tmpl", gin.H{"Title": "测试", "Project": "EEBLOG", "Tags": []gin.H{gin.H{"Active": true, "Tag": "tag", "URL": "/"}}})
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

// Backend 后端列表
func Backend(c *gin.Context) {
	// 读取文章列表, 显示响应入口
	// 文章名, 修改, 草稿, 发布
	c.HTML(http.StatusOK, "backend.tmpl", gin.H{"Title": "测试", "Project": "EEBLOG", "Tags": []gin.H{gin.H{"Active": true, "Tag": "tag", "URL": "/"}}})
}

// NewBlog 创建新文章
func NewBlog(c *gin.Context) {
	// 考虑判断参数,如果有uuid,那么就用这个id去索引文章进入编辑页面,如果索引不到,就认为是新blog
	// 如果没有, 就生成一个uuid,进入新文章页面
	c.HTML(http.StatusOK, "newblog.tmpl", gin.H{"Title": "测试", "Project": "EEBLOG", "Tags": []gin.H{gin.H{"Active": true, "Tag": "tag", "URL": "/"}}})
}
