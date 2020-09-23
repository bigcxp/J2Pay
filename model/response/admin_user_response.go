package response

import (
	"time"
)

type AdminUserPage struct {
	Total       int             `json:"total"`        // 总共多少页
	PerPage     int             `json:"per_page"`     // 当前页码
	CurrentPage int             `json:"current_page"` // 每页显示多少条
	Data        []AdminUserList `json:"data"`
}

type CasRole struct {
	Id   int    `json:"id"`   // 角色ID
	Name string `json:"name"` // 角色名
}

type AdminUserList struct {
	Id            int       `json:"id"`
	UserName      string    `json:"user_name"`       // 登录名
	Tel           string    `json:"tel"`             // 手机号码
	RealName      string    `json:"real_name"`       // 组织名称
	Balance       float64   `json:"balance"`         //账户余额
	Address       string    `json:"address"`         //商户地址
	OrderCharge   float64   `json:"order_charge"`    //提领fee
	ReturnCharge  float64   `json:"return_charge"`   //退款fee
	Remark        string    `json:"remark"`          //备注
	Status        int8      `json:"status"`          // 用户状态
	LastLoginTime time.Time `json:"last_login_time"` //最后登录时间
	Roles         []CasRole `json:"roles"`           // 角色信息
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

