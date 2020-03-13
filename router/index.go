package router

import (
	"../common"
	"../models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Index(c *gin.Context) {
	cid := c.DefaultQuery("cid", "0")
	cookie, _ := c.Cookie("user")
	// 获取用户信息
	userInfo := models.GetUserInfo(cookie)
	lastTime := "暂无"
	if userInfo.Last_time != 0 {
		lastTime = common.DateFormat(userInfo.Last_time)
	}
	condition :=  map[interface{}]interface{}{
		"status" :0 ,
		"cid":cid,
	}
	if condition["cid"] == "0" {
		delete(condition , "cid")
	}
	articles := common.ORM.From(&models.Article{}).Order("create_time").Where(condition).Select("id,title,create_time")
	if len(articles) > 0 {
		for _ , v := range articles{
			if v["CreateTime"] != "0" {
				times , _ := strconv.ParseInt(v["CreateTime"].(string) , 10 , 64)
				v["FormatTime"] = common.DateFormat(times)
			}
		}
	}

	c.HTML(200, "index.html", gin.H{
		"userInfo":  userInfo,
		"last_time": lastTime,
		"navs":      models.GetCategory(),
		"list":      articles,
	})
}