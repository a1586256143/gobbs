package router

import (
	"../common"
	"../models"
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

func DetailArticle(c *gin.Context){
	id := c.Param("id")
	condition := map[interface{}]interface{}{
		"id":id ,
	}
	find := common.ORM.From(&models.Article{}).Where(condition).Find("id,uid,title,content,create_time")
	if find["CreateTime"] != "0" {
		times , _ := strconv.ParseInt(find["CreateTime"].(string) , 10 , 64)
		find["FormatTime"] = common.DateFormat(times)
	}
	find["NickName"] = models.GetUserName(find["Uid"].(string))
	captchaId := common.GetCaptchaId()
	// 查询评论
	commentMaps := map[interface{}]interface{}{
		"aid" : find["Id"],
		"status":0 ,
	}
	comments := models.GetAllComments(commentMaps)
	commentCount := common.ORM.From(&models.Comments{}).Where(commentMaps).Count("id")
	c.HTML(http.StatusOK , "detail.html" , gin.H{
		"find" : find,
		"ImageUrl":"/captcha/" + captchaId + ".png",
		"captchaId":captchaId ,
		"comments":comments,
		"commentsCount" : commentCount,
	})
}