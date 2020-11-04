package request

//新增用户收钱地址
type AddressAdd struct {
	Num          int64 `json:"num" form:"num"`                                  // 生成钱包数量
	UserId       int   `json:"user_id" example:"1" form:"user_id"`              //组织id
	HandleStatus int   `json:"handle_status" example:"1"  form:"handle_status"` //指派状态 0：所有，1：启用，2：停用
	UseTag       int64 `json:"use_tag" example:"0" form:"use_tag"`              //-2:eth钱包,-1:热钱包,0:未占用,userId:组织充币钱包
}

//编辑收款地址
type AddressEdit struct {
	ID     []int `json:"id" binding:"required,unique,min=1" example:"1,2,3,4,5,6" form:"id"`
	UserId int   `json:"user_id" binding:"required",example:"1" form:"user_id"` //组织id
}

//删除收款地址
type AddressDel struct {
	ID []int `json:"id" binding:"required,unique,min=1" example:"1,2,3,4,5,6" form:"id"`
}

//启用停用
type OpenOrStopAddress struct {
	ID           []int `json:"id" binding:"required,unique,min=1" example:"1,2,3,4,5,6" form:"id"`
	HandleStatus int   `json:"handle_status" example:"1" form:"handle_status"` //是否启用
}

//储值
type SaveAmount struct {
	ID        []int   `json:"id" binding:"required,unique,min=1" example:"1,2,3,4,5,6" form:"id"`
	EthAmount float64 `json:"eth_amount" binding:"required" example:"1" form:"eth_amount"` //以太币数量
}

//结账
type Math struct {
	ID     []int `json:"id" binding:"required,unique,min=1" example:"1,2,3,4,5,6" form:"id"`
	Status int   `json:"status" binding:"required" example:"1：已完成，2：执行中，3：结账中" form:"status"` //状态
}

//更新余额
type UpdateAmount struct {
	ID []int `json:"id" binding:"required,unique,min=1" example:"1,2,3,4,5,6" form:"id"`
}
