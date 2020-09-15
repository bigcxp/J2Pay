package controller

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/model/request"
	"j2pay-server/pkg/util"
	"j2pay-server/service"
	"strconv"
)

// @Tags 系统公告
// @Summary 获取系统公告列表
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页显示多少条"
// @Success 200 {object} response.SystemMessagePage
// @Router /systemMessage [get]
func SystemMessage(c *gin.Context)  {
	response := util.Response{c}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	res, err := service.MessageList(page, pageSize)
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
	if err := c.ShouldBindJSON(&message); err != nil {
		response.SetValidateError(err).SetMeta(map[string]string{"Title": "标题"})
		return
	}
	if err := service.MessageAdd(message); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("添加成功")
}
