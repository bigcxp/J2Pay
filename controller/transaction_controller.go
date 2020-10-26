package controller

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/pkg/util"
	"j2pay-server/service"
	"strconv"
)

// @Tags 交易记录管理
// @Summary eth交易记录
// @Produce json
// @Param fromTime  query int false "起始时间"
// @Param toTime  query int false "到达时间"
// @Param page query int false "页码"
// @Param pageSize query int false "每页显示多少条"
// @Success 200 {object} response.AdminUserPage
// @Router /ethTransfer [get]
func EthTransfer(c *gin.Context) {
	response := util.Response{c}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	fromTime, _ := strconv.Atoi(c.Query("fromTime"))
	toTime, _ := strconv.Atoi(c.Query("toTime"))
	res, err := service.EthTxList(fromTime, toTime, page, pageSize)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(res)
}

// @Tags 交易记录管理
// @Summary 热钱包交易明细
// @Produce json
// @Param fromAddr   query string false "来源地址"
// @Param toAddr   query string false "目的地址"
// @Param fromTime  query int false "起始时间"
// @Param toTime  query int false "到达时间"
// @Param scheduleStatus  query int false "排程状态"
// @Param chainStatus  query int false "链上状态"
// @Param page query int false "页码"
// @Param pageSize query int false "每页显示多少条"
// @Success 200 {object} response.AdminUserPage
// @Router /hotTransfer [get]
func HotTransfer(c *gin.Context) {
	response := util.Response{c}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	fromTime, _ := strconv.Atoi(c.Query("fromTime"))
	toTime, _ := strconv.Atoi(c.Query("toTime"))
	scheduleStatus, _ := strconv.Atoi(c.Query("scheduleStatus"))
	chainStatus, _ := strconv.Atoi(c.Query("chainStatus"))
	fromAddr := c.Query("fromAddr")
	toAddr := c.Query("toAddr")
	res, err := service.HotTxList(fromAddr,toAddr,scheduleStatus,chainStatus,fromTime,toTime, page, pageSize)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(res)
}
