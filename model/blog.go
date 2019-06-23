package model

import (
	"log"
	"time"
)

// Blog 存放博客文章
type Blog struct {
	ID      string `xorm:"id uuid pk"`
	Title   string `xorm:"VARCHAR(64)"`
	Summary string `xorm:"VARCHAR(120)"`
	Text    string `xorm:"text"`
	Status  int64
	Created time.Time `xorm:"Created"`
}

// GetOnlineBlog 获取最近的20条
func GetOnlineBlog(n int) []*Blog {
	vb := make([]*Blog, 0, 20)
	err := Engine.Cols("id", "title", "summary").Where("status=1").Desc("created").Limit(20, n).Find(&vb)
	if err != nil {
		log.Println(err)
	}
	return vb
}

// GetResentBlog 获取最近的20条
func GetResentBlog(n int) []*Blog {
	vb := make([]*Blog, 0, 20)
	err := Engine.Cols("id", "title", "summary", "status").Desc("created").Limit(20, n).Find(&vb)
	if err != nil {
		log.Println(err)
	}
	return vb
}

// SetBlogStatus 修改博文状态
func SetBlogStatus(id string, n int64) {
	b := &Blog{}
	b.Status = n
	_, err := Engine.Cols("status").Where("id=?", id).Update(b)
	if err != nil {
		log.Println(err)
	}
}

// GetBlog 获取文章所有信息
func GetBlog(id string) *Blog {
	b := &Blog{}
	_, err := Engine.Where("id=?", id).Get(b)
	if err != nil {
		log.Println(err)
	}
	return b
}

// UpdateBlog 更新博客内容
func UpdateBlog(id, title, summary, ctx string) error {
	b := &Blog{}
	b.Title = title
	b.Summary = summary
	b.Text = ctx
	i, err := Engine.Where("id=?", id).Cols("text", "summary", "title").Update(b)
	if err != nil {
		log.Println(err)
		return err
	}
	if i == 0 {
		b.ID = id
		b.Status = 1
		b.Summary = summary
		b.Title = title
		_, err := Engine.Insert(b)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
