package models

import (
	"encoding/base64"
	"encoding/json"
	"gobbs/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义用户的struct
type User struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Last_time int64  `json:"last_time"`
	Password  string `json:"-"`
	Avatar    string `json:"avatar"`
	Views     int64  `json:"views"`
}

type Gousers struct {
	User
}

// 登录
func Login(user string) (userInfo User) {
	userInfo = User{}
	common.MysqlDb.QueryRow("select id,name,last_time,password,avatar from gousers where name = ?", user).Scan(&userInfo.Id, &userInfo.Name, &userInfo.Last_time, &userInfo.Password, &userInfo.Avatar)
	userInfo.Avatar = setAvatar(userInfo.Avatar)
	return
}

// 获取头像
func setAvatar(avatar string) string {
	if avatar == "" {
		return "http://tva1.sinaimg.cn/crop.0.0.118.118.180/5db11ff4gw1e77d3nqrv8j203b03cweg.jpg"
	}
	return avatar
}

// 获取用户昵称
func GetUserName(id string) string {
	condition := map[interface{}]interface{}{
		"id": id,
	}
	if id == "0" || id == "" {
		return "无"
	}
	find := common.ORM.From(&Gousers{}).Where(condition).Find("name")
	if find["Name"] != "" {
		return find["Name"].(string)
	}
	return ""
}

// 获取用户基本信息
func GetUserBase(uid int64) map[string]interface{} {
	condition := map[interface{}]interface{}{
		"id": uid,
	}
	find := common.ORM.From(&Gousers{}).Where(condition).Find("name,avatar")
	find["Avatar"] = setAvatar(find["Avatar"].(string))
	return find
}

// 获取用户基本信息
func GetUserProfile(uid int64, field string) map[string]interface{} {
	condition := map[interface{}]interface{}{
		"id": uid,
	}
	find := common.ORM.From(&Gousers{}).Where(condition).Find(field)
	find["Avatar"] = setAvatar(find["Avatar"].(string))
	return find
}

// 保存cookie
func SaveLoginCookie(c *gin.Context, userInfo User) bool {
	jsonData, _ := json.Marshal(userInfo)
	data := base64.StdEncoding.EncodeToString(jsonData)
	cookie := &http.Cookie{
		Name:   "user",
		Value:  data,
		MaxAge: 3600,
	}
	http.SetCookie(c.Writer, cookie)
	return true
}

// 删除cookie
func DeleteLoginCookie(c *gin.Context, cookieName string) bool {
	cookie := &http.Cookie{
		Name:   cookieName,
		MaxAge: -1,
	}
	http.SetCookie(c.Writer, cookie)
	return true
}

// 获取用户信息
func GetUserInfo(cookie string) (userInfo User) {
	userInfo = User{}
	data, _ := base64.StdEncoding.DecodeString(cookie)
	json.Unmarshal(data, &userInfo)
	return
}

// 更新用户资料
func UpdateUserProfile(id int64, telephone string, avatar string) bool {
	sql := "UPDATE gousers set telephone = ?,avatar = ? WHERE id = ?"
	res, _ := common.MysqlDb.Exec(sql, telephone, avatar, id)
	affect, _ := res.RowsAffected()
	if affect > 0 {
		return true
	}
	return false
}

// 更新密码
func UpdatePassword(id int64, password string) bool {
	sql := "UPDATE gousers set password = ? WHERE id = ?"
	res, _ := common.MysqlDb.Exec(sql, password, id)
	affect, _ := res.RowsAffected()
	if affect > 0 {
		return true
	}
	return false
}
