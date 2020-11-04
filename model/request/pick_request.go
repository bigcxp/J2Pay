package request

//提领
type PickAdd struct {
	Amount      float64 `json:"amount" binding:"required" example:"1" form:"amount"`                                   //数量
	Type        int     `json:"type" binding:"required" example:"1" form:"type"`                                       //类型 1：代发 0：提领
	UserId      int     `json:"user_id" binding:"required" example:"1"  form:"user_id"`                                //用户id
	Remark      string  `json:"remark" binding:"required,max=255" example:"备注" form:"remark"`                          //备注
	Currency    string  `json:"currency";binding:"oneof=RMB TWD 1" example:"RMB" form:"currency"`                      //换算汇率
	PickAddress string  `json:"pick_address" binding:"required,max=255" example:"0x1243cfsadfcsd" form:"pick_address"` //收款地址
}

//代发
type SendAdd struct {
	OrderCode   string  `json:"order_code" binding:"required,max=255" example:"asfasgdsasfgas" form:"order_code"`      //商户订单编号
	Amount      float64 `json:"amount" binding:"required" example:"1" form:"amount"`                                   //数量
	Type        int     `json:"type" binding:"required" example:"1" form:"type"`                                       //类型 1：代发 0：提领
	UserId      int     `json:"user_id" binding:"required" example:"1" form:"user_id"`                                 //用户id
	Remark      string  `json:"remark" binding:"required,max=255" example:"备注" form:"remark"`                          //备注
	Currency    string  `json:"currency";binding:"oneof=RMB TWD 1" example:"RMB" form:"currency"`                      //换算汇率
	PickAddress string  `json:"pick_address" binding:"required,max=255" example:"0x1243cfsadfcsd" form:"pick_address"` //收款地址
}

//取消提领 代发
type SendEdit struct {
	ID     int `json:"id" form:"id"`
	Status int `json:"status" binding:"required" example:"3" form:"status"` // 状态 0：等待中，1：执行中，2：成功，3：已取消，4：失败
}
