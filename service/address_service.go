package service

import (
	"j2pay-server/heth"
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/myerr"
)

//获取钱包列表 商户钱包列表 热钱包列表
func AddressList(address string, status, handStatus, userId, useTag, page, pageSize int) (res response.AddressPage, err error) {
	addresses := model.Address{}
	//根据usetag来判断需要查询的钱包类型
	switch useTag {
	//eth地址列表
	case -2:
		res, err = addresses.GetAllAddress(page, pageSize, "use_tag = ? ", useTag)
	//hot地址列表
	case -1:
		res, err = addresses.GetAllAddress(page, pageSize, "use_tag = ? ", useTag)
	//商户充币地址列表
	default:
		if status == 0 && address == "" && handStatus == 0 && userId == 0 {
			res, err = addresses.GetAllAddress(page, pageSize, "use_tag > 0")
		} else {
			//根据条件查询
			res, err = addresses.GetAllAddress(page, pageSize, "use_tag >0 or status = ? "+
				"or handle_status = ? or user_id = ? or address = ? ", status, handStatus, userId, address)
		}
	}
	return
}

//按照一定数量和是否启用钱包新建热钱包，eth钱包 ，检测是否有空闲钱包 为商户分配钱包
func AddAddress(addr request.AddressAdd) (err error) {
	//检测空余地址是否足够 不够则生成新的地址 数据库中设置好需要生成的地址数量
	_, err = heth.CheckAddressFree()
	if err != nil {
		return
	}
	//根据需要生成的地址执行对应方法
	switch addr.UseTag {
	//生成热钱包地址
	case -1:
		if addr.Num > 0 {
			_, err := heth.CreateHotAddress(addr)
			return err

		}
	//生成eth钱包地址
	case -2:
		if addr.Num > 0 {
			_, err := heth.CreateHotAddress(addr)
			return err
		}
	//生成用户钱包地址 将已经生成好的空闲钱包地址
	default:
		if addr.Num > 0 {
			err := heth.ToMerchantAddress(addr)
			return err
		}
	}
	return
}

//启用 停用 地址
func RestartAddr(address request.OpenOrStopAddress) (err error) {
	//根据ids查询出地址
	allAddress := model.GetAllAddress("id in (?)", address.Id)
	//遍历地址
	for _, v := range allAddress {
		//判断是否满足状态
		if v.HandleStatus == address.HandleStatus {
			return myerr.NewDbValidateError("失败：地址：" + v.UserAddress + ", assign_status_not_allow")
		} else {
			err = model.OpenOrStopAddress(address.HandleStatus, allAddress)
		}
	}
	return err

}

//更新钱包余额
func UpdateBalance(ids request.UpdateAmount)(err error)  {
	err = model.UpdateBalance(ids)
	if err != nil{
		return err
	}
	return
}