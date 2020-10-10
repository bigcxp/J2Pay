package service

import (
	"j2pay-server/model"
	"j2pay-server/model/response"
)

//获取钱包列表
func AddressList(address string, status, handStatus, userId, page, pageSize int) (res response.AddressPage, err error) {
	addresses := model.Address{}
	if status == 0 && address == "" && handStatus == 0 && userId == 0 {
		res, err = addresses.GetAllAddress(page, pageSize)
	} else {
		//根据条件查询
		res, err = addresses.GetAllAddress(page, pageSize, "status = ? or handStatus = ? or user_id = ? or address like ?", status, handStatus, userId, address)
	}
	return
}

//按照一定数量和是否启用钱包新建钱包

