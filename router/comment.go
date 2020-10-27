package router

import (
	"gobbs/common"
	"gobbs/models"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PublishComment(c *gin.Context) {
	// 校验验证码
	verifyCode := c.PostForm("captcha")
	captchaId := c.PostForm("captchaId")
	content := c.PostForm("content")
	log.Println(verifyCode, content)
	if common.VerifyCaptcha(captchaId, verifyCode) == false {
		c.JSON(200, common.Error("validImg"))
		return
	}
	cookie, _ := c.Cookie("user")
	uid := models.GetUserInfo(cookie).Id
	aid, _ := strconv.Atoi(c.PostForm("aid"))
	if models.AddComments(int(uid), aid, content, 1) == true {
		c.JSON(200, common.Success("发布成功"))
		return
	}
	c.JSON(200, common.Error("发布失败"))
}
