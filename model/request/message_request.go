package request

import "time"

type MessageAdd struct {
	CommonMessage
}

type MessageEdit struct {
	ID int `json:"id"`
	CommonMessage
}

type CommonMessage struct {
	Title     string    `json:"title" binding:"required,max=255" example:"这是一个系统公告"`           // 系统公告
	BeginTime time.Time `json:"begin_time" binding:"required" example:"2020-09-15T14:41:46+08:00"` // 开始时间
	EndTime   time.Time `json:"end_time" binding:"required" example:"2020-09-23T14:41:50+08:00"`   // 结束时间
	Users     []int     `json:"users" binding:"required,min=1" example:"1,2,3,4,5,6,7"`            //发给用户id
}
