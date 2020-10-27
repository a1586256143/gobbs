package router

import (
	"fmt"
	// "fmt"
	"gobbs/common"
	"gobbs/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 登录方法
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

// 注册方法
func Register(c *gin.Context) {
	user := c.PostForm("name")
	password := c.PostForm("password")
	telephone := c.PostForm("telephone")
	userInfo := models.Login(user)
	if userInfo.Name != "" {
		// 说明账号已经被注册了
		c.JSON(200, common.Error("账号已被注册"))
	} else {
		t := common.GetUnix()
		ret, _ := common.MysqlDb.Exec("insert into gousers(name,password,telephone,create_time,last_time) values (?,?,?,?,?)", user, common.PasswordEncode(password), telephone, t, t)
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

// 退出登录
func Logout(c *gin.Context) {
	models.DeleteLoginCookie(c, "user")
	c.Redirect(302, "/login")
}

// 登录页面
func RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

// 注册页面
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

// 个人中心
func SpacePage(c *gin.Context) {
	cid := c.DefaultQuery("cid", "0")
	uid := c.Param("uid")
	cookie, _ := c.Cookie("user")
	// 获取用户信息
	userInfo := models.GetUserInfo(cookie)
	// 自己刷新时不更新浏览数
	if uid != "" && uid != strconv.FormatInt(userInfo.Id, 10) {
		models.UpdateViews(uid)
	}
	conditionUid := uid
	if conditionUid == "" {
		conditionUid = strconv.FormatInt(userInfo.Id, 10)
	}
	condition := map[interface{}]interface{}{
		"status": 0,
		"cid":    cid,
		"uid":    conditionUid,
	}
	if condition["cid"] == "0" {
		delete(condition, "cid")
	}
	articles := models.GetAllArticle(condition)
	data := gin.H{
		"list": articles,
	}
	commonData := UserNavData(c, common.Option{K: "uid", V: uid})
	c.HTML(http.StatusOK, "space.html", common.MergeMap(data, commonData))
}

// 添加好友
func AddFriends(c *gin.Context) {
	uid, _ := strconv.ParseInt(c.Param("uid"), 10, 64)
	content := c.PostForm("content")
	if uid == 0 {
		c.JSON(200, common.Error("无效的用户信息"))
		return
	}
	cookie, _ := c.Cookie("user")
	// 获取用户信息
	userInfo := models.GetUserInfo(cookie)
	exists := models.IsExists(userInfo.Id, uid)
	if exists > 0 {
		c.JSON(200, common.Success("请求发送成功，请等待对方同意申请"))
		return
	}
	isFriends := models.IsFriends(userInfo.Id, uid)
	if isFriends > 0 {
		c.JSON(200, common.Success("你们已经是好友了"))
		return
	}
	status := models.AddFriends(userInfo.Id, uid, content)
	if status > 0 {
		c.JSON(200, common.Success("请求发送成功，请等待对方同意申请"))
		return
	}
	c.JSON(200, common.Error("请求发送失败"))
}

// 我的消息
func Message(c *gin.Context) {
	menuType := c.DefaultQuery("type", "1")
	cookie, _ := c.Cookie("user")
	// 获取用户信息
	userInfo := models.GetUserInfo(cookie)
	var list map[int]map[string]interface{}
	if menuType == "1" {
		condition := map[interface{}]interface{}{
			"status":   0,
			"pull_uid": userInfo.Id,
		}
		// 查询好友申请信息
		list = models.GetMyApplyList(condition)
	}
	data := gin.H{
		"list": list,
	}
	commonData := UserNavData(c, common.Option{})
	c.HTML(http.StatusOK, "message.html", common.MergeMap(data, commonData))
}

// 用户导航公共数据
func UserNavData(c *gin.Context, option common.Option) gin.H {
	fmt.Println(option)
	cookie, _ := c.Cookie("user")
	// 获取用户信息
	userInfo := models.GetUserInfo(cookie)
	lastTime := "暂无"
	// 处理最后登录时间
	if userInfo.Last_time != 0 {
		lastTime = common.DateFormat(userInfo.Last_time)
	}
	// 申请的数量
	applyNumber := models.GetFriendsApplyNumber(userInfo.Id)
	// 我的好友数量
	friendsNumber := models.GetFriendsNumber(userInfo.Id)
	var uid string
	uid = strconv.FormatInt(userInfo.Id, 10)
	if option.K == "uid" && option.V != "" {
		uid = option.V.(string)
	}
	fmt.Println(userInfo)
	return gin.H{
		"today":         models.GetTodayNumber(uid), // 当天的统计
		"week":          models.GetWeekNumber(uid),  // 本周的统计
		"month":         models.GetWeekNumber(uid),  // 本月的访问统计
		"userInfo":      userInfo,
		"lastTime":      lastTime,
		"messageNumber": applyNumber,
		"friendsNumber": friendsNumber,
	}
}

// 同意或拒绝好友请求
func Accept(c *gin.Context) {
	id, _ := strconv.ParseInt(c.PostForm("id"), 10, 64)
	status, _ := strconv.ParseInt(c.PostForm("status"), 10, 64)
	res := models.Accept(id, status)
	if res == false {
		c.JSON(200, common.Error("请求处理失败"))
		return
	}
	c.JSON(200, common.Success("处理成功"))
}

// 我的好友
func Friends(c *gin.Context) {
	data := gin.H{}
	commonData := UserNavData(c, common.Option{})
	c.HTML(http.StatusOK, "friends.html", common.MergeMap(data, commonData))
}
