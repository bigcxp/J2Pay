package request

type ParameterEdit struct {
	ID           int     `json:"id" form:"id"`
	Confirmation int     `json:"confirmation" example:"12" form:"confirmation"`  // 交易确认数
	GasLimit     int     `json:"gas_limit"  example:"1" form:"gas_limit"`        //gas Limit
	GasPrice     float64 `json:"gas_price"  example:"88.00000" form:"gas_price"` // GasPrice
}
