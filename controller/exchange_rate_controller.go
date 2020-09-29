package controller

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/model/request"
	"j2pay-server/pkg/util"
	"j2pay-server/service"
	"strconv"
)

// @Tags 汇率管理
// @Summary 所有汇率列表
// @Produce json
// @Success 200 {object} response.RatePage
// @Router /rate [get]
func RateList(c *gin.Context) {
	response := util.Response{c}
	detail := service.GetAllRate()
	response.SuccessData(detail)
}

// @Tags 汇率管理
// @Summary 获取汇率详情
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} response.RateList
// @Router /rate/{id} [get]
func RateDetail(c *gin.Context) {
	response := util.Response{c}
	id, _ := strconv.Atoi(c.Param("id"))
	detail, err := service.RateDetail(uint(id))
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(detail)
}

// @Tags 汇率管理
// @Summary 更新汇率
// @Produce json
// @Param id path int true "ID"
// @Param body body request.RateEdit true "汇率"
// @Router /rate/{id} [put]
func RateEdit(c *gin.Context) {
	response := util.Response{c}
	var rate request.RateEdit
	rate.Id, _ = strconv.Atoi(c.Param("id"))
	if err := c.ShouldBindJSON(&rate); err != nil {
		response.SetValidateError(err)
		return
	}
	if err := service.UpdateRate(rate); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("成功")
}
