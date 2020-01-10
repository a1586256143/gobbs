package common

import (
	"crypto/sha1"
	"fmt"
	"io"
	"time"
)

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