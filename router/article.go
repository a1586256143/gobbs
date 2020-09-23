package router

import (
	"gobbs/common"
	"gobbs/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func PublishArticle(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	cid, _ := strconv.ParseInt(c.PostForm("cid"), 10, 64)
	// 校验开始
	validate := []common.Validate{
		{"title", "标题不能为空",},
		{"content", "内容不能为空",},
		{"cid", "分类不能为空",},
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
func PublishArticlePage(c *gin.Context) {
	c.HTML(http.StatusOK, "publish.html", gin.H{
		"category": models.GetCategory(),
	})
}

/**
 * 刷新验证码
 * @param {[type]} c *gin.Context [description]
 */
func ReloadCaptcha(c *gin.Context){
	captchaId := common.GetCaptchaId()
	c.JSON(http.StatusOK , gin.H{
		"ImageUrl":"/captcha/" + captchaId + ".png",
		"captchaId":captchaId ,
	})
}

func DetailArticle(c *gin.Context){
	id := c.Param("id")
	condition := map[interface{}]interface{}{
		"id":id ,
	}
	find := common.ORM.From(&models.Article{}).Where(condition).Find("id,uid,title,content,create_time")
	if find["CreateTime"] != "0" {
		find["FormatTime"] = common.FormatTime(find["CreateTime"].(string))
	}
	find["NickName"] = models.GetUserName(find["Uid"].(string))
	captchaId := common.GetCaptchaId()
	// 查询评论
	commentMaps := map[interface{}]interface{}{
		"aid" : find["Id"],
		"status":1 ,
	}
	comments := models.GetAllComments(commentMaps)
	for _ , v := range comments {
		if v["CreateTime"] != "0" {
			v["FormatTime"] = common.FormatTime(v["CreateTime"].(string));
		}
	}
	commentCount := common.ORM.From(&models.Comments{}).Where(commentMaps).Count("id")
	c.HTML(http.StatusOK , "detail.html" , gin.H{
		"find" : find,
		"ImageUrl":"/captcha/" + captchaId + ".png",
		"captchaId":captchaId ,
		"comments":comments,
		"commentsCount" : commentCount,
	})
}