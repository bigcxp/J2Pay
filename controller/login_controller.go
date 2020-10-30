package controller

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/middleware"
	"j2pay-server/model/request"
	"j2pay-server/pkg/setting"
	"j2pay-server/pkg/util"
	"j2pay-server/service"
	"net/http"
)

// @Tags 登录操作
// @Summary 登录操作
// @Produce json
// @Router /login [get]
func LoginIndex(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{
		"code": http.StatusOK,
	})
}

// @Tags 登录操作
// @Summary 登录操作
// @Produce json
// @Param body body request.LoginUser true "用户"
// @Success 200
// @Router /login [post]
func Login(c *gin.Context) {
	response := util.Response{c}
	user := request.LoginUser{}
	if err := c.ShouldBindJSON(&user); err != nil {
		response.SetValidateError(err)
		return
	}
	token, err := service.Login(&user)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	//更新用户最后登录时间 并更新用户的GooGle验证码路径
	service.EditToken(user.Username)
	c.SetCookie(middleware.JwtName, token, setting.JwtConf.ExpTime*3600, "/", setting.ApplicationConf.Domain, false, true)
	response.SuccessData(token)
}
