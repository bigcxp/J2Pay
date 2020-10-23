package request

import "time"

type ParameterEdit struct {
	ID           int       `json:"id"`
	Confirmation int       `json:"confirmation" example:"12"`     // 交易确认数
	GasLimit     int       `json:"gas_limit"  example:"1"`        //gas Limit
	GasPrice     float64   `json:"gas_price"  example:"88.00000"` // GasPrice
	UpdateAt     time.Time `json:"create_at"`                     //更新时间
}
