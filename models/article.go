package models

import (
	"../common"
)

type Article struct {
	Id int64
	Title string
	Content string
	Status int
	CreateTime int64
	Cid int64
	FormatTime string
}

// 添加一篇文章
func AddArticle(uid int64 , title string , content string , cid int64) int64 {
	ret , _ := common.MysqlDb.Exec("insert into article(uid , title,content,create_time,status,cid) values (?,?,?,?,?,?)" , uid , title , content , common.GetUnix() , 0 , cid)
	rowsAffected , _ := ret.RowsAffected()
	lastInsertID,_ := ret.LastInsertId()
	if rowsAffected > 0 {
		return lastInsertID
	}else {
		return 0
	}
}