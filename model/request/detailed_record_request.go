package request

import "time"

type DetailedAdd struct {
	IdCode        string    `json:"id_code"`        //系统编号
	Amount        float64   `json:"amount"`         //金额
	TXID          string    `json:"txid"`           //交易hash
	FromAddress   string    `json:"from_address"`   //发款地址
	ChargeAddress string    `json:"charge_address"` //收款地址
	CreateAt      time.Time `json:"create_at"`      //时间
}

type DetailedEdit struct {
	Id        int    `json:"id"`         //id
	OrderCode string `json:"order_code"` //商户订单编号
	IsBind    int    `json:"is_bind"`    // 1：解绑 2：绑定
}
