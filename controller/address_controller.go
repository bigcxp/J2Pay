package controller

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/model/request"
	"j2pay-server/pkg/util"
	"j2pay-server/service"
	"strconv"
)

// @Tags 地址管理
// @Summary 钱包地址列表
// @Produce json
// @Param status  query int false "状态 0：所有，1：已完成，2：执行中，3：结账中"
// @Param handleStatus  query int false "指派状态 0：所有，1：启用，2：停用"
// @Param address    query string false "钱包地址"
// @Param userId query int true "商户组织id"
// @Param useTag query int true "钱包占用类型 -2:eth钱包地址,-1：hot钱包地址,1:商户充币钱包地址"
// @Param page query int false "页码"
// @Param pageSize query int false "每页显示多少条"
// @Success 200 {object} response.AddressPage
// @Router /addrList [get]
func AddrList(c *gin.Context) {
	response := util.Response{c}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "1"))
	handStatus, _ := strconv.Atoi(c.DefaultQuery("handStatus", "0"))
	userId, _ := strconv.Atoi(c.DefaultQuery("userId", "0"))
	useTag, _ := strconv.Atoi(c.DefaultQuery("useTag", "1"))
	address := c.Query("address")
	res, err := service.AddressList(address, status, handStatus, userId, useTag, page, pageSize)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(res)
}

// @Tags 地址管理
// @Summary 为商户分配地址 生成热钱包地址 生成eth地址
// @Produce json
// @Param body body request.AddressAdd true "钱包地址"
// @Success 200
// @Router /createAddress [post]
func CreateAddress(c *gin.Context) {
	response := util.Response{c}
	var addr request.AddressAdd
	if err := c.ShouldBindJSON(&addr); err != nil {
		response.SetValidateError(err)
		return
	}
	if err := service.AddAddress(addr); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("成功")
}

// @Tags 地址管理
// @Summary 启用 停用地址
// @Produce json
// @Param body body request.OpenOrStopAddress true "地址id"
// @Router /addrRestart [post]
func AddrRestart(c *gin.Context) {
	response := util.Response{c}
	var addr request.OpenOrStopAddress
	if err := c.ShouldBindJSON(&addr); err != nil {
		response.SetValidateError(err)
		return
	}
	if err := service.RestartAddr(addr); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("成功")
}

// @Tags 地址管理
// @Summary 更新余额
// @Produce json
// @Param body body request.UpdateAmount true "地址id"
// @Router /updateBalance [post]
func UpdateBalance(c *gin.Context) {
	response := util.Response{c}
	var addr request.UpdateAmount
	if err := c.ShouldBindJSON(&addr); err != nil {
		response.SetValidateError(err)
		return
	}
	if err := service.UpdateBalance(addr); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("成功")
}
