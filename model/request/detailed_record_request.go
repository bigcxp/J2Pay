package request

type DetailedAdd struct {
	IdCode        string  `json:"id_code" binding:"required,max=255" example:"test"`         //系统编号
	Amount        float64 `json:"amount"  binding:"required" example:"1"`                    //金额
	TXID          string  `json:"txid"`                                                      //交易hash
	FromAddress   string  `json:"from_address"  binding:"required,max=255" example:"test"`   //发款地址
	ChargeAddress string  `json:"charge_address"  binding:"required,max=255" example:"test"` //收款地址
	CreateAt      int64   `json:"create_at"`                                                 //时间
}

type DetailedEdit struct {
	ID        int    `json:"id"`                                                   //id
	OrderCode string `json:"order_code" binding:"required,max=255" example:"test"` //商户订单编号
	IsBind    int    `json:"is_bind"`                                              // 1：解绑 2：绑定
}
