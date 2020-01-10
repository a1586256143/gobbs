package main

import (
	"./common"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"./models"
)

// 中间件
func AuthMiddleware() gin.HandlerFunc{
	return func(c *gin.Context){
		if cookie , err := c.Cookie("user"); err == nil {
			userInfo := models.GetUserInfo(cookie)
			if userInfo.Id != 0{
				c.Next()
				return
			}
		}
		c.Redirect(302 , "/login")
		c.Abort()
		return
	}
}

func main(){
	router := gin.Default()
	// 加载HTML模板
	router.LoadHTMLGlob("templates/*")
	// 引入静态文件
	router.Static("/static","./static")
	// 注入模板函数
	router.SetFuncMap(template.FuncMap{
		"dateFormat" : common.DateFormat ,
	})
	// 首页
	router.GET("/" , AuthMiddleware() , func(c *gin.Context){
		cid := c.DefaultQuery("cid" , "0")
		fmt.Println("cid = " , cid)
		list := []struct{
			Name string
			Id int
		}{
			{Name : "标题1" , Id:1} ,
			{Name : "标题2" , Id:2} ,
		}
		cookie , _ := c.Cookie("user")
		// 获取用户信息
		userInfo := models.GetUserInfo(cookie)
		lastTime := "暂无"
		if userInfo.Last_time != 0 {
			lastTime = common.DateFormat(userInfo.Last_time)
		}
		c.HTML(200 , "index.html" , gin.H{
			"userInfo" : userInfo,
			"last_time" : lastTime ,
			"list" : list ,
		})
	})

	// 登录页面
	router.GET("/login" , func(c *gin.Context){
		c.HTML(http.StatusOK , "login.html" , gin.H{})
	})

	// 注册页面
	router.GET("/register" , func(c *gin.Context){
		c.HTML(http.StatusOK , "register.html" , gin.H{})
	})

	// 处理登录
	router.POST("/login" , func(c *gin.Context){
		user := c.PostForm("name")
		password := c.PostForm("password")
		fmt.Println("password" , common.PasswordEncode(password))
		// 调用登录方法
		userInfo := models.Login(user)
		if userInfo.Id != 0{
			// 保存用户登陆Cookie信息
			models.SaveLoginCookie(c , userInfo)
			t := common.GetUnix()
			common.MysqlDb.Exec("update gousers set last_time = ? where id = ?" , t , userInfo.Id)
			c.JSON(200 , gin.H{
				"status" : 0 ,
				"msg" : "登录成功",
			})
		}else {
			c.JSON(200 , gin.H{
				"status" : 1 ,
				"msg" : "登录失败",
			})
		}
	})

	// 处理注册
	router.POST("/register" , func(c *gin.Context){
		user := c.PostForm("name")
		userInfo := models.Login(user)
		if userInfo.Name != ""{
			// 说明账号已经被注册了
			c.JSON(200 , gin.H{
				"status" : 1 ,
				"msg" : "账号已被注册",
			})
		}else {
			t := common.GetUnix()
			ret , _ := common.MysqlDb.Exec("insert into gousers(name,telephone,create_time,last_time) values (?,?)" , user , "123456" , t , t)
			rowsAffected , _ := ret.RowsAffected()
			lastInsertID,_ := ret.LastInsertId()
			if rowsAffected > 0 {
				userInfo = models.User{Id : int(lastInsertID) , Name : user , Last_time:t}
				// 保存登录信息
				models.SaveLoginCookie(c , userInfo)
				c.JSON(200 , gin.H{
					"status" : 0 ,
					"msg" : "注册成功，将自动登录...",
				})
			}else {
				c.JSON(200 , gin.H{
					"status" : 1 ,
					"msg" : "注册失败，请重试",
				})
			}
		}
	})

	// 退出登录
	router.GET("/logout" , func(c *gin.Context){
		models.DeleteLoginCookie(c , "user")
		c.Redirect(302 , "/login")
	})

	// 发布页面
	router.GET("/publish" , func(c *gin.Context){
		c.HTML(http.StatusOK , "publish.html" , gin.H{})
	})

	router.Run(":9999")
}