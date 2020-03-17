package models

import (
	"../common"
	"strconv"
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

// 获取所有文章
func GetAllArticle(condition map[interface{}]interface{})map[int]map[string]interface{}{
	articles := common.ORM.From(&Article{}).Order("create_time").Where(condition).Select("id,title,create_time")
	if len(articles) > 0 {
		for _ , v := range articles{
			if v["CreateTime"] != "0" {
				times , _ := strconv.ParseInt(v["CreateTime"].(string) , 10 , 64)
				v["FormatTime"] = common.DateFormat(times)
			}
		}
	}
	return articles
}