package request

//创建eth交易
type EthTxAdd struct {
	To      string  `json:"to" binding:"required,max=255" example:"0xabcd"`
	Balance float64 `json:"balance" binding:"required" example:"20" form:"balance"`
}

//编辑eth交易
type EthTxEdit struct {
	ID      int64   `json:"id" form:"id"`
	Balance float64 `json:"balance" binding:"required" example:"20" form:"balance"`
}
