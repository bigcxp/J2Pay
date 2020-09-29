package controller

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/model/request"
	"j2pay-server/pkg/util"
	"j2pay-server/service"
	"strconv"
)

// @Tags 系统参数管理
// @Summary 系统参数详情
// @Produce json
// @Success 200 {object} response.Parameter
// @Router /system [get]
func SystemParameter(c *gin.Context) {
	response := util.Response{c}
	detail := service.GetDetail()
	response.SuccessData(detail)
}

// @Tags 系统参数管理
// @Summary 更新系统参数
// @Produce json
// @Param id path int true "ID"
// @Param body body request.ParameterEdit true "系统参数"
// @Router /system/{id} [put]
func SystemParameterEdit(c *gin.Context) {
	response := util.Response{c}
	var parameter request.ParameterEdit
	parameter.Id, _ = strconv.Atoi(c.Param("id"))
	if err := c.ShouldBindJSON(&parameter); err != nil {
		response.SetValidateError(err)
		return
	}
	if err := service.UpdateParameter(parameter); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("成功")
}

// @Tags 系统参数管理
// @Summary 更新GasPrice
// @Produce json
// @Param id path int true "ID"
// @Param body body request.ParameterEdit true "系统参数"
// @Router /systemGasPrice/{id} [put]
func SystemGasPriceEdit(c *gin.Context) {
	response := util.Response{c}
	var parameter request.ParameterEdit
	parameter.Id, _ = strconv.Atoi(c.Param("id"))
	if err := c.ShouldBindJSON(&parameter); err != nil {
		response.SetValidateError(err)
		return
	}
	if err := service.UpdateGasPrice(parameter); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("成功")
}


