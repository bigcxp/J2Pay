package response

import "time"

//热钱包交易记录返回对象
type HotTransactionPage struct {
	Total       int              `json:"total"`        // 总共多少页
	PerPage     int              `json:"per_page"`     // 当前页码
	CurrentPage int              `json:"current_page"` // 每页显示多少条
	Data        []HotTransaction `json:"data"`
}

type HotTransaction struct {
	ID             int64     `json:"id"`
	SystemCode     string    `json:"system_code"`     //系统编号
	From           string    `json:"from"`            //打币地址
	To             string    `json:"to"`              //充币地址
	Type           int       `json:"type"`            //类型:1:代发,2:排程结账,3:手动结账
	GasFee         float64   `json:"gas_fee"`         //gas费
	Balance        float64   `json:"balance"`         //余额
	ScheduleStatus int       `json:"schedule_status"` //排程状态：1：等待中，:成功,2：失败,3:执行中
	TXID           string    `json:"txid"`            //交易hash
	ChainStatus    int       `json:"chain_status"''"` //链上状态
	CreateTime     time.Time `json:"create_time"`     //创建时间戳
}
