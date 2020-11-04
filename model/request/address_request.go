package request

//新增用户收钱地址
type AddressAdd struct {
	Num          int64 `json:"num"`                       // 生成钱包数量
	UserId       int   `json:"user_id" example:"1"`       //组织id
	HandleStatus int   `json:"handle_status" example:"1"` //指派状态 0：所有，1：启用，2：停用
	UseTag       int64 `json:"use_tag" example:"0"`       //-2:eth钱包,-1:热钱包,0:未占用,userId:组织充币钱包
}

//编辑收款地址
type AddressEdit struct {
	ID     []int `json:"id" binding:"required,unique,min=1" example:"1,2,3,4,5,6"`
	UserId int   `json:"user_id" binding:"required",example:"1"` //组织id
}

//删除收款地址
type AddressDel struct {
	ID []int `json:"id" binding:"required,unique,min=1" example:"1,2,3,4,5,6"`
}

//启用停用
type OpenOrStopAddress struct {
	ID           []int `json:"id" binding:"required,unique,min=1" example:"1,2,3,4,5,6"`
	HandleStatus int   `json:"handle_status" example:"1"` //是否启用
}

//储值
type SaveAmount struct {
	ID        []int   `json:"id" binding:"required,unique,min=1" example:"1,2,3,4,5,6"`
	EthAmount float64 `json:"eth_amount" binding:"required" example:"1"` //以太币数量
}

//结账
type Math struct {
	ID     []int `json:"id" binding:"required,unique,min=1" example:"1,2,3,4,5,6"`
	Status int   `json:"status" binding:"required" example:"1：已完成，2：执行中，3：结账中"` //状态
}

//更新余额
type UpdateAmount struct {
	ID []int `json:"id" binding:"required,unique,min=1" example:"1,2,3,4,5,6"`
}
