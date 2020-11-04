package request

type DetailedAdd struct {
	IdCode        string  `json:"id_code" binding:"required,max=255" example:"test" form:"id_code"`                //系统编号
	Amount        float64 `json:"amount"  binding:"required" example:"1" form:"amount"`                            //金额
	TXID          string  `json:"txid"  form:"txid"`                                                               //交易hash
	FromAddress   string  `json:"from_address"  binding:"required,max=255" example:"test" form:"from_address"`     //发款地址
	ChargeAddress string  `json:"charge_address"  binding:"required,max=255" example:"test" form:"charge_address"` //收款地址
}

type DetailedEdit struct {
	ID        int    `json:"id" form:"id"`                                                           //id
	OrderCode string `json:"order_code" binding:"required,max=255" example:"test" form:"order_code"` //商户订单编号
	IsBind    int    `json:"is_bind" form:"is_bind"`                                                 // 1：解绑 2：绑定
}
