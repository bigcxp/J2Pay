package request

type ReturnAdd struct {
	OrderCode string  `json:"order_code" binding:"required,max=255" example:"asfasgdsasfgas" form:"order_code"` //商户订单编号
	Amount    float64 `json:"amount" binding:"required" example:"1" form:"amount"`                              //实收金额
	UserId    int     `json:"user_id" binding:"required" example:"1" form:"user_id"`                            //用户id
}
