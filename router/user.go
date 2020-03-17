package router

import (
	"../common"
	"../models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	user := c.PostForm("name")
	password := c.PostForm("password")
	// 调用登录方法
	userInfo := models.Login(user)
	if userInfo.Id != 0 {
		if common.PasswordEncode(password) != userInfo.Password {
			c.JSON(200, common.Error("用户名或密码错误"))
			return
		}
		// 保存用户登陆Cookie信息
		models.SaveLoginCookie(c, userInfo)
		t := common.GetUnix()
		common.MysqlDb.Exec("update gousers set last_time = ? where id = ?", t, userInfo.Id)
		c.JSON(200, common.Success("登录成功"))
	} else {
		c.JSON(200, common.Error("登录失败"))
	}
}

func Register(c *gin.Context) {
	user := c.PostForm("name")
	userInfo := models.Login(user)
	if userInfo.Name != "" {
		// 说明账号已经被注册了
		c.JSON(200, common.Error("账号已被注册"))
	} else {
		t := common.GetUnix()
		ret, _ := common.MysqlDb.Exec("insert into gousers(name,telephone,create_time,last_time) values (?,?)", user, "123456", t, t)
		rowsAffected, _ := ret.RowsAffected()
		lastInsertID, _ := ret.LastInsertId()
		if rowsAffected > 0 {
			userInfo = models.User{Id: lastInsertID, Name: user, Last_time: t}
			// 保存登录信息
			models.SaveLoginCookie(c, userInfo)
			c.JSON(200, common.Success("注册成功，将自动登录..."))
		} else {
			c.JSON(200, common.Error("注册失败，请重试"))
		}
	}
}

func Logout(c *gin.Context) {
	models.DeleteLoginCookie(c, "user")
	c.Redirect(302, "/login")
}

func RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}