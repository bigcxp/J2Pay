package service

import (
	"j2pay-server/model"
)

// 提领订单列表
func PickList(fromDate string,toDate string,status int,name string, page, pageSize int) (res model.PickUpPage, err error) {
	pick := model.Pick{}
	if name == "" && status ==0 && fromDate =="" &&toDate ==""{
		res, err = pick.GetAll(page, pageSize)
	} else {
		//将时间进行转换
		res, err = pick.GetAll(page, pageSize, "status = ? or real_name like ? and UNIX_TIMESTAMP(created_at)>=? and  UNIX_TIMESTAMP(created_at) <=?", status, name,fromDate,toDate)
	}
	return
}


// 提领订单详情
func PickDetail(id uint) (res model.Pick, err error) {
	pick := model.Pick{}
	pick.ID = id
	res,err := pick.getDetail()
	if err != nil {
		return
	}
	return
}


