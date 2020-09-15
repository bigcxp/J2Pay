package controller

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/pkg/util"
)

// @Tags 管理员首页
// @Summary
// @Produce json
// @Success 200 {object} response.SystemMessagePage
// @Router /index [get]
func IndexSystem(c *gin.Context)  {
	response := util.Response{c}
	response.SuccessData("待写")
}
