package common

import (
	"crypto/sha1"
	"fmt"
	"reflect"

	// "gobbs/models"
	"html/template"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 校验的struct
type Validate struct {
	Field string
	Msg   string
}

/**
 * 验证助手
 * @param {[type]} c        gin.Context [description]
 * @param {[type]} validate []Validate) (success      bool [description]
 */
func ValidateHelper(c gin.Context, validate []Validate) (success bool) {
	success = true
	for _, v := range validate {
		if c.PostForm(v.Field) == "" && success == true {
			c.JSON(200, Error(v.Msg))
			success = false
			break
		}
	}
	return
}

// 选项、参数
type Option struct {
	K string      // 名字
	V interface{} // 值
}

// 选项是否为空
func (o Option) IsEmpty() bool {
	return reflect.DeepEqual(o, Option{})
}

// 获取时间戳
func GetUnix() int64 {
	return time.Now().Unix()
}

// 格式化时间
func DateFormat(unix int64) string {
	formatTime := time.Unix(unix, 0).Format("2006-01-02 15:04:05")
	return formatTime
}

// 密码加密
func PasswordEncode(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

//如果是小写字母, 则变换为大写字母
func StrFirstToUpper(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArry := []rune(str)
	if strArry[0] >= 97 && strArry[0] <= 122 {
		strArry[0] -= 32
	}
	return string(strArry)
}

// 把字段名转换为驼峰名create_time转换CreateTime
func AutoTuoFeng(item string) string {
	items := strings.Split(item, "_")
	for i, v := range items {
		items[i] = StrFirstToUpper(v)
	}
	key := strings.Join(items, "")
	return key
}

// 格式化时间
func FormatTime(time string) string {
	times, _ := strconv.ParseInt(time, 10, 64)
	return DateFormat(times)
}

// 字符转html
func Str2html(content string) interface{} {
	var a interface{}
	a = content
	fmt.Println(a)
	return template.HTML(a.(string))
}

// 获取用户ID
func GetUid(c *gin.Context) (cookie string) {
	cookie, _ = c.Cookie("user")
	return
}

// 获取今天的日期
func GetToday() string {
	date := time.Now()
	return fmt.Sprintf("%d%d%d", date.Year(), date.Month(), date.Day())
}

// 获取一周的日期
func GetWeek() (startDay, endDay string) {
	date := time.Now()
	var (
		weekDay = map[string]int{
			"Sunday":    0,
			"Monday":    1,
			"Tuesday":   2,
			"Wednesday": 3,
			"Thursday":  4,
			"Friday":    5,
			"Saturday":  6,
		}
	)
	// 将获取的WeekDay转换成string
	end := (date.Weekday()).String()
	plus := int64(86400 * weekDay[end])
	owe := int64(86400 * (7 - 1 - weekDay[end]))
	endDay = time.Unix(date.Unix()+plus+86399, 0).Format("20060102")
	startDay = time.Unix(date.Unix()-owe, 0).Format("20060102")
	return startDay, endDay
}

// 合并两个map
// @param item2 合并到 item1
func MergeMap(item1, item2 map[string]interface{}) map[string]interface{} {
	for item, value := range item2 {
		item1[item] = value
	}
	return item1
}

// 获取数据类型
func TypeOf(v interface{}) string {
	return fmt.Sprintf("%T", v)
}
