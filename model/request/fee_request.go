package request

type FeeAdd struct {
	Amount float64 `json:"amount" binding:"required" example:"1"`  //实收金额金额
	UserId int     `json:"user_id" binding:"required" example:"1"` //用户id
}

type FeeEdit struct {
	Id     int `json:"id"`     //id
	Status int `json:"status"` //状态 1：执行中 2：已完成
}
