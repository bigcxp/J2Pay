package service

import (
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/myerr"
	"j2pay-server/pkg/casbin"
	"j2pay-server/pkg/logger"
	_ "j2pay-server/pkg/util"
)

// 系统公告列表
func MessageList(page, pageSize int) (res response.SystemMessagePage, err error) {
	systemMessage := model.SystemMessage{}
	res, err = systemMessage.GetAll(page, pageSize)
	if err != nil {
		return
	}
	users := model.GetAllUser()
	mappings := model.GetMessageUserMapping()
	for i, v := range res.Data {
		_, ok := mappings[v.Id]
		if !ok {
			continue
		}
		res.Data[i].Users = []response.UserNames{}
		for _, user := range mappings[v.Id] {
			if _, ok := users[user]; !ok {
				logger.Logger.Error("用户获取错误: message_id = ", v.Id)
				continue
			}
			res.Data[i].Users = append(res.Data[i].Users, users[user])
		}
	}
	return
}

// 创建公告
func MessageAdd(message request.MessageAdd) error {
	defer casbin.ClearEnforcer()
	m := model.SystemMessage{
		Title: message.Title,
		BeginTime: message.BeginTime,
		EndTime: message.EndTime,
	}
	// 判断用户是否存在
	hasUsers, err := model.GetUsersByWhere("id in (?)", message.Users)
	if err != nil {
		return err
	}
	if len(hasUsers) != len(message.Users) {
		return myerr.NewDbValidateError("选择的用户不存在")
	}
	return m.Create(message.Users)
}

// 编辑公告
func MessageEdit(message request.MessageEdit) error {
	defer casbin.ClearEnforcer()
	m := model.SystemMessage{
		Id:       message.Id,
		Title: message.Title,
		BeginTime: message.BeginTime,
		EndTime: message.EndTime,
	}
	// 判断用户是否存在
	hasUsers, err := model.GetUsersByWhere("id in (?)", message.Users)
	if err != nil {
		return err
	}
	if len(hasUsers) != len(message.Users) {
		return myerr.NewDbValidateError("选择的用户不存在")
	}
	return m.Edit(message.Users)
}

// 删除公告
func MessageDel(id int) error {
	defer casbin.ClearEnforcer()
	m := model.SystemMessage{
		Id: id,
	}
	return m.Del()
}

