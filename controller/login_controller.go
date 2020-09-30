package controller

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/middleware"
	"j2pay-server/model/request"
	"j2pay-server/pkg/logger"
	"j2pay-server/pkg/setting"
	"j2pay-server/pkg/util"
	"j2pay-server/service"
)

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

	id, ok := c.Get("session_id")
	if !ok {
		response.ErrorMsg("GET SESSION ERROR")
		return
	}
	idStr := id.(string)
	token, err := service.Login(&user, idStr)
	if err != nil {
		response.SetOtherError(err)
		return
	}
	//更新用户信息
	service.EditToken(token, user.Username)
	c.SetCookie(middleware.JwtName, token, setting.JwtConf.ExpTime*3600, "/", setting.ApplicationConf.Domain, false, true)
	response.SuccessData(token)
}

// @Tags 登录操作
// @Summary 验证码
// @Success 200
// @Router /captcha [get]
func Captcha(c *gin.Context) {
	id, ok := c.Get("session_id")
	if !ok {
		logger.Logger.Error("获取 session 异常")
		return
	}
	idStr := id.(string)
	image, err := service.NewCaptcha(idStr)
	if err != nil {
		logger.Logger.Error("创建验证码异常：", err)
		return
	}
	c.Header("content-type", "image/png")
	_, err = image.WriteTo(c.Writer)
	if err != nil {
		logger.Logger.Error("c.Writer.Write() 异常：", err)
	}
}
