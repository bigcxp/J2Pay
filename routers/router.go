package routers

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"j2pay-server/controller"
	_ "j2pay-server/docs"
	"j2pay-server/middleware"
	"j2pay-server/pkg/setting"
)

func InitRouter() *gin.Engine {
	gin.SetMode(setting.ApplicationConf.Env)
	r := gin.New()
	// swagger 文档输出
	if setting.ApplicationConf.Env == "debug" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// 加入通用中间件
	r.Use(
		gin.Recovery(),           // recovery 防止程序奔溃
		middleware.Logger(),      // 日志记录
		middleware.NoFound(),     // 404 处理
		middleware.MakeSession(), // session支持
		middleware.ErrorHandle(), // 错误处理
	)

	r.GET("/captcha", controller.Captcha)
	r.POST("/login", controller.Login)

	// 加入鉴权中间件
	r.Use(middleware.JWT())
	r.GET("/userInfo", controller.UserInfo)
	r.Use(middleware.Authentication())
	// 用户
	{
		r.GET("/auth/role", controller.RoleTree)
		r.GET("/adminUser", controller.UserIndex)
		r.GET("/adminUser/:id", controller.UserDetail)
		r.POST("/adminUser", controller.UserAdd)
		r.PUT("/adminUser/:id", controller.UserEdit)
		r.DELETE("/adminUser/:id", controller.UserDel)
	}

	// 角色
	{
		r.GET("/auth/tree", controller.AuthTree)
		r.GET("/role", controller.RoleIndex)
		r.GET("/role/:id", controller.RoleDetail)
		r.POST("/role", controller.RoleAdd)
		r.PUT("/role/:id", controller.RoleEdit)
		r.DELETE("/role/:id", controller.RoleDel)
	}
	//公告
	{
		r.GET("/systemMessage", controller.SystemMessage)
		r.GET("/systemMessageByUser", controller.SystemMessageByUserId)
		r.POST("/systemMessage", controller.SystemMessageAdd)
		r.DELETE("/systemMessage/:id", controller.SystemMessageDel)
		r.PUT("/systemMessage/:id", controller.SystemMessageEdit)

	}
	//首页
	{
		r.GET("/index", controller.IndexSystem)
	}
	//商户提领
	{
		r.GET("/merchantPick",controller.PickIndex)
		r.POST("/merchantPick",controller.PickAdd)
		r.GET("/merchantPick/:id",controller.PickDetail)
		r.POST("/notify",controller.PickNotify)
	}
	//订单
	{
		r.GET("/order",controller.OrderList)
		r.POST("/order",controller.OrderAdd)
		r.GET("/order/:id",controller.OrderDetail)
	}
	//订单退款
	{
		r.GET("/return",controller.ReturnList)
		r.POST("/return",controller.ReturnAdd)
		r.GET("/return/:id",controller.ReturnDetail)
	}
	//手续费结账
	{
		r.GET("/fee",controller.FeeList)
		r.POST("/fee",controller.Settle)
	}
	return r
}
