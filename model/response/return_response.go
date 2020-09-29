package response

import (
	"time"
)

//退款订单列表
type ReturnPage struct {
	Total       int          `json:"total"`        // 总共多少页
	PerPage     int          `json:"per_page"`     // 当前页码
	CurrentPage int          `json:"current_page"` // 每页显示多少条
	Data        []ReturnList `json:"data"`
}

type ReturnList struct {
	ID          uint      `json:"id"`           //ID
	UserId      int       `json:"user_id"`      //商户id
	SystemCode string    `json:"system_code"` //系统订单编号
	OrderCode  string    `json:"order_code"`  //商户订单编号
	RealName   string    `json:"real_name"`   //组织名称
	Amount     float64   `json:"amount"`      //实收金额
	FinishTime time.Time `json:"finishTime"`  //时间
	Status     int       `json:"status"`      //状态 1：退款等待中，2：退款中，3：退款失败，4：已退款
}
