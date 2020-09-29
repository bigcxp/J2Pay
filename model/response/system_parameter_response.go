package response

import "time"

//返回系统参数对象
type Parameter struct {
	ID           uint      `json:"id"`           //ID
	Confirmation int       `json:"confirmation"` // 交易确认数
	GasLimit     int       `json:"gas_limit"`    //gas Limit
	GasPrice     float64   `json:"gas_price"`    // GasPrice
	EthFee       float64   `json:"eth_fee"`      //ETH 最小矿工费
	UpdatedAt    time.Time `json:"updated_at"`   //更新时间
}
