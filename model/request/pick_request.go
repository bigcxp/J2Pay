package request

type PickAdd struct {
	OrderCode   string  `binding:"required,max=255" example:"asfasgdsasfgas"` //商户订单编号
	Amount      float64 `json:"amount" binding:"required" example:"1"` //数量
	Type        int     `json:"type" binding:"required" example:"1"` //类型 1：代发 0：收款
	UserId      int     `json:"user_id" binding:"required" example:"1"` //用户id
	Remark      string  `binding:"required,max=255" example:"备注"` //备注
	PickAddress string  `binding:"required,max=255" example:"0x1243cfsadfcsd"` //收款地址
}
