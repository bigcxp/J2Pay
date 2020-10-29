package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Tags 首页
// @Summary 首页
// @Produce json
// @Router /main [get]
func MainIndex(c *gin.Context) {
	c.HTML(200,"main.html" ,gin.H{
		"code": http.StatusOK,
	})
}
