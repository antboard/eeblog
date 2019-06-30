package handle

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/spf13/viper"

	"github.com/antboard/eeblog/mdex"
	"github.com/antboard/eeblog/model"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// Tags 顶部标签
type Tags struct {
	Active bool
	Tag    string
	URL    string
}

// BlogSummary 博文摘要
type BlogSummary struct {
	Title     string
	Summary   string
	URL       string
	Online    bool   // 上线or草稿
	URLEdit   string // 编辑页面ß
	URLDraft  string // 草稿状态, 点此变成上线
	URLOnline string // 上线状态, 点此变成草稿
}

// MainPage 首页渲染数据
type MainPage struct {
	Title   string         // 网页标题
	Project string         // 网站名称
	Tags    []*Tags        // 导航标签
	Blogs   []*BlogSummary // 文集列表

	BigTitle   string // 首页大标题
	BigSummary string // 首页大摘要

	BlogTitle   string        // 博文题目
	BlogSummary string        // 博文摘要
	BlogCtx     template.HTML // 博文内容
}

var store *sessions.CookieStore

func init() {
	// cookie加密秘钥
	store = sessions.NewCookieStore([]byte("eeblog-antboard-secret"))
}

// Index ...
func Index(c *gin.Context) {
	vbs := model.GetOnlineBlog(0)
	mp := &MainPage{}
	mp.Title = "eeblog"
	mp.Project = "EEBLOG"
	mp.Tags = make([]*Tags, 0, 10)
	tg := new(Tags)
	tg.Active = true
	tg.Tag = "首页"
	tg.URL = "/"
	mp.Tags = append(mp.Tags, tg)
	mp.BigTitle = "Hello EEBLOG"
	mp.BigSummary = "这就是一个比较好的电子工程学博客"
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
	session, err := store.Get(c.Request, "user-sess")
	if err == nil {
		user, ok := session.Values["user"]
		log.Println(user, ok)
		name := viper.Get("name").(string)
		if ok && (user == name) {
			c.Redirect(http.StatusFound, "/backend/")
			return
		}
	}

	c.HTML(http.StatusOK, "login.tmpl", nil)
}

// Plogin ...
func Plogin(c *gin.Context) {
	mp := make(map[string]string)
	err := c.BindJSON(&mp)
	if err != nil {
		log.Println(err)
	}

	name := mp["name"]
	pwd := mp["password"]
	if (name == viper.Get("name")) &&
		(pwd == viper.Get("password")) {
		//   设置cookies, 转型正常后台
		session, err := store.Get(c.Request, "user-sess")
		if err == nil {
			session.Values["user"] = name
			session.Save(c.Request, c.Writer)
			// post 不能重定向. 需要让前端重定向
			c.JSON(http.StatusOK, gin.H{"url": "/backend/"})
		}

	}
	c.String(http.StatusOK, "")
}

// Backend 后端列表
func Backend(c *gin.Context) {
	// 读取文章列表, 显示响应入口
	// 文章名, 修改, 草稿, 发布
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
	for _, v := range vbs {
		blog := new(BlogSummary)
		blog.Title = v.Title
		blog.Summary = v.Summary
		blog.URL = "/eeblog/" + v.ID
		blog.Online = (v.Status == 1)
		blog.URLOnline = "/draft/" + v.ID
		blog.URLDraft = "/online/" + v.ID
		blog.URLEdit = "/edit/" + v.ID
		mp.Blogs = append(mp.Blogs, blog)
	}

	c.HTML(http.StatusOK, "backend.tmpl", mp)
	// c.HTML(http.StatusOK, "backend.tmpl", gin.H{"Title": "测试", "Project": "EEBLOG", "Tags": []gin.H{gin.H{"Active": true, "Tag": "tag", "URL": "/"}}})
}

// NewBlog 创建新文章
func NewBlog(c *gin.Context) {
	// 考虑判断参数,如果有uuid,那么就用这个id去索引文章进入编辑页面,如果索引不到,就认为是新blog
	id := c.Param("id")
	id = strings.Trim(id, "/")
	if id == "" {
		uuidex, _ := uuid.NewV4()
		URL := `/edit/` + uuidex.String()
		c.Redirect(http.StatusFound, URL)
	}
	c.String(http.StatusOK, "what's wrong??")
}

// Blog 博文阅读页
func Blog(c *gin.Context) {
	id := c.Param("id")
	blog := model.GetBlog(id)
	mp := &MainPage{}
	mp.Title = "eeblog"
	mp.Project = "EEBLOG"
	mp.Tags = make([]*Tags, 0, 10)
	tg := new(Tags)
	tg.Active = true
	tg.Tag = "首页"
	tg.URL = "/"
	mp.Tags = append(mp.Tags, tg)
	mp.BlogTitle = blog.Title
	mp.BlogSummary = blog.Summary
	mp.BlogCtx = template.HTML(blog.Text)
	c.HTML(http.StatusOK, "blog.tmpl", mp)
}

// Draft 设置为草稿
func Draft(c *gin.Context) {
	id := c.Param("id")
	model.SetBlogStatus(id, 0)
	c.Redirect(http.StatusFound, "/backend/")
}

// Online 设置为上线
func Online(c *gin.Context) {
	id := c.Param("id")
	model.SetBlogStatus(id, 1)
	c.Redirect(http.StatusFound, "/backend/")
}

// Edit 编辑博文
func Edit(c *gin.Context) {
	id := c.Param("id")
	blog := model.GetBlog(id)
	mp := &MainPage{}
	mp.Title = "eeblog"
	mp.Project = "EEBLOG"
	mp.Tags = make([]*Tags, 0, 10)
	tg := new(Tags)
	tg.Active = true
	tg.Tag = "首页"
	tg.URL = "/"
	mp.Tags = append(mp.Tags, tg)
	mp.BigTitle = blog.Title
	mp.BigSummary = blog.Summary
	mp.BlogCtx = template.HTML(blog.Text)
	c.HTML(http.StatusOK, "edit.tmpl", mp)
}

// PEdit 提交编辑的文章
func PEdit(c *gin.Context) {
	id := c.Param("id")
	mp := make(map[string]string)
	err := c.BindJSON(&mp)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "")
		return
	}
	// log.Println(mp)
	err = model.UpdateBlog(id, mp["title"], mp["summary"], mp["ctx"])
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "")
		return
	}
	// post 不能重定向. 需要让前端重定向
	c.JSON(http.StatusOK, gin.H{"url": "/eeblog/" + id})
}

// DumyBlog 测试markdown接口
func DumyBlog(c *gin.Context) {
	src := `$(80,50)
	U10-P8-NSTC12[1:VCC,8:GND(ddd)](3,2,8)
	U11-P4-NEEPROM[1:VCC,4:GND](10,12,5)
	$`
	var buf bytes.Buffer
	if err := mdex.MD.Convert([]byte(src), &buf); err != nil {
		panic(err)
	}
	// log.Println(buf.String())

	mp := &MainPage{}
	mp.Title = "eeblog"
	mp.Project = "EEBLOG"
	mp.Tags = make([]*Tags, 0, 10)
	tg := new(Tags)
	tg.Active = true
	tg.Tag = "首页"
	tg.URL = "/"
	mp.Tags = append(mp.Tags, tg)
	mp.BlogTitle = "blog.Title"
	mp.BlogSummary = "blog.Summary"
	mp.BlogCtx = template.HTML(buf.String())
	c.HTML(http.StatusOK, "blog.tmpl", mp)
}
