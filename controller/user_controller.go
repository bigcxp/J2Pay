package controller

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/model"
	"j2pay-server/pkg/util"
)


// @Tags 账户管理
// @Summary 获取账户信息
// @Produce json
// @Router /userInfo [get]
func UserInfo(c *gin.Context) {
	response := util.Response{c}
	user, hasUser := c.Get("user")
	if !hasUser {
		response.Error("用户未登录")
		return
	}
	userInfo := user.(*util.Claims)
	//创建map
	res := make(map[string]interface{}, 2)
	res["role"] = model.GetAccountRole(userInfo.ID)
	res["auth"] = model.GetAccountAuth(userInfo.ID)


	response.SuccessUserInfo(res,userInfo)
}
