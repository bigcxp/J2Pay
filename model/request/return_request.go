package request

type ReturnAdd struct {
	OrderCode   string  `binding:"required,max=255" example:"asfasgdsasfgas"` //商户订单编号
	Amount      float64 `json:"amount" binding:"required" example:"1"` //实收金额
	UserId      int     `json:"user_id" binding:"required" example:"1"` //用户id
}
