package request

//编辑收款地址
type AddressEdit struct {
	Id     int `json:"-"`
	UserId int `json:"user_id" binding:"required" example:"1"` //组织id
}

//启用停用
type OpenOrStopAddress struct {
	Id         int `json:"-"`
	UseTag int64 `json:"hand_status" binding:"oneof=0 1" example:"0：停用，1：启用"` //是否启用
}

//储值
type SaveAmount struct {
	Id        int     `json:"-"`
	EthAmount float64 `json:"eth_amount" binding:"required" example:"1"` //以太币数量
}

//结账
type Math struct {
	Id     int `json:"-"`
	Status int `json:"status" binding:"required" example:"1：已完成，2：执行中，3：结账中"` //状态
}

//更新余额
type UpdateAmount struct {
	Id[]     int `json:"-"`
}
