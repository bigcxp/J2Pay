package request

//登录所需实体
type LoginUser struct {
	Username   string `json:"username" binding:"required,max=255" example:"admin"` // 用户名
	Password   string `json:"password" binding:"required,max=255" example:"admin"` // 密码
	GoogleCode string `json:"google_code" example:"952721"`                        // google 动态验证码
}

//新增账户
type AccountAdd struct {
	CommonAccount
	Password   string `json:"password" binding:"required,max=255" example:"test" form:"password"`        // 密码
	RePassword string `json:"re_password" binding:"required,max=255" example:"admin" form:"re_password"` // 确认密码

}

//编辑账户
type AccountEdit struct {
	ID     int64  `json:"id" form:"id"`
	Status int    `json:"status"`                  //是否启用 1：开启 2：关闭
	Token  string `json:"token"`                   //token
	IsOpen int    `json:"is_open" form:"is_open" ` //是否开启双重验证 1：开启 2：关闭
	RID    int    `json:"rid" form:"rid" binding:"required" example:"1"`
}

//开启google验证
type Google struct {
	ID         int    `json:"id" form:"id"`
	IsOpen     int    `json:"is_open" binding:"oneof=0 1" example:"1" form:"is_open"`                     //是否开启双重验证 1：开启 2：关闭
	GoogleCode string `json:"google_code" binding:"required,max=255" example:"852079" form:"google_code"` //google验证码
	Code       string `json:"code" form:"code"`                                                           //动态码
}

//账户
type CommonAccount struct {
	UID      int64  `json:"uid" example:"1" form:"uid"`                                           //所属组织名称
	UserName string `json:"user_name" binding:"required,max=255" example:"test" form:"user_name"` // 使用者名称
	RID      int    `json:"rid" form:"rid" example:"1"`                                           // 所属角色
}
