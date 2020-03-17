package router

import (
	"../common"
	"../models"
	"github.com/gin-gonic/gin"
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
	articles := models.GetAllArticle(condition)

	c.HTML(200, "index.html", gin.H{
		"userInfo":  userInfo,
		"last_time": lastTime,
		"navs":      models.GetCategory(),
		"list":      articles,
	})
}

func Captcha(c *gin.Context)  {
	common.ServeHTTP(c.Writer , c.Request)
}