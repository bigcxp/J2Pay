package response

import (
	"time"
)

type FeePage struct {
	Total       int       `json:"total"`        // 总共多少页
	PerPage     int       `json:"per_page"`     // 当前页码
	CurrentPage int       `json:"current_page"` // 每页显示多少条
	Data        []FeeList `json:"data"`
}

type FeeList struct {
	ID         uint      `json:"id"`          //ID
	RealName   string    `json:"real_name"`   //组织名称
	Amount     float64   `json:"amount"`      //金额
	CreatedAt  time.Time `json:"created_at"`  //创建时间
	FinishTime time.Time `json:"finish_time"` //完成时间
	Status     int       `json:"status"`      //状态 1：执行中 2：已完成
}
