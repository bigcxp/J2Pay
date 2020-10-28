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
func OrderList(fromDate string, toDate string, status int, chargeAddress string, txid string, orderCode string, userId int,page, pageSize int) (res response.OrderPage, err error) {
	order := model.Order{}
	if userId == 0 {
		if chargeAddress == "" && status == 0 && fromDate == "" && toDate == "" && orderCode == "" && txid == "" {
			res, err = order.GetAllMerchantOrder(page, pageSize)
		} else {
			//将时间进行转换
			res, err = order.GetAllMerchantOrder(page, pageSize, "status = ? or charge_address like ? or order_code like ?  or txid like ? or UNIX_TIMESTAMP(created_at)>=? or  UNIX_TIMESTAMP(created_at) <=?", status, chargeAddress, orderCode, txid, fromDate, toDate)
		}
	}else {
		if chargeAddress == "" && status == 0 && fromDate == "" && toDate == "" && orderCode == "" && txid == "" {
			res, err = order.GetAllMerchantOrder(page, pageSize,"user_id = ?",userId)
		} else {
			//将时间进行转换
			res, err = order.GetAllMerchantOrder(page, pageSize, "user_id = ? and status = ? or charge_address like ? or order_code like ?  or txid like ? or UNIX_TIMESTAMP(created_at)>=? or  UNIX_TIMESTAMP(created_at) <=?",userId, status, chargeAddress, orderCode, txid, fromDate, toDate)
		}
	}

	return
}

// 订单详情
func OrderDetail(id uint) (res response.RealOrderList, err error) {
	order := model.Order{}
	order.ID = id
	res, err = order.GetDetail()
	if err != nil {
		return
	}
	return res, err
}

// 新增订单
func OrderAdd(order request.OrderAdd) error {
	defer casbin.ClearEnforcer()
	o := model.Order{
		OrderCode: order.OrderCode,
		IdCode:    util.RandString(20),
		Amount:    order.Amount,
		UserId:    order.UserId,
		Remark:    order.Remark,
	}
	//逻辑处理 待完善 ==》随机分配地址 实际收款明细表code
	//判断订单编号是否重复
	if hasCode := model.GetOrderByWhere("order_code = ?", order.OrderCode); hasCode.ID > 0 {
		return myerr.NewDbValidateError("商户订单编号重复！")
	}
	return o.Create()
}

//修改订单
func OrderEdit(order request.OrderEdit) error {
	defer casbin.ClearEnforcer()
	o := model.Order{}
	o.ID = uint(order.ID)
	return o.UpdateOrder(order)
}
