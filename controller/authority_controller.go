// 权限一般都是程序员添加的，无需写增删改的逻辑
// 只需查的逻辑
package controller

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/pkg/util"
	"j2pay-server/service"
)

// @Tags 角色管理
// @Summary 获取权限树
// @Produce json
// @Success 200 {object} response.RolePage
// @Router /auth/tree [get]
func AuthTree(c *gin.Context) {
	response := util.Response{c}
	res := service.AuthTreeCache()
	response.SuccessData(res)
}

// @Tags 角色管理
// @Summary 权限列表
// @Produce json
// @Success 200 {object} response.RolePage
// @Router /auth/list [get]
func AuthList(c *gin.Context) {
	response := util.Response{c}
	res := service.AuthListCache()
	response.SuccessData(res)
}