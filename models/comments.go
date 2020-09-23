package models

import (
	"gobbs/common"
)

type Comments struct{
	Id int				`json:"id"`
	Uid int				`json:"uid"`
	Aid int				`json:"aid"`
	Content string		`json:"content"`
	CreateTime int64	`json:"create_time"`
	Status int			`json:"status"`
}

// 添加一条评论
func AddComments(uid , aid int , content string , status int) bool {
	rows , _ := common.ORM.Insert(&Comments{Uid : uid , Aid : aid , Content : content , Status:status , CreateTime:common.GetUnix()})
	if rows > 0 {
		return true
	}
	return false
}

// 获取指定条件的所有评论
func GetAllComments(commentMaps map[interface{}]interface{}) map[int]map[string]interface{} {
	comments := common.ORM.From(&Comments{}).Where(commentMaps).Order("create_time").Select("content,id,create_time")
	return comments
}