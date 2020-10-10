package controller

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/model/request"
	"j2pay-server/pkg/util"
	"j2pay-server/service"
	"net/http"
	"strconv"
)

// @Tags 退款订单管理
// @Summary 退款订单管理
// @Produce json
// @Router /returnIndex [get]
func ReturnIndex(c *gin.Context) {
	c.HTML(200,"return.html" ,gin.H{
		"code": http.StatusOK,
	})
}

// @Tags 退款订单管理
// @Summary 退款订单列表
// @Produce json
// @Param status  query int false "状态 1：退款等待中，2：退款中，3：退款失败，4：已退款"
// @Param orderCode    query string false "商户订单编号"
// @Param from_date  query string false "起"
// @Param to_date    query string false "至"
// @Param page query int false "页码"
// @Param pageSize query int false "每页显示多少条"
// @Success 200 {object} response.ReturnPage
// @Router /return [get]
func ReturnList(c *gin.Context) {
	response := util.Response{c}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	orderCode := c.Query("")
	FromDate := c.Query("orderCode")
	ToDate := c.Query("to_date")
	status, _ := strconv.Atoi(c.Query("status"))
	res, err := service.ReturnList(FromDate,ToDate,status,orderCode,page,pageSize)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(res)
}

// @Tags 退款订单管理
// @Summary 退款订单详情
// @Produce json
// @Param id path uint true "ID"
// @Success 200 {object} response.ReturnList
// @Router /return/{id} [get]
func ReturnDetail(c *gin.Context) {
	response := util.Response{c}
	id, _ := strconv.Atoi(c.Param("id"))
	detail, err := service.ReturnDetail(uint(id))
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(detail)
}

// @Tags 退款订单管理
// @Summary 新增退款订单
// @Produce json
// @Param body body request.ReturnAdd true "新增退款订单"
// @Success 200
// @Router /return [post]
func ReturnAdd(c *gin.Context) {
	response := util.Response{c}
	var returns request.ReturnAdd
	if err := c.ShouldBindJSON(&returns); err != nil {
		response.SetValidateError(err)
		return
	}
	if err := service.ReturnAdd(returns); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("成功")

}
