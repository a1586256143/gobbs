package common

import (
	"crypto/sha1"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"strings"
	"time"
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
			c.JSON(200 , Error(v.Msg))
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

//如果是小写字母, 则变换为大写字母
func StrFirstToUpper(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArry := []rune(str)
	if strArry[0] >= 97 && strArry[0] <= 122  {
		strArry[0] -=  32
	}
	return string(strArry)
}

// 把字段名转换为驼峰名create_time转换CreateTime
func AutoTuoFeng(item string) string {
	items := strings.Split(item , "_")
	for i , v := range items {
		items[i] = StrFirstToUpper(v)
	}
	key := strings.Join(items , "")
	return key
}