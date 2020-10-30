package service

import (
	"github.com/gin-gonic/gin"
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/myerr"
	"j2pay-server/pkg/casbin"
	"j2pay-server/pkg/util"
	"time"
)

// 订单列表
func OrderList(fromDate string, toDate string, status int, chargeAddress string, txid string, orderCode string, userId int, page, pageSize int) (res response.OrderPage, err error) {
	order := model.Order{}
	if userId == 0 {
		if chargeAddress == "" && status == 0 && fromDate == "" && toDate == "" && orderCode == "" && txid == "" {
			res, err = order.GetAllMerchantOrder(page, pageSize)
		} else {
			//将时间进行转换
			res, err = order.GetAllMerchantOrder(page, pageSize, "status = ? or charge_address like ? or order_code like ?  or txid like ? or UNIX_TIMESTAMP(created_at)>=? or  UNIX_TIMESTAMP(created_at) <=?", status, chargeAddress, orderCode, txid, fromDate, toDate)
		}
	} else {
		if chargeAddress == "" && status == 0 && fromDate == "" && toDate == "" && orderCode == "" && txid == "" {
			res, err = order.GetAllMerchantOrder(page, pageSize, "user_id = ?", userId)
		} else {
			//将时间进行转换
			res, err = order.GetAllMerchantOrder(page, pageSize, "user_id = ? and status = ? or charge_address like ? or order_code like ?  or txid like ? or UNIX_TIMESTAMP(created_at)>=? or  UNIX_TIMESTAMP(created_at) <=?", userId, status, chargeAddress, orderCode, txid, fromDate, toDate)
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
func OrderAdd(order request.OrderAdd) (error) {
	defer casbin.ClearEnforcer()
	//如果是RMB或TWD 换算成USDT
	var amount float64
	switch order.Currency {
	case "RMB":
		detail, err := TypeDetail(order.Currency)
		if err != nil {
			return err
		}
		amount = detail.OriginalRate * order.Amount
	case "TWB":
		detail, err := TypeDetail(order.Currency)
		if err != nil {
			return err
		}
		amount = detail.OriginalRate * order.Amount

	default:
		amount = order.Amount
	}
	var c *gin.Context

	//判断订单编号是否重复
	if hasCode := model.GetOrderByWhere("order_code = ?", order.OrderCode); hasCode.ID > 0 {
		return myerr.NewDbValidateError("商户订单编号重复！")
	}
	//如果用户id为0 商户端操作 获取当前登录用户
	if order.UserId == 0 {
		user, hasUser := c.Get("user")
		if !hasUser {
			return myerr.NewNormalValidateError("用户未登录")
		}
		userInfo := user.(*util.Claims)
		//检测当前用户的收钱地址是否满足 随机查询一个状态为已完成的地址
		address, err := model.GetAddress(userInfo.Id)
		if err != nil {
			return myerr.NewNormalValidateError("用户收款地址不足")
		}
		o := model.Order{
			OrderCode:  order.OrderCode,
			IdCode:     util.RandString(20),
			Amount:     amount,
			Address:    address.UserAddress,
			UserId:     userInfo.Id,
			CreateTime: order.Uts,
			ExprireTime:order.Uts+3600,
			Remark:     order.Remark,
		}
		return o.Create()
		//管理员添加订单
	} else {
		//检测当前用户的收钱地址是否满足 随机查询一个状态为已完成的地址
		address, err := model.GetAddress(order.UserId)
		if err != nil {
			return myerr.NewNormalValidateError("用户收款地址不足")
		}
		o := model.Order{
			OrderCode:  order.OrderCode,
			IdCode:     util.RandString(20),
			Amount:     amount,
			Address:    address.UserAddress,
			UserId:     order.UserId,
			CreateTime: order.Uts,
			Remark:     order.Remark,
		}
	}

	o := model.Order{
		OrderCode: order.OrderCode,
		IdCode:    util.RandString(20),
		Amount:    order.Amount,
		UserId:    order.UserId,
		Remark:    order.Remark,
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
