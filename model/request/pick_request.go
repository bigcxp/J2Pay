package request

//提领
type PickAdd struct {
	Amount      float64 `json:"amount" binding:"required" example:"1"`                             //数量
	Type        int     `json:"type" binding:"required" example:"1"`                               //类型 1：代发 0：提领
	UserId      int     `json:"user_id" binding:"required" example:"1"`                            //用户id
	Remark      string  `json:"remark" binding:"required,max=255" example:"备注"`                    //备注
	Currency    string  `json:"currency";binding:"oneof=RMB TWD 1" example:"RMB"`                  //换算汇率
	PickAddress string  `json:"pick_address" binding:"required,max=255" example:"0x1243cfsadfcsd"` //收款地址
}

//代发
type SendAdd struct {
	OrderCode   string  `json:"order_code" binding:"required,max=255" example:"asfasgdsasfgas"`    //商户订单编号
	Amount      float64 `json:"amount" binding:"required" example:"1"`                             //数量
	Type        int     `json:"type" binding:"required" example:"1"`                               //类型 1：代发 0：提领
	UserId      int     `json:"user_id" binding:"required" example:"1"`                            //用户id
	Remark      string  `json:"remark" binding:"required,max=255" example:"备注"`                    //备注
	Currency    string  `json:"currency";binding:"oneof=RMB TWD 1" example:"RMB"`                  //换算汇率
	PickAddress string  `json:"pick_address" binding:"required,max=255" example:"0x1243cfsadfcsd"` //收款地址
}

//取消提领 代发
type SendEdit struct {
	ID     int `json:"id"`
	Status int `json:"status" binding:"required" example:"3"` // 状态 0：等待中，1：执行中，2：成功，3：已取消，4：失败
}
