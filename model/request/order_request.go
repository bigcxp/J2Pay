package request

//新增订单
type OrderAdd struct {
	Uts       int64   `json:"uts" binding:"required" example:"1231244520" form:"uts"`                           //时间戳
	UID       int64   `json:"uid" binding:"" form:"uid"`                                                        //组织ID
	OrderCode string  `json:"order_code" binding:"required,max=255" example:"asfasgdsasfgas" form:"order_code"` //商户订单编号
	Amount    float64 `json:"amount" binding:"required" example:"1" form:"amount"`                              //数量
	Remark    string  `example:"备注" json:"remark" form:"remark"`                                                //备注
}

//编辑订单
type OrderEdit struct {
	ID           int     `json:"id" form:"id"`
	Status       int     `json:"status" binding:"required" example:"3" form:"status"`               // 状态 -1：收款中，1：已完成，2：异常，3：退款等待中，4：退款中，5：退款失败，6：已退款，7：：已过期
	Address      string  `json:"address" binding:"required,max=255" form:"address"`                 //收款地址
	ShouldAmount float64 `json:"should_amount" binding:"required" example:"1" form:"should_amount"` //应收金额
	ExprireTime  int64   `json:"exprire_time" form:"exprire_time"`                                  //过期时间

}
