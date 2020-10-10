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
// @Router /adminUserIndex [get]
func AdminUserIndex(c *gin.Context) {
	c.HTML(200,"admin-user.html" ,gin.H{
		"code": http.StatusOK,
	})
}

// @Tags 账户管理
// @Summary 获取账户列表
// @Produce json
// @Param Pid  query int false "默认0：查商户列表，1:查商户账户列表"
// @Param page query int false "页码"
// @Param pageSize query int false "每页显示多少条"
// @Success 200 {object} response.AdminUserPage
// @Router /adminUser [get]
func UserIndex(c *gin.Context) {
	response := util.Response{c}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	//name := c.Query("name")
	Pid,_ :=strconv.Atoi(c.Query("Pid"))
	//if utf8.RuneCountInString(name) > 32 {
	//	name = string([]rune(name)[:32])
	//}
	res, err := service.UserList(Pid, page, pageSize)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(res)
}

// @Tags 账户管理
// @Summary 获取角色树
// @Produce json
// @Success 200 {object} response.AdminUserPage
// @Router /auth/role [get]
func RoleTree(c *gin.Context) {
	response := util.Response{c}
	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	res, err := service.RoleTree(id, 0)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(res)
}

// @Tags 账户管理
// @Summary 获取账户详情
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.AdminUserList
// @Router /adminUser/{id} [get]
func UserDetail(c *gin.Context) {
	response := util.Response{c}
	id, _ := strconv.Atoi(c.Param("id"))
	detail, err := service.UserDetail(id)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(detail)
}

// @Tags 账户管理
// @Summary 添加账户
// @Produce json
// @Param body body request.UserAdd true "用户"
// @Router /adminUser [post]
func UserAdd(c *gin.Context) {
	response := util.Response{c}
	var user request.UserAdd
	if err := c.ShouldBindJSON(&user); err != nil {
		response.SetValidateError(err)
		return
	}
	if err := service.UserAdd(user); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("添加成功")
}

// @Tags 账户管理
// @Summary 编辑账户
// @Produce json
// @Param id path int true "用户ID"
// @Param body body request.UserEdit true "用户"
// @Router /adminUser/{id} [put]
func UserEdit(c *gin.Context) {
	response := util.Response{c}
	var user request.UserEdit
	user.Id, _ = strconv.Atoi(c.Param("id"))
	if err := c.ShouldBindJSON(&user); err != nil {
		response.SetValidateError(err)
		return
	}
	if err := service.UserEdit(user); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("编辑成功")
}

// @Tags 账户管理
// @Summary 删除账户
// @Produce json
// @Param id path int true "用户ID"
// @Router /adminUser/{id} [delete]
func UserDel(c *gin.Context) {
	response := util.Response{c}
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.UserDel(id); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("删除成功")
}



