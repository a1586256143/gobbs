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
	cookie, _ := c.Cookie("user")
	// 获取用户信息
	userInfo := models.GetUserInfo(cookie)
	condition := map[interface{}]interface{}{
		"status": 0,
		"cid":    cid,
		"uid":    strconv.FormatInt(userInfo.Id, 10),
	}
	if condition["cid"] == "0" {
		delete(condition, "cid")
	}
	articles := models.GetAllArticle(condition)
	data := gin.H{
		"list": articles,
	}
	commonData := UserNavData(c, common.Option{})
	c.HTML(http.StatusOK, "space.html", common.MergeMap(data, commonData))
}

// 其它用户个人中心
func OtherSpacePage(c *gin.Context) {
	cid := c.DefaultQuery("cid", "0")
	uid := c.Param("uid")
	models.UpdateViews(uid)
	condition := map[interface{}]interface{}{
		"status": 0,
		"cid":    cid,
		"uid":    uid,
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
	fmt.Println("options is ", option)
	var (
		userInfo      models.User
		applyNumber   int64
		friendsNumber int64
		uid           string
		other         int
		titlePrefix   string
	)
	lastTime := "暂无"
	if option.IsEmpty() {
		cookie, _ := c.Cookie("user")
		// 获取用户信息
		userInfo = models.GetUserInfo(cookie)
		// 处理最后登录时间
		if userInfo.Last_time != 0 {
			lastTime = common.DateFormat(userInfo.Last_time)
		}
		// 申请的数量
		applyNumber = models.GetFriendsApplyNumber(userInfo.Id)
		// 我的好友数量
		friendsNumber = models.GetFriendsNumber(userInfo.Id)
		uid = strconv.FormatInt(userInfo.Id, 10)
		other = 0
		titlePrefix = "我"
	} else if option.K == "uid" && option.V != "" {
		parseUid, _ := strconv.ParseInt(option.V.(string), 10, 64)
		base := models.GetUserBase(parseUid)
		userInfo = models.User{Id: parseUid, Name: base["Name"].(string), Avatar: base["Avatar"].(string)}
		uid = option.V.(string)
		other = 1
		titlePrefix = "他(她)"
	}
	return gin.H{
		"today":         models.GetTodayNumber(uid), // 当天的统计
		"week":          models.GetWeekNumber(uid),  // 本周的统计
		"month":         models.GetWeekNumber(uid),  // 本月的访问统计
		"userInfo":      userInfo,
		"lastTime":      lastTime,
		"messageNumber": applyNumber,
		"friendsNumber": friendsNumber,
		"other":         other,
		"titlePrefix":   titlePrefix,
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
	list := models.GetMyFriendsList(commonData["userInfo"].(models.User).Id)
	data["list"] = list
	c.HTML(http.StatusOK, "friends.html", common.MergeMap(data, commonData))
}

// 删除好友
func DelFriends(c *gin.Context) {
	id := c.DefaultPostForm("id", "0")
	if id == "0" {
		c.JSON(200, common.Error("好友不存在"))
		return
	}
	cookie, _ := c.Cookie("user")
	uid := models.GetUserInfo(cookie).Id
	status := models.DeleteFriend(id, uid)
	if status {
		c.JSON(200, common.Success("删除成功"))
		return
	}
	c.JSON(200, common.Error("删除失败"))
}

// 个人资料
func Profile(c *gin.Context) {
	data := gin.H{}
	commonData := UserNavData(c, common.Option{})
	profile := models.GetUserProfile(commonData["userInfo"].(models.User).Id, "*")
	data["profile"] = profile
	c.HTML(http.StatusOK, "profile.html", common.MergeMap(data, commonData))
}

// 保存个人资料
func SaveProfile(c *gin.Context) {
	phone := c.DefaultPostForm("telephone", "")
	avatar := c.DefaultPostForm("avatar", "")
	// 校验开始
	validate := []common.Validate{
		{"telephone", "手机号不能为空"},
	}
	success := common.ValidateHelper(*c, validate)
	if success == false {
		return
	}
	cookie, _ := c.Cookie("user")
	uid := models.GetUserInfo(cookie).Id
	status := models.UpdateUserProfile(uid, phone, avatar)
	if status == false {
		c.JSON(200, common.Error("更新失败，请重新尝试"))
		return
	}
	c.JSON(200, common.Success("更新成功"))
}

// 修改密码
func Password(c *gin.Context) {
	if c.Request.Method == "POST" {
		password := c.PostForm("password")
		newPassword := c.PostForm("new_password")
		rePassword := c.PostForm("re_password")
		// 校验开始
		validate := []common.Validate{
			{"password", "原密码不能为空"},
			{"new_password", "新密码不能为空"},
			{"re_password", "确认密码不能为空"},
		}
		success := common.ValidateHelper(*c, validate)
		if success == false {
			return
		}
		if rePassword != newPassword {
			c.JSON(200, common.Error("两次密码输入不一致"))
			return
		}
		cookie, _ := c.Cookie("user")
		uid := models.GetUserInfo(cookie).Id
		userPassword := models.GetUserProfile(uid, "password,avatar")
		if userPassword["Password"] != common.PasswordEncode(password) {
			c.JSON(200, common.Error("原密码不正确"))
			return
		}
		status := models.UpdatePassword(uid, common.PasswordEncode(newPassword))
		if status == false {
			c.JSON(200, common.Error("修改失败"))
			return
		}
		models.DeleteLoginCookie(c, "user")
		c.JSON(200, common.Success("修改成功，请重新登录"))
		return
	}
	data := gin.H{}
	commonData := UserNavData(c, common.Option{})
	// profile := models.GetUserProfile(commonData["userInfo"].(models.User).Id, "*")
	// data["profile"] = profile
	c.HTML(http.StatusOK, "password.html", common.MergeMap(data, commonData))
}
