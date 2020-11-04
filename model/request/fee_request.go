package request

//手续费

type FeeAdd struct {
	Amount float64 `json:"amount" binding:"required" example:"1" form:"amount"`   //实收金额金额
	UserId int     `json:"user_id" binding:"required" example:"1" form:"user_id"` //用户id
}

type FeeEdit struct {
	ID     int `json:"id" form:"id"`         //id
	Status int `json:"status" form:"status"` //状态 1：执行中 2：已完成
}
