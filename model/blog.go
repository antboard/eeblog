package model

import "log"

// Blog 存放博客文章
type Blog struct {
	ID      string `xorm:"id uuid"`
	Title   string `xorm:"VARCHAR(64)"`
	Summary string `xorm:"VARCHAR(120)"`
	Text    string `xorm:"text"`
	Status  int64
}

// GetOnlineBlog 获取最近的20条
func GetOnlineBlog(n int) []*Blog {
	vb := make([]*Blog, 0, 20)
	err := Engine.Cols("id", "title", "summary").Where("status=1").Limit(20, n).Find(&vb)
	if err != nil {
		log.Println(err)
	}
	return vb
}

// GetResentBlog 获取最近的20条
func GetResentBlog(n int) []*Blog {
	vb := make([]*Blog, 0, 20)
	err := Engine.Cols("id", "title", "summary", "status").Limit(20, n).Find(&vb)
	if err != nil {
		log.Println(err)
	}
	return vb
}

// SetBlogStatus 修改博文张婷
func SetBlogStatus(id string, n int64) {
	b := &Blog{}
	b.Status = n
	_, err := Engine.Cols("status").Where("id=?", id).Update(b)
	if err != nil {
		log.Println(err)
	}
}
