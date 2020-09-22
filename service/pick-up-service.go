package service

import (
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/myerr"
	"j2pay-server/pkg/casbin"
	"j2pay-server/pkg/util"
)

// 提领订单列表
func PickList(fromDate string, toDate string, status int, name string, code string,types int, page, pageSize int) (res model.PickUpPage, err error) {
	pick := model.Pick{}
	if name == "" && status == 0 && fromDate == "" && toDate == "" && code == "" && types == 0 {
		res, err = pick.GetAll(page, pageSize)
	} else {
		//将时间进行转换
		res, err = pick.GetAll(page, pageSize, "status = ? or type = ? or real_name like ? or is_code like ? or UNIX_TIMESTAMP(created_at)>=? or  UNIX_TIMESTAMP(created_at) <=?", status,types, name, code, fromDate, toDate)
	}
	return
}

// 提领订单详情
func PickDetail(id uint) (res model.Pick, err error) {
	pick := model.Pick{}
	pick.ID = id
	res, err = pick.GetDetail()
	if err != nil {
		return
	}
	return res, err
}

// 提领,代发
func PickAdd(pick request.PickAdd) error {
	defer casbin.ClearEnforcer()
	p := model.Pick{
		OrderCode:   pick.OrderCode,
		Amount:      pick.Amount,
		Type:        pick.Type,
		UserId:      pick.UserId,
		Remark:      pick.Remark,
		PickAddress: pick.PickAddress,
		IdCode:      util.RandString(20),
	}
	//逻辑处理 待完善 ==》判断金额 如果足够 进行交易，返回交易结果 账户减余额
	// 1.判断转账金额是否足够
	if user := model.GetUserByWhere("id = ?", pick.UserId); user.Balance < pick.Amount {
		return myerr.NewDbValidateError("余额不足")
	}
	return p.Create()
}
