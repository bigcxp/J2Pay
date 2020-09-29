package service

import (
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/myerr"
	"j2pay-server/pkg/casbin"
	"j2pay-server/pkg/util"
)

// 订单列表
func ReturnList(fromDate string, toDate string, status int, orderCode string, page, pageSize int) (res response.ReturnPage, err error) {
	returnMoney := model.Return{}
	if  status == 0 && fromDate == "" && toDate == "" && orderCode == "" {
		res, err = returnMoney.GetAll(page, pageSize)
	} else {
		//将时间进行转换
		res, err = returnMoney.GetAll(page, pageSize, "status = ? or order_code like ?  or UNIX_TIMESTAMP(created_at)>=? or  UNIX_TIMESTAMP(created_at) <=?", status, orderCode, fromDate, toDate)
	}
	return
}

// 订单详情
func ReturnDetail(id uint) (res response.ReturnList, err error) {
	returns := model.Return{}
	returns.ID = id
	res, err = returns.GetDetail()
	if err != nil {
		return
	}
	return res, err
}

// 新增订单
func ReturnAdd(returns request.ReturnAdd) error {
	defer casbin.ClearEnforcer()
	r := model.Return{
		SystemCode: util.RandString(20),
		OrderCode:  returns.OrderCode,
		Amount:     returns.Amount,
		UserId:     returns.UserId,

	}
	//逻辑处理 待完善 ==》随机分配地址 实际收款明细表code
	//判断订单编号是否重复
	if hasCode := model.GetReturnByWhere("order_code = ?", returns.OrderCode); hasCode.ID > 0 {
		return myerr.NewDbValidateError("商户订单编号重复！")
	}
	return r.Create()
}
