package response

//返回收款地址列表
type AddressPage struct {
	Total       int       `json:"total"`        // 总共多少页
	PerPage     int       `json:"per_page"`     // 当前页码
	CurrentPage int       `json:"current_page"` // 每页显示多少条
	Data        []Address `json:"data"`
}

//收款地址对象
type Address struct {
	ID          uint    `json:"id"`           //ID
	UserId      int     `json:"user_id"`      //商户id
	RealName    string  `json:"real_name"`    //组织名称
	UserAddress string  `json:"user_address"` //地址
	EthAmount   float64 `json:"eth_amount"`   //以太币余额
	UsdtAmount  float64 `json:"usdt_amount"`  //泰达币余额
	CreateTime  int64   `json:"create_time"`  //创建时间戳
	UpdateTime  int64   `json:"update_time"`  //更新时间戳
	Status      int     `json:"status"`       //状态 1：已完成，2：执行中，3：结账中
	UseTag     int     `json:"user_tag"`     //占用状态 0：停用，1：启用

}
