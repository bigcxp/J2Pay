package controller

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/model/request"
	"j2pay-server/pkg/util"
	"j2pay-server/service"
	"net/http"
	"strconv"
)


// @Tags 实收明细记录
// @Summary 实收明细记录
// @Produce json
// @Router /detailedRecordIndex [get]
func DetailedRecordIndex(c *gin.Context) {
	c.HTML(200,"detailed.html" ,gin.H{
		"code": http.StatusOK,
	})
}

// @Tags 实收明细记录
// @Summary 实收明细记录列表
// @Produce json
// @Param status  query int false "状态 1：未绑定，2：已绑定"
// @Param idCode    query string false "系统编号"
// @Param txid    query string false "交易hash"
// @Param address    query string false "收款地址"
// @Param from_date  query string false "起"
// @Param to_date    query string false "至"
// @Param page query int false "页码"
// @Param pageSize query int false "每页显示多少条"
// @Success 200 {object} response.Erc20Page
// @Router /detail [get]
func DetailedList(c *gin.Context) {
	response := util.Response{c}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	idCode := c.Query("idCode")
	txid := c.Query("txid")
	address := c.Query("address")
	FromDate := c.Query("from_date")
	ToDate := c.Query("to_date")
	status, _ := strconv.Atoi(c.Query("status"))
	res, err := service.Erc20List(c,idCode,address,txid,FromDate,ToDate,status,page,pageSize)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(res)
}

// @Tags 实收明细记录
// @Summary 实收记录详情
// @Produce json
// @Param id path uint true "ID"
// @Success 200 {object} response.Erc20List
// @Router /detail/{id} [get]
func DetailedDetail(c *gin.Context) {
	response := util.Response{c}
	id, _ := strconv.Atoi(c.Param("id"))
	detail, err := service.Erc20Detail(id)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(detail)
}

// @Tags 实收明细记录
// @Summary 新增实收明细记录
// @Produce json
// @Param body body request.Erc20Add true "实收记录"
// @Success 200
// @Router /detail [post]
func DetailedAdd(c *gin.Context)  {
	response := util.Response{c}
	var detailed request.Erc20Add
	if err := c.ShouldBind(&detailed); err != nil {
		response.SetValidateError(err)
		return
	}
	if err := service.Erc20Add(detailed); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("成功")

}

// @Tags 实收明细记录
// @Summary 绑定 | 解绑 订单
// @Produce json
// @Param id path int true "ID"
// @Param body body request.Erc20Edit true "明细记录"
// @Router /detail/{id} [put]
func DetailedEdit(c *gin.Context) {
	response := util.Response{c}
	var detail request.Erc20Edit
	detail.ID, _ = strconv.Atoi(c.Param("id"))
	if err := c.ShouldBind(&detail); err != nil {
		response.SetValidateError(err)
		return
	}
	if err := service.IsBindOrder(detail); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("成功")
}