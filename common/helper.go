package common

import (
	"crypto/sha1"
	"fmt"
	"io"
	"time"
	"github.com/gin-gonic/gin"
)

// 校验的struct
type Validate struct{
	Field string 
	Msg string
}

/**
 * 验证助手
 * @param {[type]} c        gin.Context [description]
 * @param {[type]} validate []Validate) (success      bool [description]
 */
func ValidateHelper(c gin.Context , validate []Validate) (success bool) {
	success = true
	for _ , v := range validate {
		if c.PostForm(v.Field) == "" && success == true {
			c.JSON(200 , gin.H{
				"status" : 1 ,
				"msg" : v.Msg ,
			})
			success = false
			break
		}
	}
	return
}

// 获取时间戳
func GetUnix() int64 {
	return time.Now().Unix()
}

// 格式化时间
func DateFormat(unix int64) string {
	formatTime := time.Unix(unix , 0).Format("2006-01-02 15:04:05")
	return formatTime
}

// 密码加密
func PasswordEncode(data string) string {
	t := sha1.New()
	io.WriteString(t , data)
	return fmt.Sprintf("%x",t.Sum(nil))
}

func Select(condition map[interface{}]interface{}){
	for k , v := range condition {
		switch v.(type){
		case string : 
			fmt.Println("is string" , v)
		case float64 : 
			fmt.Println("is float64" , v)
		case int : 
			fmt.Println("is int" , v)
		default : 
			fmt.Println("unkown" , k , v)
		}
	}
}