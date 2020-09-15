package service

import (
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/pkg/casbin"
	_ "j2pay-server/pkg/util"
)

// 系统公告列表
func MessageList(page, pageSize int) (res response.SystemMessagePage, err error) {
	systemMessage := model.SystemMessage{}
	return systemMessage.GetAll(page, pageSize)
}

// 创建公告
func MessageAdd(message request.MessageAdd) error {
	defer casbin.ClearEnforcer()
	m := model.SystemMessage{
		Title: message.Title,
		OrganizeName: message.OrganizeName,
		BeginTime: message.BeginTime,
		EndTime: message.EndTime,
	}
	return m.Create()
}

// 编辑公告
func MessageEdit(message request.MessageEdit) error {
	defer casbin.ClearEnforcer()
	m := model.SystemMessage{
		Id:       message.Id,
		Title: message.Title,
		OrganizeName: message.OrganizeName,
		BeginTime: message.BeginTime,
		EndTime: message.EndTime,
	}
	return m.Edit()
}

// 删除用户
func MessageDel(id int) error {
	defer casbin.ClearEnforcer()
	m := model.SystemMessage{
		Id: id,
	}
	return m.Del()
}

