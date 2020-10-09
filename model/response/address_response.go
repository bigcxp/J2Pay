package response

import "time"

//返回收款地址列表
type AddressPage struct {
	Total       int             `json:"total"`        // 总共多少页
	PerPage     int             `json:"per_page"`     // 当前页码
	CurrentPage int             `json:"current_page"` // 每页显示多少条
	Data        []AddressList `json:"data"`
}

//收款地址对象
type AddressList struct {
	ID          uint      `json:"id"`           //ID
	UserId      int       `json:"user_id"`      //商户id
	RealName    string    `json:"real_name"`    //组织名称
	UserAddress string    `json:"user_address"` //地址
	EthAmount   float64   `json:"eth_amount"`   //以太币余额
	UsdtAmount  float64   `json:"usdt_amount"`  //泰达币余额
	UpdateAt    time.Time `json:"update_at"`    //创建时间
	SearchTime  time.Time `json:"search_time"`  //最后查询余额时间
	Status      int       `json:"status"`       //状态 0：所有，1：已完成，2：执行中，3：结账中
	HandStatus  int       `json:"status"`       //指派状态 0：所有，1：启用，2：停用

}
