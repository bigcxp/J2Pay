package controller

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/pkg/util"
	"j2pay-server/service"
	"strconv"
	"unicode/utf8"
)

// @Tags 商户管理
// @Summary 商户提领列表
// @Produce json
// @Param status  query int false "-1：等待中，1:执行中，2：成功，3：取消，4，失败"
// @Param name    query string false "组织名称"
// @Param from_date  query string false "起"
// @Param to_date    query string false "至"
// @Param page query int false "页码"
// @Param pageSize query int false "每页显示多少条"
// @Success 200 {object} model.PickUpPage
// @Router /merchantPick [get]
func PickIndex(c *gin.Context) {
	response := util.Response{c}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.Query("name")
	FromDate := c.Query("name")
	ToDate := c.Query("name")
	status, _ := strconv.Atoi(c.Query("status"))
	if utf8.RuneCountInString(name) > 32 {
		name = string([]rune(name)[:32])
	}

	res, err := service.PickList(FromDate,ToDate,status, name, page, pageSize)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(res)
}

// @Tags 商户管理
// @Summary 获取提领详情
// @Produce json
// @Param id path uint true "ID"
// @Success 200 {object} response.AdminUserList
// @Router /merchantPick/{id} [get]
func PickDetail(c *gin.Context) {
	response := util.Response{c}
	id, _ := strconv.Atoi(c.Param("id"))
	detail, err := service.PickDetail(uint(id))
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(detail)
}
