package router

import (
	"gobbs/common"
	"gobbs/models"

	"github.com/gin-gonic/gin"
)

// 首页
func Index(c *gin.Context) {
	cid := c.DefaultQuery("cid", "0")
	cookie := common.GetUid(c)
	// 获取用户信息
	userInfo := models.GetUserInfo(cookie)
	lastTime := "暂无"
	if userInfo.Last_time != 0 {
		lastTime = common.DateFormat(userInfo.Last_time)
	}
	condition := map[interface{}]interface{}{
		"status": 0,
		"cid":    cid,
	}
	if condition["cid"] == "0" {
		delete(condition, "cid")
	}
	articles := models.GetAllArticle(condition)

	c.HTML(200, "index.html", gin.H{
		"userInfo":  userInfo,
		"last_time": lastTime,
		"navs":      models.GetCategory(),
		"list":      articles,
	})
}

// 验证码
func Captcha(c *gin.Context) {
	common.ServeHTTP(c.Writer, c.Request)
}
