package request

type OrderAdd struct {
	OrderCode   string  `binding:"required,max=255" example:"asfasgdsasfgas"` //商户订单编号
	Amount      float64 `json:"amount" binding:"required" example:"1"` //数量
	UserId      int     `json:"user_id" binding:"required" example:"1"` //用户id
	Remark      string  `binding:"required,max=255" example:"备注"` //备注
}