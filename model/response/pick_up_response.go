package response

import (
	"time"
)

//管理端提领返回列表
type PickUpPage struct {
	Total       int        `json:"total"`        // 总共多少页
	PerPage     int        `json:"per_page"`     // 当前页码
	CurrentPage int        `json:"current_page"` // 每页显示多少条
	Data        []PickList `json:"data"`
}

//管理端代发返回列表
type SendPage struct {
	Total       int        `json:"total"`        // 总共多少页
	PerPage     int        `json:"per_page"`     // 当前页码
	CurrentPage int        `json:"current_page"` // 每页显示多少条
	Data        []SendList `json:"data"`
}

//商户端提领代发返回列表
type MerchantPickSendPage struct {
	Total       int                `json:"total"`        // 总共多少页
	PerPage     int                `json:"per_page"`     // 当前页码
	CurrentPage int                `json:"current_page"` // 每页显示多少条
	TotalAmount float64            `json:"total_amount"` // 提领总额
	TotalFee    float64            `json:"total_fee"`    //总手续费
	TotalReduce float64            `json:"total_reduce"` //总减少金额
	Data        []MerchantPickList `json:"data"`
}

//管理端提领对象
type PickList struct {
	ID          uint      `json:"id"`           //ID
	UserId      int       `json:"user_id"`      //商户id
	RealName    string    `json:"real_name"`    //组织名称
	IdCode      string    `json:"id_code"`      //系统编号
	Status      int       `json:"status"`       //状态 0：等待中，1：执行中，2：成功，3：已取消，4：失败
	SendAddress string    `json:"send_address"` //代发地址
	Amount      float64   `json:"amount"`       //金额
	TXID        string    `json:"txid"`         //交易hash
	CreateAt    time.Time `json:"create_at"`    //建立时间
	FinishTime  time.Time `json:"finish_time"`  //完成时间
	Type        int       `json:"type"`         //类型 1：提领，2：代发
}

//管理端代发对象

type SendList struct {
	ID          uint      `json:"id"`           //ID
	UserId      int       `json:"user_id"`      //商户id
	RealName    string    `json:"real_name"`    //组织名称
	IdCode      string    `json:"id_code"`      //系统编号
	OrderCode   string    `json:"order_code"`   //商户订单编号
	Status      int       `json:"status"`       //状态 0：等待中，1：执行中，2：成功，3：已取消，4：失败
	SendAddress string    `json:"send_address"` //代发地址
	Amount      float64   `json:"amount"`       //金额
	Fee         float64   `json:"fee"`          //手续费
	DelMoney    float64   `json:"del_money"`    //扣除商户余额
	TXID        string    `json:"txid"`         //交易hash
	CreateAt    time.Time `json:"create_at"`    //建立时间
	FinishTime  time.Time `json:"finish_time"`  //完成时间
	Type        int       `json:"type"`         //类型 1：提领，2：代发
	Remark      string    `json:"remark"`       //备注
}

//商户端提领和代发返回对象
type MerchantPickList struct {
	ID          uint      `json:"id"`           //ID
	UserId      int       `json:"user_id"`      //商户id
	RealName    string    `json:"real_name"`    //组织名称
	IdCode      string    `json:"id_code"`      //系统编号
	OrderCode   string    `json:"order_code"`   //商户订单编号
	Status      int       `json:"status"`       //状态 0：等待中，1：执行中，2：成功，3：已取消，4：失败
	SendAddress string    `json:"send_address"` //代发地址
	Amount      float64   `json:"amount"`       //金额
	Fee         float64   `json:"fee"`          //手续费
	GasFee      float64   `json:"fee"`          //手续费
	DelMoney    float64   `json:"del_money"`    //扣除商户余额
	TXID        string    `json:"txid"`         //交易hash
	CreateAt    time.Time `json:"create_at"`    //建立时间
	FinishTime  time.Time `json:"finish_time"`  //完成时间
	Type        int       `json:"type"`         //类型 1：提领，2：代发
	Remark      string    `json:"remark"`       //备注
}
