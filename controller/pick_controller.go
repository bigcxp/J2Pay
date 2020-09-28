package controller

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/pkg/util"
	"j2pay-server/service"
	"strconv"
)

// @Tags 商户提领代发管理
// @Summary 商户提领、代发列表
// @Produce json
// @Param status  query int false "-1：等待中，1:执行中，2：成功，3：取消，4，失败"
// @Param type  query int false "0：所有，1：代发，2：提领"
// @Param orderCode    query string false "商户订单编号"
// @Param userId query int false "商户id"
// @Param from_date  query string false "起"
// @Param to_date    query string false "至"
// @Param page query int false "页码"
// @Param pageSize query int false "每页显示多少条"
// @Success 200 {object} response.MerchantPickSendPage
// @Router /merchantPick [get]
func MerchantPickIndex(c *gin.Context) {
	response := util.Response{c}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	code := c.Query("orderCode")
	FromDate := c.Query("from_date")
	ToDate := c.Query("to_date")
	status, _ := strconv.Atoi(c.Query("status"))
	userId, _ := strconv.Atoi(c.Query("userId"))
	types, _ := strconv.Atoi(c.Query("type"))
	res, err := service.MerchantPickList(FromDate,ToDate,code,types,userId,status,page,pageSize)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(res)
}

// @Tags 商户提领代发管理
// @Summary 管理端商户提领列表
// @Produce json
// @Param status  query int false "-1：等待中，1:执行中，2：成功，3：取消，4，失败"
// @Param type  query int false "0：所有，1：代发，2：提领"
// @Param userId query int false "商户id"
// @Param from_date  query string false "起"
// @Param to_date    query string false "至"
// @Param page query int false "页码"
// @Param pageSize query int false "每页显示多少条"
// @Success 200 {object} response.PickUpPage
// @Router /pick [get]
func PickIndex(c *gin.Context) {
	response := util.Response{c}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	FromDate := c.Query("from_date")
	ToDate := c.Query("to_date")
	status, _ := strconv.Atoi(c.Query("status"))
	userId, _ := strconv.Atoi(c.Query("userId"))
	types, _ := strconv.Atoi(c.DefaultQuery("type", "2"))
	res, err := service.PickUpList(FromDate,ToDate,status,types,userId,page,pageSize)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(res)
}

// @Tags 商户提领代发管理
// @Summary 管理端商户代发列表
// @Produce json
// @Param status  query int false "-1：等待中，1:执行中，2：成功，3：取消，4，失败"
// @Param type  query int false "0：所有，1：代发，2：提领"
// @Param orderCode    query string false "商户订单编号"
// @Param userId query int false "商户id"
// @Param from_date  query string false "起"
// @Param to_date    query string false "至"
// @Param page query int false "页码"
// @Param pageSize query int false "每页显示多少条"
// @Success 200 {object} response.SendPage
// @Router /send [get]
func SendIndex(c *gin.Context) {
	response := util.Response{c}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	code := c.Query("orderCode")
	FromDate := c.Query("from_date")
	ToDate := c.Query("to_date")
	status, _ := strconv.Atoi(c.Query("status"))
	userId, _ := strconv.Atoi(c.Query("userId"))
	types, _ := strconv.Atoi(c.Query("type"))
	res, err := service.SendList(FromDate,ToDate,code,status,types,userId,page,pageSize)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(res)
}

// @Tags 商户提领代发管理
// @Summary 商户提领代发详情
// @Produce json
// @Param id path uint true "ID"
// @Success 200 {object} response.MerchantPickList
// @Router /merchantPick/{id} [get]
func MerchantPickDetail(c *gin.Context) {
	response := util.Response{c}
	id, _ := strconv.Atoi(c.Param("id"))
	detail, err := service.MerchantPickDetail(uint(id))
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(detail)
}

// @Tags 商户提领代发管理
// @Summary 管理端商户提领详情
// @Produce json
// @Param id path uint true "ID"
// @Success 200 {object} response.PickList
// @Router /pick/{id} [get]
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

// @Tags 商户提领代发管理
// @Summary 管理端商户代发详情
// @Produce json
// @Param id path uint true "ID"
// @Success 200 {object} response.SendList
// @Router /send/{id} [get]
func SendDetail(c *gin.Context) {
	response := util.Response{c}
	id, _ := strconv.Atoi(c.Param("id"))
	detail, err := service.SendDetail(uint(id))
	if err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessData(detail)
}



// @Tags 商户提领代发管理
// @Summary 提领
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

// @Tags 商户提领代发管理
// @Summary 代发
// @Produce json
// @Param body body request.SendAdd true "商户代发"
// @Success 200
// @Router /merchantSend [post]
func SendAdd(c *gin.Context)  {
	response := util.Response{c}
	var send request.SendAdd
	if err := c.ShouldBindJSON(&send); err != nil {
		response.SetValidateError(err)
		return
	}
	if err := service.SendAdd(send); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("成功")

}

// @Tags 商户提领代发管理
// @Summary 修改提领订单状态
// @Produce json
// @Param id path int true "ID"
// @Param body body request.SendEdit true "提领订单"
// @Router /pick/{id} [put]
func PickEdit(c *gin.Context) {
	response := util.Response{c}
	var send request.SendEdit
	send.Id, _ = strconv.Atoi(c.Param("id"))
	if err := c.ShouldBindJSON(&send); err != nil {
		response.SetValidateError(err)
		return
	}
	if err := service.CancelPick(send); err != nil {
		response.SetOtherError(err)
		return
	}
	response.SuccessMsg("成功")
}

// @Tags 商户提领代发管理
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
