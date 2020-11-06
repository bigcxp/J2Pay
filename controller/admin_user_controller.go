package controller

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/model/request"
	"j2pay-server/pkg/util"
	"j2pay-server/service"
	"net/http"
	"strconv"
)

// @Tags 组织管理
// @Summary 组织管理
// @Produce json
// @Router /adminUserIndex [get]
func AdminUserIndex(c *gin.Context) {
	c.HTML(200,"adminUser.html" ,gin.H{
		"code": http.StatusOK,
	})
}

// @Tags 组织管理
// @Summary 获取组织列表
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页显示多少条"
// @Success 200 {object} response.AdminUserPage
// @Router /adminUser [get]
func UserIndex(c *gin.Context) {
	response := util.Response{c}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	res, err := service.UserList( page, pageSize)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(res)
}



// @Tags 组织管理
// @Summary 获取组织详情
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.AdminUserList
// @Router /adminUser/{id} [get]
func UserDetail(c *gin.Context) {
	response := util.Response{c}
	id, _ := strconv.Atoi(c.Param("id"))
	detail, err := service.UserDetail(int64(id))
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(detail)
}

// @Tags 组织管理
// @Summary 添加组织
// @Produce json
// @Param body body request.UserAdd true "用户"
// @Router /adminUser [post]
func UserAdd(c *gin.Context) {
	response := util.Response{c}
	var user request.UserAdd
	if err := c.ShouldBind(&user); err != nil {
		response.SetValidateError(err)
		return
	}
	if err := service.UserAdd(user); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("添加成功")
}

// @Tags 组织管理
// @Summary 编辑组织
// @Produce json
// @Param id path int true "组织ID"
// @Param body body request.UserEdit true "用户"
// @Router /adminUser/{id} [put]
func UserEdit(c *gin.Context) {
	response := util.Response{c}
	var user request.UserEdit
	account1, _ := strconv.Atoi(c.Param("id"))
	user.ID = int64(account1)
	if err := c.ShouldBind(&user); err != nil {
		response.SetValidateError(err)
		return
	}
	user.ID = int64(account1)
	if err := service.UserEdit(user); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("编辑成功")
}

// @Tags 组织管理
// @Summary 删除组织
// @Produce json
// @Param id path int true "用户ID"
// @Router /adminUser/{id} [delete]
func UserDel(c *gin.Context) {
	response := util.Response{c}
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.UserDel(int64(id)); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("删除成功")
}



