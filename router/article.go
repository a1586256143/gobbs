package router

import (
	"fmt"
	"gobbs/common"
	"gobbs/models"

	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 发布文章
func PublishArticle(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	cid, _ := strconv.ParseInt(c.PostForm("cid"), 10, 64)
	// 校验开始
	validate := []common.Validate{
		{"title", "标题不能为空"},
		{"content", "内容不能为空"},
		{"cid", "分类不能为空"},
	}
	success := common.ValidateHelper(*c, validate)
	if success == false {
		return
	}
	cookie, _ := c.Cookie("user")
	uid := models.GetUserInfo(cookie).Id
	status := models.AddArticle(uid, title, content, cid)
	if status == 0 {
		c.JSON(200, common.Error("发布失败，请重新尝试"))
		return
	}
	c.JSON(200, common.Success("发布成功"))
}

// 发布文章页面
func PublishArticlePage(c *gin.Context) {
	c.HTML(http.StatusOK, "publish.html", gin.H{
		"category": models.GetCategory(),
	})
}

/**
 * 刷新验证码
 * @param {[type]} c *gin.Context [description]
 */
func ReloadCaptcha(c *gin.Context) {
	captchaId := common.GetCaptchaId()
	c.JSON(http.StatusOK, gin.H{
		"ImageUrl":  "/captcha/" + captchaId + ".png",
		"captchaId": captchaId,
	})
}

// 帖子详情
func DetailArticle(c *gin.Context) {
	id := c.Param("id")
	condition := map[interface{}]interface{}{
		"id": id,
	}
	find := common.ORM.From(&models.Article{}).Where(condition).Find("id,uid,title,content,create_time")
	if find["CreateTime"] != "0" {
		find["FormatTime"] = common.FormatTime(find["CreateTime"].(string))
	}
	var content interface{}
	content = find["Content"]
	find["Content"] = template.HTML(content.(string))
	find["NickName"] = models.GetUserName(find["Uid"].(string))
	captchaId := common.GetCaptchaId()
	// 查询评论
	commentMaps := map[interface{}]interface{}{
		"aid":    find["Id"],
		"status": 1,
	}
	comments := models.GetAllComments(commentMaps)
	for _, v := range comments {
		if v["CreateTime"] != "0" {
			v["FormatTime"] = common.FormatTime(v["CreateTime"].(string))
		}
	}
	// 获取用户信息
	userInfo := models.GetUserInfo(common.GetUid(c))
	commentCount := common.ORM.From(&models.Comments{}).Where(commentMaps).Count("id")
	// 是否为好友
	pullUid, _ := strconv.ParseInt(find["Uid"].(string), 10, 64)
	var isFriend int64
	isFriend = 0
	if pullUid != userInfo.Id {
		isFriend = models.IsFriends(userInfo.Id, pullUid)
	}
	fmt.Println("isFriend", isFriend)
	find["Uid"] = pullUid
	c.HTML(http.StatusOK, "detail.html", gin.H{
		"isFriend":      isFriend,
		"userInfo":      userInfo,
		"find":          find,
		"ImageUrl":      "/captcha/" + captchaId + ".png",
		"captchaId":     captchaId,
		"comments":      comments,
		"commentsCount": commentCount,
		"login":         common.GetUid(c),
	})
}

// 上传图片方法
func uploadHanlder(c *gin.Context) (savePath string, err error) {
	header, err := c.FormFile("file")
	if err != nil {
		fmt.Println("upload error", err)
		return "", err
	}
	date := time.Now()
	savePath = fmt.Sprintf("upload/%d%d%d", date.Year(), date.Month(), date.Day())
	exists, _ := PathExists(savePath)
	if !exists {
		os.MkdirAll(savePath, 0777)
	}
	savePath = savePath + "/" + header.Filename
	if err := c.SaveUploadedFile(header, savePath); err != nil {
		fmt.Println("upload error", err)
		return "", err
	}
	fmt.Println("Upload success")
	return "/" + savePath, nil
}

// 上传图片控制器
func UploadImg(c *gin.Context) {
	path, err := uploadHanlder(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 404,
			"msg":  "上传错误",
			"data": "{}",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "",
		"data": gin.H{
			"src": path,
		},
	})
}

// 目录是否存在
// @param string path 校验的目录
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, nil
}
