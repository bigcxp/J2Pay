package response

import "time"

type SystemMessagePage struct {
	Total       int                 `json:"total"`        // 总共多少页
	PerPage     int                 `json:"per_page"`     // 当前页码
	CurrentPage int                 `json:"current_page"` // 每页显示多少条
	Data        []SystemMessageList `json:"data"`
}

type AdminName struct {
	Name string `json:"name"` // 组织名
}

type SystemMessageList struct {
	Id           int    `json:"id"`            //id
	Title        string `json:"title"`         // 标题
	BeginTime    time.Time `json:"begin_time"`    //开始时间
	EndTime      time.Time `json:"end_time"`      //结束时间
	Status       int8   `json:"status"`        //是否作废 0：否，1：是
	OrganizeName string `json:"organize_name"` //组织名称
}
