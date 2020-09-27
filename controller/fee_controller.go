package controller

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/model/request"
	"j2pay-server/pkg/util"
	"j2pay-server/service"
	"strconv"
)

// @Tags 手续费结账
// @Summary 手续费列表
// @Produce json
// @Param status  query int false "状态 1：等待中，2：已完成"
// @Param userId    query int false "商户id"
// @Param from_date  query string false "起"
// @Param to_date    query string false "至"
// @Param page query int false "页码"
// @Param pageSize query int false "每页显示多少条"
// @Success 200 {object} response.FeePage
// @Router /fee [get]
func FeeList(c *gin.Context) {
	response := util.Response{c}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	userId,_:= strconv.Atoi(c.Query("userId"))
	FromDate := c.Query("fromDate")
	ToDate := c.Query("toDate")
	status, _ := strconv.Atoi(c.Query("status"))
	res, err := service.FeeList(FromDate,ToDate,status,userId,page,pageSize)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(res)
}


// @Tags 手续费结账
// @Summary 结账
// @Produce json
// @Param id path int true "ID"
// @Param body body request.FeeEdit true "手续费"
// @Router /fee/{id} [put]
func Settle(c *gin.Context) {
		response := util.Response{c}
		var fee request.FeeEdit
		fee.Id, _ = strconv.Atoi(c.Param("id"))
		if err := c.ShouldBindJSON(&fee); err != nil {
			response.SetValidateError(err)
			return
		}
		if err := service.FeeSettle(fee); err != nil {
			response.SetOtherError(err)
			return
		}
		response.SuccessMsg("编辑成功")

}
