package request

import "time"

type MessageAdd struct {
	CommonMessage
}

type MessageEdit struct {
	Id int `json:"-"`
	CommonMessage
}

type CommonMessage struct {
	Title        string    `binding:"required,max=255" example:"这是一个系统公告"` // 系统公告
	BeginTime    time.Time `binding:"required" example:"2020-09-14"`       // 开始时间
	EndTime      time.Time `binding:"required" example:"2020-09-18"`       // 结束时间
	OrganizeName string    `binding:"required,min=1" example:"张三,李四"`
}
