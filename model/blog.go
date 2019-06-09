package model

import "log"

// Blog 存放博客文章
type Blog struct {
	ID      string `xorm:"id uuid"`
	Title   string `xorm:"VARCHAR(64)"`
	Summary string `xorm:"VARCHAR(120)"`
	Text    string `xorm:"text"`
}

// GetResentBlog 获取最近的20条
func GetResentBlog(n int) []*Blog {
	vb := make([]*Blog, 0, 20)
	err := Engine.Cols("id", "title", "summary").Limit(20, n).Find(&vb)
	if err != nil {
		log.Println(err)
	}
	return vb
}
