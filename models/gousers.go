package models

import (
	"gobbs/common"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 定义用户的struct
type User struct{
	Id int64 			`json:"id"`
	Name string 	`json:"name"`
	Last_time int64 	`json:"last_time"`
	Password string `json:"-"`
}

type Gousers struct{
	User
}

// 登录
func Login(user string) (userInfo User) {
	userInfo = User{}
	common.MysqlDb.QueryRow("select id,name,last_time,password from gousers where name = ?", user).Scan(&userInfo.Id, &userInfo.Name , &userInfo.Last_time , &userInfo.Password)
	return
}

// 获取用户昵称
func GetUserName(id string) string {
	condition := map[interface{}]interface{}{
		"id" : id ,
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

// 保存cookie
func SaveLoginCookie(c *gin.Context , userInfo User) bool {
	jsonData, _ := json.Marshal(userInfo)
	data := base64.StdEncoding.EncodeToString(jsonData)
	cookie := &http.Cookie{
		Name:  "user",
		Value: data,
		MaxAge : 3600,
	}
	http.SetCookie(c.Writer, cookie)
	return true
}

// 删除cookie
func DeleteLoginCookie(c *gin.Context , cookieName string) bool {
	cookie := &http.Cookie{
		Name:  cookieName,
		MaxAge : -1,
	}
	http.SetCookie(c.Writer, cookie)
	return true
}

// 获取用户信息
func GetUserInfo(cookie string)(userInfo User){
	userInfo = User{}
	data , _ := base64.StdEncoding.DecodeString(cookie)
	json.Unmarshal(data, &userInfo)
	return
}