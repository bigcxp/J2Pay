package request

//实收明细请求参数

//新增交易订单明细
type Erc20Add struct {
	TxID       string `json:"tx_id" binding:"required,max=255" example:"0xasfasf" form:"tx_id"`               //交易Id
	From       string `json:"from" binding:"required,max=255" example:"oxvfswedfvgs" form:"from"`     //打币地址
	To         string `json:"to" binding:"required,max=255" example:"oxvfswedfvgs" form:"to"`         //收币地址                             //数量
	Balance    string `json:"balance" binding:"required" example:"1234"  form:"balance"`              //金额                                            //备注
	CreateTime string `json:"create_time" binding:"required" example:"2016-01-01" form:"create_time"` // 创建时间
	Remark     string `json:"remark" binding:"" form:"remark"`                                        //备注
}

//修改交易顶单明细
type Erc20Edit struct {
	ID      int  `json:"id" form:"id"`                                                   //ID
	OrderId string `json:"order_id" binding:"required" example:"0xasfasf" form:"order_id"` //订单编号
	Status  int    `json:"status" binding:"oneof=1 2" example:"1" form:"status"`           //是否绑定订单 1 未绑定 2已绑定
}
