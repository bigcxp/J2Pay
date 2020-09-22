package model

type PickUpPage struct {
	Total       int     `json:"total"`        // 总共多少页
	PerPage     int     `json:"per_page"`     // 当前页码
	CurrentPage int     `json:"current_page"` // 每页显示多少条
	TotalAmount float64 `json:"total_amount"` // 提领总额
	TotalFee    float64 `json:"total_fee"`    //总手续费
	TotalReduce float64 `json:"total_reduce"` //总减少金额
	Data        []Pick  `json:"data"`
}
