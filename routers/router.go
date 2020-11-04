package routers

import (
	"j2pay-server/controller"
	_ "j2pay-server/docs"
	"j2pay-server/middleware"
	"j2pay-server/pkg/setting"
	"net/http"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter() *gin.Engine {
	gin.SetMode(setting.ApplicationConf.Env)
	r := gin.New()
	// swagger 文档输出
	if setting.ApplicationConf.Env == "debug" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	//静态文件处理
	r.LoadHTMLGlob("views/**/*")
	r.StaticFS("/static", http.Dir("./static"))

	// 加入通用中间件
	r.Use(
		//gin.Recovery(),           // recovery 防止程序奔溃
		middleware.Logger(),      // 日志记录
		middleware.NoFound(),     // 404 处理
		middleware.MakeSession(), // session支持
		middleware.ErrorHandle(), // 错误处理
	)
	//不需要授权
	r.POST("/login", controller.Login)
	r.GET("/login", controller.LoginIndex)
	r.GET("/index", controller.Index)
	r.GET("/", controller.Indexs)
	r.GET("adminUserIndex", controller.AdminUserIndex)
	r.GET("/roleIndex", controller.IndexRole)
	r.GET("/systemMessageIndex", controller.SystemMessageIndex)
	r.GET("/systemIndex", controller.SystemIndex)
	r.GET("/pickIndex", controller.IndexPick)
	r.GET("rateIndex", controller.RateIndex)
	r.GET("/parameterIndex", controller.ParameterIndex)
	r.GET("/detailedRecordIndex", controller.DetailedRecordIndex)
	r.GET("/returnIndex", controller.ReturnIndex)
	r.GET("/orderIndex", controller.OrderIndex)
	r.GET("/main", controller.MainIndex)
	//加入签名中间件
	//r.Use(middleware.SetUp())
	//创建新订单（充币）
	r.POST("/order", controller.OrderAdd)
	//提领,代发
	r.POST("/merchantPick", controller.PickAdd)
	r.POST("/merchantSend", controller.SendAdd)

	//加入jwt中间件
	r.Use(middleware.JWT())
	//登录后能做的操作
	//代发 提领 eth转账 erc20代币转账 生成用户钱包地址
	//1.生成钱包地址 =》热钱包地址 eth钱包地址 分配商户充币地址
	r.POST("/createAddress", controller.CreateAddress)
	//2.获取钱包列表
	r.GET("/addrList", controller.AddrList)
	//3.启用 停用地址
	r.POST("addrRestart", controller.AddrRestart)
	//4.更新余额
	r.POST("updateBalance", controller.UpdateBalance)
	//5.编辑地址
	r.POST("addrEdit", controller.AddrEdit)
	//6.删除地址
	r.DELETE("addrDel", controller.AddrDel)
	//7.交易记录
	r.GET("ethTransfer", controller.EthTransfer)
	r.GET("hotTransfer", controller.HotTransfer)
	r.GET("/userInfo", controller.UserInfo)
	// 加入鉴权中间件
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
		r.GET("/auth/list", controller.AuthList)
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
		r.PUT("/password/:id", controller.UpdatePassword)
		r.PUT("/google/:id", controller.GoogleValidate)
	}
	//商户提领 代发
	{
		r.GET("/merchantPick", controller.MerchantPickIndex)
		r.GET("/merchantPick/:id", controller.MerchantPickDetail)

		r.GET("/pick", controller.PickIndex)
		r.GET("/pick/:id", controller.PickDetail)

		r.GET("/send", controller.SendIndex)
		r.GET("/send/:id", controller.SendDetail)

		r.PUT("/pick/:id", controller.PickEdit)

		r.POST("/notify", controller.PickNotify)
	}
	//订单
	{
		r.GET("/order", controller.OrderList)
		r.GET("/order/:id", controller.OrderDetail)
		r.PUT("/order/:id", controller.OrderEdit)
		r.POST("/orderNotify", controller.OrderNotify)
	}
	//订单退款
	{
		r.GET("/return", controller.ReturnList)
		r.POST("/return", controller.ReturnAdd)
		r.GET("/return/:id", controller.ReturnDetail)
	}
	//手续费结账
	{
		r.GET("/feeIndex", controller.FeeIndex)
		r.GET("/fee", controller.FeeList)
		r.POST("/fee", controller.Settle)
	}
	//实收明细记录
	{
		r.GET("/detail", controller.DetailedList)
		r.GET("/detail/:id", controller.DetailedDetail)
		r.POST("/detail", controller.DetailedAdd)
		r.PUT("/detail/:id", controller.DetailedEdit)
	}
	//系统参数管理
	{

		r.GET("/system", controller.SystemParameter)
		r.PUT("/system/:id", controller.SystemParameterEdit)
		r.PUT("/systemGasPrice/:id", controller.SystemGasPriceEdit)
	}
	//汇率管理
	{
		r.GET("/rate", controller.RateList)
		r.GET("/rate/:id", controller.RateDetail)
		r.PUT("/rate/:id", controller.RateEdit)

	}
	return r
}
