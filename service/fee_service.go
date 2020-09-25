package service

import (
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/pkg/casbin"
)

// 手续费订单列表
func FeeList(fromDate string, toDate string, status int, userId int, page, pageSize int) (res response.FeePage, err error) {
	fee := model.Fee{}
	if  status == 0 && fromDate == "" && toDate == "" && userId == 0 {
		res, err = fee.GetAll(page, pageSize)
	} else {
		//将时间进行转换
		res, err = fee.GetAll(page, pageSize, "status = ? or user_id = ? or UNIX_TIMESTAMP(created_at)>=? or  UNIX_TIMESTAMP(created_at) <=?", status, userId, fromDate, toDate)
	}
	return
}

//手续费结账  =>修改状态 计算手续费
func FeeSettle(fee request.FeeEdit)  error{
	defer casbin.ClearEnforcer()
	f := model.Fee{}
	f.ID = uint(fee.Id)
	//逻辑
	return f.Settlement(f.ID)
}

// 创建手续费订单
func FeeAdd(fee request.FeeAdd) error {
	defer casbin.ClearEnforcer()
	f := model.Fee{
		Amount:    fee.Amount,//需要计算手续费
		UserId:    fee.UserId,

	}
	return f.Create()
}
