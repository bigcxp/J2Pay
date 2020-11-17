package response

import "time"

//实收明细返回参数
//商户实收订单列表
type Erc20Page struct {
	Total       int         `json:"total"`        // 总共多少页
	PerPage     int         `json:"per_page"`     // 当前页码
	CurrentPage int         `json:"current_page"` // 每页显示多少条
	Data        []Erc20List `json:"data"`
}

//商户实收记录订单对象
type Erc20List struct {
	ID          int       `json:"id"`           //ID
	RealName    string    `json:"real_name"`    //组织名称
	OrderId     string    `json:"order_id"`     //商户订单编号
	ToAddress   string    `json:"to_address"`   //收款地址
	BalanceReal string    `json:"balance_real"` //金额
	Status      int       `json:"status"`       //状态 1：未绑定，2：已绑定
	TxID        string    `json:"tx_id"`        // 交易id
	Remark      string    `json:"remark"`       //备注
	CreateTime  int64     `json:"create_time"`  //创建时间戳
	Create      time.Time `json:"create"`       //创建时间
}
