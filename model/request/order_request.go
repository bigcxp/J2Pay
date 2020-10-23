package request

import "time"

type OrderAdd struct {
	OrderCode string  `binding:"required,max=255" example:"asfasgdsasfgas"` //商户订单编号
	Amount    float64 `json:"amount" binding:"required" example:"1"`        //数量
	UserId    int     `json:"user_id" binding:"required" example:"1"`       //用户id
	Remark    string  `binding:"required,max=255" example:"备注"`             //备注
}

type OrderEdit struct {
	ID           int       `json:"id"`
	Status       int       `json:"status" binding:"required" example:"3"`        // 状态 -1：收款中，1：已完成，2：异常，3：退款等待中，4：退款中，5：退款失败，6：已退款，7：：已过期
	Address      string    `json:"address" binding:"required,max=255"`           //收款地址
	ShouldAmount float64   `json:"should_amount" binding:"required" example:"1"` //应收金额
	ExprireTime  time.Time `json:"exprire_time"`                                 //过期时间

}
