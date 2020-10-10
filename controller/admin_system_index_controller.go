package controller

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/model/request"
	"j2pay-server/pkg/util"
	"j2pay-server/service"
	"strconv"
)


// @Tags 管理员首页
// @Summary 系统首页
// @Produce json
// @Router /systemIndex [get]
func Index(c *gin.Context) {
	response := util.Response{c}
	response.Index("登录")
}

// @Tags 管理员首页
// @Summary 首页数据
// @Produce json
// @Success 200 {object} response.SystemMessagePage
// @Router /index [get]
func SystemIndex(c *gin.Context) {
	response := util.Response{c}
	response.SuccessData("待写")
}

// @Tags 管理员首页
// @Summary 修改密码
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} response.Password
// @Router /password/{id} [put]
func UpdatePassword(c *gin.Context) {
	response := util.Response{c}
	Id, _ := strconv.Atoi(c.Param("id"))
	if err := c.Bind(&Id); err != nil {
		response.SetValidateError(err)
		return
	}
	password, err1 := service.UpdatePassword(Id)
	if err1 != nil {
		return
	}
	response.SuccessData(password)
}

// @Tags 管理员首页
// @Summary 开启google验证
// @Produce json
// @Param id path int true "ID"
// @Param body body request.Google true "google参数"
// @Router /google/{id} [put]
func GoogleValidate(c *gin.Context) {
	response := util.Response{c}
	var google request.Google
	google.Id, _ = strconv.Atoi(c.Param("id"))
	if err := c.ShouldBindJSON(&google); err != nil {
		response.SetValidateError(err)
		return
	}
	if err := service.OpenGoogle(google); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("成功")
}
