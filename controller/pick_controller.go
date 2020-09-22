package controller

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/pkg/util"
	"j2pay-server/service"
	"strconv"
	"unicode/utf8"
)

// @Tags 商户管理
// @Summary 商户提领列表
// @Produce json
// @Param status  query int false "-1：等待中，1:执行中，2：成功，3：取消，4，失败"
// @Param type  query int false "0:全部，1：代发，2：提领"
// @Param name    query string false "组织名称"
// @Param code    query string false "系统编号"
// @Param from_date  query string false "起"
// @Param to_date    query string false "至"
// @Param page query int false "页码"
// @Param pageSize query int false "每页显示多少条"
// @Success 200 {object} model.PickUpPage
// @Router /merchantPick [get]
func PickIndex(c *gin.Context) {
	response := util.Response{c}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.Query("name")
	code := c.Query("code")
	FromDate := c.Query("name")
	ToDate := c.Query("name")
	status, _ := strconv.Atoi(c.Query("status"))
	types, _ := strconv.Atoi(c.Query("type"))
	if utf8.RuneCountInString(name) > 32 {
		name = string([]rune(name)[:32])
	}
	res, err := service.PickList(FromDate,ToDate,status, name,code,types, page, pageSize)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(res)
}

// @Tags 商户管理
// @Summary 获取提领详情
// @Produce json
// @Param id path uint true "ID"
// @Success 200 {object} model.Pick
// @Router /merchantPick/{id} [get]
func PickDetail(c *gin.Context) {
	response := util.Response{c}
	id, _ := strconv.Atoi(c.Param("id"))
	detail, err := service.PickDetail(uint(id))
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(detail)
}

// @Tags 商户管理
// @Summary 提领，代发
// @Produce json
// @Param body body request.PickAdd true "商户提领"
// @Success 200
// @Router /merchantPick [post]
func PickAdd(c *gin.Context)  {
	response := util.Response{c}
	var pick request.PickAdd
	if err := c.ShouldBindJSON(&pick); err != nil {
		response.SetValidateError(err)
		return
	}
	if err := service.PickAdd(pick); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("成功")

}

// @Tags 商户管理
// @Summary 通知
// @Produce json
// @Param token  query string true "token"
// @Success 302
// @Router /notify [post]
func PickNotify(c *gin.Context)  {
	response := util.Response{c}
	token := c.Query("token")
	adminUser := model.GetUserByWhere("token =?", token)
	c.Redirect(302, adminUser.DaiUrl)
	c.Abort()
	response.SuccessMsg("成功")
}
