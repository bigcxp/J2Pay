package response

import (
	"time"
)

type DetailedRecordPage struct {
	Total       int            `json:"total"`        // 总共多少页
	PerPage     int            `json:"per_page"`     // 当前页码
	CurrentPage int            `json:"current_page"` // 每页显示多少条
	Data        []DetailedList `json:"data"`         //数据
}

type DetailedList struct {
	ID            uint      `json:"id"`             //ID
	IdCode        string    `json:"id_code"`        //系统编号
	OrderCode     string    `json:"order_code"`     //订单编号
	RealName      string    `json:"real_name"`      //组织名称
	Amount        float64   `json:"amount"`         //金额
	TXID          string    `json:"txid"`           //交易hash
	Remark        string    `json:"remark"`         //备注
	ChargeAddress string    `json:"charge_address"` //收款地址
	Status        int       `json:"status"`         //状态 1：未绑定 2：已绑定
	CreateAt      time.Time `json:"create_at"`      //时间
	OrderId       int       `json:"order_id"`       // 商户订单id
	UserId        int       `json:"user_id"`        //商户id
}
