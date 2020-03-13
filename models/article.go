package models

import (
	"../common"
)

type Article struct {
	Id int64			`json:"id"`
	Title string		`json:"title"`
	Content string		`json:"content"`
	Status int			`json:"status"`
	CreateTime int64	`json:"create_time"`
	Cid int64			`json:"cid"`
	Uid int64 			`json:"uid"`
	FormatTime string
}

// 添加一篇文章
func AddArticle(uid int64 , title string , content string , cid int64) int64 {
	lastId , _ := common.ORM.Insert(&Article{Title:title , Content:content , Status:0 , CreateTime:common.GetUnix() , Cid:cid , Uid:uid})
	if lastId > 0 {
		return lastId
	}else {
		return 0
	}
}