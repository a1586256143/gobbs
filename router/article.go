package router

import (
	"../common"
	"../models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func PublishArticle(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	cid, _ := strconv.ParseInt(c.PostForm("cid"), 10, 64)
	// 校验开始
	validate := []common.Validate{
		{"title", "标题不能为空",},
		{"content", "内容不能为空",},
		{"cid", "分类不能为空",},
	}
	success := common.ValidateHelper(*c, validate)
	if success == false {
		return
	}
	cookie, _ := c.Cookie("user")
	uid := models.GetUserInfo(cookie).Id
	status := models.AddArticle(uid, title, content, cid)
	if status == 0 {
		c.JSON(200, gin.H{
			"status": 1,
			"msg":    "发布失败，请重新尝试",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 0,
		"msg":    "发布成功",
	})
}
func PublishArticlePage(c *gin.Context) {
	c.HTML(http.StatusOK, "publish.html", gin.H{
		"category": models.GetCategory(),
	})
}