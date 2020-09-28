package service

import (
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/myerr"
	"j2pay-server/pkg/casbin"
	"j2pay-server/pkg/util"
)

// 商户端  提领订单列表 + 代发订单列表
func MerchantPickList(fromDate, toDate, code string, types, userId, status, page, pageSize int) (res response.MerchantPickSendPage, err error) {
	pick := model.Pick{}
	if userId == 0 {
		return res, myerr.NewDbValidateError("没有此用户信息")
	}
	if status == 0 && fromDate == "" && toDate == "" && code == "" && types == 0 {
		res, err = pick.GetAll(page, pageSize, "user_id = ?", userId)
	} else {
		//将时间进行转换
		res, err = pick.GetAll(page, pageSize, "user_id = ? and status = ? or type = ? or real_name like ? or is_code like ? or UNIX_TIMESTAMP(created_at)>=? or  UNIX_TIMESTAMP(created_at) <=?", userId, status, types, code, fromDate, toDate)
	}

	return
}

// 管理端 提领订单列表
func PickUpList(fromDate, toDate string, status, types, userId, page, pageSize int) (res response.PickUpPage, err error) {
	pick := model.Pick{}

	if status == 0 && fromDate == "" && toDate == "" && userId == 0 {
		res, err = pick.GetAllPick(page, pageSize, "type = ?", types)
	} else {
		//将时间进行转换
		res, err = pick.GetAllPick(page, pageSize, "type = ? and status = ?  or user_id = ?  or UNIX_TIMESTAMP(created_at)>=? or  UNIX_TIMESTAMP(created_at) <=?", types, status, userId, fromDate, toDate)
	}
	return
}

// 管理端  代发订单列表
func SendList(fromDate, toDate, orderCode string, status, types, userId, page, pageSize int) (res response.SendPage, err error) {
	pick := model.Pick{}
	if orderCode == "" && status == 0 && fromDate == "" && toDate == "" && userId == 0 {
		res, err = pick.GetAllSend(page, pageSize, "type = ?", types)
	} else {
		//将时间进行转换
		res, err = pick.GetAllSend(page, pageSize, "type = ? and status = ? or orderCode like ? or user_id = ? or UNIX_TIMESTAMP(created_at)>=? or  UNIX_TIMESTAMP(created_at) <=?", types, status, orderCode, userId, fromDate, toDate)
	}
	return
}

// 商户端提领订单或者代发订单详情
func MerchantPickDetail(id uint) (res response.MerchantPickList, err error) {
	pick := model.Pick{}
	pick.ID = id
	res, err = pick.GetPickSendDetail()
	if err != nil {
		return
	}
	return res, err
}

// 管理端提领订单详情
func PickDetail(id uint) (res response.PickList, err error) {
	pick := model.Pick{}
	pick.ID = id
	res, err = pick.GetPickDetail()
	if err != nil {
		return
	}
	return res, err
}

// 管理端代发订单详情
func SendDetail(id uint) (res response.SendList, err error) {
	pick := model.Pick{}
	pick.ID = id
	res, err = pick.GetSendDetail()
	if err != nil {
		return
	}
	return res, err
}

// 提领
func PickAdd(pick request.PickAdd) error {
	defer casbin.ClearEnforcer()
	p := model.Pick{
		IdCode:      util.RandString(20),
		Amount:      pick.Amount,
		Type:        pick.Type,
		UserId:      pick.UserId,
		Remark:      pick.Remark,
		PickAddress: pick.PickAddress,
	}
	//逻辑处理 待完善 ==》判断金额 如果足够 进行交易，返回交易结果 账户减余额
	// 1.判断转账金额是否足够
	if user := model.GetUserByWhere("id = ?", pick.UserId); user.Balance < pick.Amount {
		return myerr.NewDbValidateError("余额不足")
	}
	return p.Create()
}

// 代发
func SendAdd(send request.SendAdd) error {
	defer casbin.ClearEnforcer()
	p := model.Pick{
		OrderCode:   send.OrderCode,
		IdCode:      util.RandString(20),
		Amount:      send.Amount,
		Type:        send.Type,
		UserId:      send.UserId,
		Remark:      send.Remark,
		PickAddress: send.PickAddress,
	}
	//逻辑处理 待完善 ==》判断金额 如果足够 进行交易，返回交易结果  修改订单状态 账户减余额
	// 1.判断转账金额是否足够
	if user := model.GetUserByWhere("id = ?", send.UserId); user.Balance < send.Amount {
		return myerr.NewDbValidateError("余额不足")
	}
	return p.Create()
}

//取消提领 代发
func CancelPick(send request.SendEdit)  error{
	defer casbin.ClearEnforcer()
	p := model.Pick{}
	p.ID = uint(send.Id)
	//逻辑
	return p.CancelPick(send.Id,send.Status)
}
