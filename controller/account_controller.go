package controller

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/model/request"
	"j2pay-server/pkg/util"
	"j2pay-server/service"
	"net/http"
	"strconv"
)

// @Tags 账户管理
// @Summary 账户管理
// @Produce json
// @Router /accountIndex [get]
func AccountIndex(c *gin.Context) {
	c.HTML(200, "account.html", gin.H{
		"code": http.StatusOK,
	})
}


// @Tags 账户管理
// @Summary 获取角色树
// @Produce json
// @Success 200 {object} response.AdminUserPage
// @Router /auth/role [get]
func RoleTree(c *gin.Context) {
	response := util.Response{c}
	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	res, err := service.RoleTree(id)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(res)
}

// @Tags 账户管理
// @Summary 获取账户列表
// @Produce json
// @Param username  query string false "账户名称"
// @Param page query int false "页码"
// @Param pageSize query int false "每页显示多少条"
// @Success 200 {object} response.AccountPage
// @Router /accountList [get]
func AccountList(c *gin.Context) {
	response := util.Response{c}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	username := c.Query("username")
	res, err := service.AccountList(username, page, pageSize)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(res)
}

// @Tags 账户管理
// @Summary 获取账户详情
// @Produce json
// @Param id path int true "账户ID"
// @Success 200 {object} response.AccountList
// @Router /account/{id} [get]
func AccountDetail(c *gin.Context) {
	response := util.Response{c}
	id, _ := strconv.Atoi(c.Param("id"))
	detail, err := service.AccountDetail(int64(id))
	if err != nil {
		 response.SetOtherError(err)
		return
	}
	response.SuccessData(detail)
}

// @Tags 账户管理
// @Summary 添加账户
// @Produce json
// @Param body body request.AccountAdd true "账户"
// @Router /account [post]
func Account(c *gin.Context) {
	response := util.Response{c}
	var account request.AccountAdd
	if err := c.ShouldBind(&account); err != nil {
		response.SetValidateError(err)
		return
	}
	if err := service.AccountAdd(account); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("添加成功")
}

// @Tags 账户管理
// @Summary 编辑账户
// @Produce json
// @Param id path int true "用户ID"
// @Param body body request.AccountEdit true "账户"
// @Router /account/{id} [put]
func AccountEdit(c *gin.Context) {
	response := util.Response{c}
	var account request.AccountEdit
	account1, _ := strconv.Atoi(c.Param("id"))
	account.ID = int64(account1)
	if err := c.ShouldBind(&account); err != nil {
		response.SetValidateError(err)
		return
	}
	if err := service.AccountEdit(account); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("编辑成功")
}

// @Tags 账户管理
// @Summary 删除账户
// @Produce json
// @Param id path int true "用户ID"
// @Router /account/{id} [delete]
func AccountDel(c *gin.Context) {
	response := util.Response{c}
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.UserDel(int64(id)); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("删除成功")
}

//
