package router

import (
	"gobbs/common"
	"github.com/gin-gonic/gin"
	"log"
	"gobbs/models"
	"strconv"
)

func PublishComment(c *gin.Context) {
	// 校验验证码
	verifyCode := c.PostForm("captcha")
	captchaId := c.PostForm("captchaId")
	content := c.PostForm("content")
	log.Println(verifyCode , content)
	if common.VerifyCaptcha(captchaId , verifyCode) == false {
		c.JSON(200 , common.Error("验证码错误"))
		return
	}
	cookie, _ := c.Cookie("user")
	uid := models.GetUserInfo(cookie).Id
	aid , _ := strconv.Atoi(c.PostForm("aid"))
	if models.AddComments(int(uid), aid , content , 1) == true {
		c.JSON(200 , common.Success("发布成功"))
		return
	}
	c.JSON(200 , common.Error("发布失败"))
}