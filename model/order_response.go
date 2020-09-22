package model

type OrderPage struct {
	Total          int     `json:"total"`           // 总共多少页
	PerPage        int     `json:"per_page"`        // 当前页码
	CurrentPage    int     `json:"current_page"`    // 每页显示多少条
	TotalAmount    float64 `json:"total_amount"`    //总订单金额
	MerchantAmount float64 `json:"merchant_amount"` //总商户总实收金额
	ReallyAmount   float64 `json:"really_amount"`   //总实收金额
	TotalFee       float64 `json:"total_fee"`       //总手续费
	Data []Order `json:"data"`
}
