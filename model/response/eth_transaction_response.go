package response

//eth钱包交易返回对象
type EthTransactionPage struct {
	Total       int              `json:"total"`        // 总共多少页
	PerPage     int              `json:"per_page"`     // 当前页码
	CurrentPage int              `json:"current_page"` // 每页显示多少条
	Data        []EthTransaction `json:"data"`
}

type EthTransaction struct {
	ID             int64   `json:"id"`
	From           string  `json:"from"`            //打币地址
	To             string  `json:"to"`              //充币地址
	Balance        float64 `json:"balance"`         //余额
	ScheduleStatus int     `json:"schedule_status"` //排程状态：1：等待中，:成功,2：失败,3:执行中
	TXID           string  `json:"txid"`            //交易hash
	ChainStatus    int     `json:"chain_status"''"` //链上状态
	CreateTime     int64   `json:"create_time"`     //创建时间戳
}
