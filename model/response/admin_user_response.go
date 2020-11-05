package response

import "time"

type AdminUserPage struct {
	Total       int             `json:"total"`        // 总共多少页
	PerPage     int             `json:"per_page"`     // 当前页码
	CurrentPage int             `json:"current_page"` // 每页显示多少条
	Data        []AdminUserList `json:"data"`
}

//主账户
type Account struct {
	Token    string `json:"token"`     //Token
}

//组织详情
type AdminUserList struct {
	UserName string `json:"user_name"` //主账号
	Account       Account `json:"account"`          // 主账户
	RealName      string  `json:"real_name"`        // 组织名称
	WhitelistIP   string  `json:"whitelist_ip" `    //IP白名单
	Address       string  `json:"address"`          // 商户地址
	Balance       float64 `json:"balance" `         //余额
	ReturnUrl     string  `json:"return_url"`       // 回传URL
	DaiUrl        string  `json:"dai_url" `         // 代发URL
	Remark        string  `json:"remark" `          // 备注
	IsCollection  int     `json:"is_collection" `   //是否开启收款功能 1：是 0：否
	IsCreation    int     `json:"is_creation" `     //是否开启手动建单 1：是 0：否
	More          int     `json:"more" `            //地址多单收款
	OrderType     int     `json:"order_type" `      //订单手续费类型 1：百分比 0：固定
	OrderCharge   float64 `json:"order_charge"`     //订单手续费
	ReturnType    int     `json:"return_type" `     //退款手续费类型 1：百分比 0：固定
	ReturnCharge  float64 `json:"return_charge" `   //退款手续费
	IsDai         int     `json:"is_dai" `          //是否启用代发功能
	DaiType       int     `json:"dai_type" `        //代发手续费类型 1：百分比 0：固定
	DaiCharge     float64 `json:"dai_charge"`       //代发手续费
	PickType      int     `json:"pick_type" `       //提领手续费类型 1：百分比 0：固定
	PickCharge    float64 `json:"pick_charge" `     //提领手续费
	IsGas         int     `json:"is_gas" `          //是否启用gas预估 1：是 0：否
	Examine       float64 `json:"examine" `         //代发审核
	DayTotalCount float64 `json:"day_total_count" ` //每日交易总量
	MaxOrderCount float64 `json:"max_order_count" ` //最大交易总量
	MinOrderCount float64 `json:"min_order_count"`  //最小交易总量
	Limit         float64 `json:"limit" `           //结账限制
	UserLessTime  int64   `json:"user_less_time" `  //订单无效时间

}

type AdminUserMessagePage struct {
	Total       int                    `json:"total"`        // 总共多少页
	PerPage     int                    `json:"per_page"`     // 当前页码
	CurrentPage int                    `json:"current_page"` // 每页显示多少条
	Data        []AdminUserMessageList `json:"data"`
}

type AdminSystemMessage struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`      //标题
	BeginTime time.Time `json:"begin_time"` //开始时间
	EndTime   time.Time `json:"end_time"`   //结束时间
}

type AdminUserMessageList struct {
	Id             int                  `json:"id"`
	UserName       string               `json:"user_name"`      //用户名
	SystemMessages []AdminSystemMessage `json:"systemMessages"` //系统公告
}
