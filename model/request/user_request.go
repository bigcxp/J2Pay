package request

type LoginUser struct {
	Username   string `json:"username" binding:"required,max=255" example:"admin"` // 用户名
	Password   string `json:"password" binding:"required,max=255" example:"admin"` // 密码
	GoogleCode string `json:"google_code" example:"952721"`                        // google 动态验证码
}

type UserAdd struct {
	CommonUser
	Password string `json:"password" binding:"required,max=255" example:"test" form:"password"` // 密码

}

type UserEdit struct {
	ID     int    `json:"id" form:"id"`
	IsOpen int    `json:"is_open" form:"is_open" ` //是否开启双重验证 0：关闭 1：开启
	Code   string `json:"code" form:"code"`        //动态码
	CommonUser
	Password string `json:"password" example:"test" form:"password"` // 密码（非必填）

}

type Google struct {
	ID         int    `json:"id" form:"id"`
	IsOpen     int    `json:"is_open" binding:"oneof=0 1" example:"1" form:"is_open"`                     //是否开启google双重验证 默认0：不开启 1：开启
	GoogleCode string `json:"google_code" binding:"required,max=255" example:"852079" form:"google_code"` //google验证码
}

type CommonUser struct {
	UserName      string  `json:"user_name" binding:"required,max=255" example:"test" form:"user_name"`       // 账号
	RealName      string  `json:"real_name" binding:"required,max=255" example:"test" form:"real_name"`       // 组织名称
	Pid           int     `json:"pid" form:"pid" binding:"required,max=11" example:"1"`                       // 所属组织ID
	Status        int8    `json:"status" form:"status" binding:"oneof=0 1" example:"1"`                       // 状态 1：正常 0：禁用
	WhitelistIP   string  `json:"whitelist_ip" binding:"" example:"多个地址之间用逗号隔开" form:"whitelist_ip"`   //IP白名单
	Tel           string  `json:"tel" binding:"required,max=12" example:"17585534067" form:"tel"`             // 电话号码
	Address       string  `json:"address" binding:"required,max=255" example:"test" form:"address"`           // 商户地址
	ReturnUrl     string  `json:"return_url" binding:"required,max=255" example:"test" form:"return_url"`     // 回传URL
	DaiUrl        string  `json:"dai_url" binding:"required,max=255" example:"test" form:"dai_url"`           // 代发URL
	Remark        string  `json:"remark" binding:"required,max=255" example:"test" form:"remark"`             // 备注
	IsCollection  int     `json:"is_collection" form:"is_collection" binding:"oneof=0 1" example:"1"`         //是否开启收款功能 1：是 0：否
	IsCreation    int     `json:"is_creation" form:"is_creation" binding:"oneof=0 1" example:"1"`             //是否开启手动建单 1：是 0：否
	More          int     `json:"more" form:"more" binding:"required,max=11" example:"1"`                     //地址多单收款
	OrderType     int     `json:"order_type" form:"order_type" binding:"oneof=0 1" example:"1"`               //订单手续费类型 1：百分比 0：固定
	OrderCharge   float64 `json:"order_charge" form:"order_charge" binding:"required" example:"1"`            //订单手续费
	ReturnType    int     `json:"return_type" form:"return_type" binding:"oneof=0 1" example:"1"`             //退款手续费类型 1：百分比 0：固定
	ReturnCharge  float64 `json:"return_charge" form:"return_charge" binding:"required" example:"1"`          //退款手续费
	IsDai         int     `json:"is_dai" form:"is_dai" binding:"oneof=0 1" example:"1"`                       //是否启用代发功能
	DaiType       int     `json:"dai_type" form:"dai_type" binding:"oneof=0 1" example:"1"`                   //代发手续费类型 1：百分比 0：固定
	DaiCharge     float64 `json:"dai_charge" form:"dai_charge" binding:"required" example:"1"`                //代发手续费
	PickType      int     `json:"pick_type" form:"pick_type" binding:"oneof=0 1" example:"1"`                 //提领手续费类型 1：百分比 0：固定
	PickCharge    float64 `json:"pick_charge" form:"pick_charge" binding:"required" example:"1"`              //提领手续费
	IsGas         int     `json:"is_gas" form:"is_gas" binding:"oneof=0 1" example:"1"`                       //是否启用gas预估 1：是 0：否
	Examine       float64 `json:"examine" form:"examine" binding:"required" example:"1"`                      //代发审核
	DayTotalCount float64 `json:"day_total_count" form:"day_total_count" binding:"required" example:"1"`      //每日交易总量
	MaxOrderCount float64 `json:"max_order_count" form:"max_order_count" binding:"required" example:"1"`      //最大交易总量
	MinOrderCount float64 `json:"min_order_count" form:"min_order_count" binding:"required" example:"1"`      //最小交易总量
	Limit         float64 `json:"limit" form:"limit" binding:"required" example:"1"`                          //结账限制
	UserLessTime  int64   `json:"user_less_time" form:"user_less_time" binding:"required,max=11" example:"1"` //订单无效时间
	Roles         []int   `json:"roles" form:"roles" binding:"required,min=1" example:"1,2"`                  // 所属角色
}
