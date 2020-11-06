package service

import (
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/myerr"
	"j2pay-server/pkg/casbin"
	"j2pay-server/pkg/logger"
	_ "j2pay-server/pkg/util"
	"time"
)

// 系统所有公告列表
func MessageList(title string, page, pageSize int) (res response.SystemMessagePage, err error) {
	systemMessage := model.SystemMessage{}
	if title == "" {
		res, err = systemMessage.GetAll(page, pageSize)
	} else {
		res, err = systemMessage.GetAll(page, pageSize, "title like ?", "%"+title+"%")
	}
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
				logger.Logger.Error("商户获取错误: message_id = ", v.Id)
				continue
			}
			res.Data[i].Users = append(res.Data[i].Users, users[user])
		}
	}
	return
}

//根据用户名获取公告列表
func MessageListByUser(uid, page, pageSize int) (res response.AdminUserMessagePage, err error) {
	user := model.AdminUser{}

	if uid != 0 {
		res, err = user.GetAllMessage(page, pageSize, "uid = ?", uid)

	} else {
		return
	}

	if err != nil {
		return
	}
	messages := model.GetAllMessage()
	mappings := model.GetUserMessageMapping()
	for i, v := range res.Data {
		_, ok := mappings[v.Id]
		if !ok {
			continue
		}
		res.Data[i].SystemMessages = []response.AdminSystemMessage{}
		for _, message := range mappings[v.Id] {
			if _, ok := messages[message]; !ok {
				logger.Logger.Error("系统公告获取错误: user_id = ", v.Id)
				continue
			}
			res.Data[i].SystemMessages = append(res.Data[i].SystemMessages, messages[message])
		}
	}
	return
}

// 创建公告
func MessageAdd(message request.MessageAdd) error {
	defer casbin.ClearEnforcer()
	m := model.SystemMessage{
		Title:     message.Title,
		BeginTime: time.Unix(message.BeginTime, 0),
		EndTime:   time.Unix(message.EndTime, 0),
	}
	// 判断用户是否存在
	hasUsers, err := model.GetUsersByWhere("id in (?)", message.Users)
	if err != nil {
		return err
	}
	if len(hasUsers) != len(message.Users) {
		return myerr.NewDbValidateError("选择的商户不存在")
	}
	return m.Create(message.Users)
}

// 编辑公告
func MessageEdit(message request.MessageEdit) error {
	defer casbin.ClearEnforcer()
	m := model.SystemMessage{
		ID:        message.ID,
		Title:     message.Title,
		BeginTime: time.Unix(message.BeginTime, 0),
		EndTime:   time.Unix(message.EndTime, 0),
	}
	// 判断用户是否存在
	hasUsers, err := model.GetUsersByWhere("id in (?)", message.Users)
	if err != nil {
		return err
	}
	if len(hasUsers) != len(message.Users) {
		return myerr.NewDbValidateError("选择的商户不存在")
	}
	return m.Edit(message.Users)
}

// 删除公告
func MessageDel(id int) error {
	defer casbin.ClearEnforcer()
	m := model.SystemMessage{
		ID: id,
	}
	return m.Del()
}
