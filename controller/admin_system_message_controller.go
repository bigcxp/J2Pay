package controller

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/model/request"
	"j2pay-server/pkg/util"
	"j2pay-server/service"
	"net/http"
	"strconv"
	"unicode/utf8"
)



// @Tags 系统公告
// @Summary 系统公告
// @Produce json
// @Router /systemMessageIndex [get]
func SystemMessageIndex(c *gin.Context) {
	c.HTML(200,"systemMessage.html" ,gin.H{
		"code": http.StatusOK,
	})
}

// @Tags 系统公告
// @Summary 获取系统公告列表
// @Produce json
// @Param title query string false "标题"
// @Param page query int false "页码"
// @Param pageSize query int false "每页显示多少条"
// @Success 200 {object} response.SystemMessagePage
// @Router /systemMessage [get]
func SystemMessage(c *gin.Context)  {
	response := util.Response{c}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	title := c.Query("title")
	res, err := service.MessageList(title,page, pageSize)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(res)

}

// @Tags 系统公告
// @Summary 获取用户公告列表
// @Produce json
// @Param username query string true "用户名"
// @Param page query int false "页码"
// @Param pageSize query int false "每页显示多少条"
// @Success 200 {object} response.AdminUserMessagePage
// @Router /systemMessageByUser [get]
func SystemMessageByUserId(c *gin.Context)  {
	response := util.Response{c}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	username := c.Query("username")

	if utf8.RuneCountInString(username) > 32 {
		username = string([]rune(username)[:32])
	}
	res, err := service.MessageListByUser(username, page, pageSize)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(res)


}


// @Tags 系统公告
// @Summary 添加系统公告
// @Produce json
// @Param body body request.MessageAdd true "系统公告"
// @Success 200
// @Router /systemMessage [post]
func SystemMessageAdd(c *gin.Context) {
	response := util.Response{c}
	var message request.MessageAdd
	if err := c.ShouldBind(&message); err != nil {
		response.SetValidateError(err).SetMeta(map[string]string{"Title": "标题"})
		return
	}
	if err := service.MessageAdd(message); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("添加成功")
}

// @Tags 系统公告
// @Summary 删除公告
// @Produce json
// @Param id path int true "公告ID"
// @Router /systemMessage/{id} [delete]
func SystemMessageDel(c *gin.Context)  {
	response := util.Response{c}
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.MessageDel(id); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("删除成功")
}

// @Tags 系统公告
// @Summary 编辑公告
// @Produce json
// @Param id path int true "公告ID"
// @Param body body request.MessageEdit true "公告"
// @Router /systemMessage/{id} [put]
func SystemMessageEdit(c *gin.Context)  {
	response := util.Response{c}
	var  systemMessage request.MessageEdit
	systemMessage.ID, _ = strconv.Atoi(c.Param("id"))
	if err := c.ShouldBind(&systemMessage); err != nil {
		response.SetValidateError(err)
		return
	}
	if err := service.MessageEdit(systemMessage); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("编辑成功")
}
