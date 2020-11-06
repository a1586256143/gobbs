package main

import (
	"gobbs/common"
	"gobbs/models"
	handler "gobbs/router"
	"html/template"

	"github.com/gin-gonic/gin"
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
	// 引入静态文件
	router.Static("/static", "./static")
	router.Static("/upload", "./upload")
	// 注入模板函数
	router.SetFuncMap(template.FuncMap{
		"dateFormat": common.DateFormat,
		"str2html":   common.Str2html,
	})
	// 加载HTML模板
	router.LoadHTMLGlob("templates/*")
	// 首页
	router.GET("/", AuthMiddleware(), handler.Index)

	// 详情页面
	router.GET("/article/:id", handler.DetailArticle)

	// 发布页面
	router.GET("/publish", AuthMiddleware(), handler.PublishArticlePage)

	// 发布处理
	router.POST("/publish", AuthMiddleware(), handler.PublishArticle)

	// 发布评论
	router.POST("/comments", AuthMiddleware(), handler.PublishComment)

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
	router.GET("/captcha/:captchaId", handler.Captcha)

	// 重载验证码
	router.GET("/captchaReload", handler.ReloadCaptcha)

	// 上传图片
	router.POST("/upload", AuthMiddleware(), handler.UploadImg)

	// 个人空间
	router.GET("/space", AuthMiddleware(), handler.SpacePage)

	// 其它用户的空间
	router.GET("/space/:uid", handler.OtherSpacePage)

	// 加好友
	router.POST("/add-friends/:uid", AuthMiddleware(), handler.AddFriends)

	// 我的消息
	router.GET("/message", AuthMiddleware(), handler.Message)

	// 同意或拒绝请求
	router.POST("/accept", AuthMiddleware(), handler.Accept)

	// 我的好友
	router.GET("/friends", AuthMiddleware(), handler.Friends)

	// 删除好友
	router.POST("/deleteFriend", AuthMiddleware(), handler.DelFriends)

	// 个人信息
	router.GET("/profile", AuthMiddleware(), handler.Profile)

	// 修改个人信息
	router.POST("/profile", AuthMiddleware(), handler.SaveProfile)

	// 修改密码
	router.GET("/password", AuthMiddleware(), handler.Password)
	router.POST("/password", AuthMiddleware(), handler.Password)

	_ = router.Run(":9999")
}
