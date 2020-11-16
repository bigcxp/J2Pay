package request

//提领
type WithDrawAdd struct {
	RealName string `json:"real_name"`                  //商户
	Symbol   string `json:"symbol" binding:"required"`  //币种
	Balance  string `json:"balance" binding:"required"` //金额
	Remark   string `json:"remark" binding:""`          //备注
}

//代发
type SendAdd struct {
	OrderCode string `json:"order_code" binding:"required,max=255" example:"asfasgdsasfgas" form:"order_code"` //商户订单编号
	RealName  string `json:"real_name"`                                                                        //商户名称
	Symbol    string `json:"symbol" binding:"required"`                                                        //币种
	Address   string `json:"address" binding:"required"`                                                       //收款地址
	Balance   string `json:"balance" binding:"required"`                                                       //金额
	Remark    string `json:"remark" binding:""`                                                                //备注
}

//取消提领 代发
type SendEdit struct {
	ID     int   `json:"id" form:"id"`
	Status int64 `json:"status" binding:"required" example:"3" form:"status"` // 状态 0：等待中，1：执行中，2：成功，3：已取消，4：失败
}
