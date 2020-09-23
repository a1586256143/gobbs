package main

import (
	"gobbs/common"
	"gobbs/models"
	handler "gobbs/router"
	"github.com/gin-gonic/gin"
	"html/template"
)

// 中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Cookie("user"); err == nil {
			userInfo := models.GetUserInfo(cookie)
			if userInfo.Id != 0 {
				c.Next()
				return
			}
		}
		c.Redirect(302, "/login")
		c.Abort()
		return
	}
}
func main() {
	router := gin.Default()
	// 加载HTML模板
	router.LoadHTMLGlob("templates/*")
	// 引入静态文件
	router.Static("/static", "./static")
	// 注入模板函数
	router.SetFuncMap(template.FuncMap{
		"dateFormat": common.DateFormat,
	})
	// 首页
	router.GET("/", AuthMiddleware(), handler.Index)

	// 详情页面
	router.GET("/article/:id", AuthMiddleware(), handler.DetailArticle)

	// 发布页面
	router.GET("/publish" , AuthMiddleware(), handler.PublishArticlePage)

	// 发布处理
	router.POST("/publish" , AuthMiddleware(), handler.PublishArticle)

	// 发布评论
	router.POST("/comments" , AuthMiddleware() , handler.PublishComment)

	// 登录页面
	router.GET("/login", handler.LoginPage)

	// 注册页面
	router.GET("/register", handler.RegisterPage)

	// 处理登录
	router.POST("/login", handler.Login)

	// 处理注册
	router.POST("/register", handler.Register)

	// 退出登录
	router.GET("/logout", handler.Logout)

	// 获取验证码
	router.GET("/captcha/:captchaId" , handler.Captcha)

	_ = router.Run(":9999")
}
