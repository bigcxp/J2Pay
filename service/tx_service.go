package service

import (
	"j2pay-server/model"
	"j2pay-server/model/response"
)

//获取所有的eth交易列表
func EthTxList(fromTime,toTime, page, pageSize int) (res response.EthTransactionPage, err error) {
	Ethtx := model.EthTransaction{}
	if fromTime == 0 && toTime == 0 {
		res, err = Ethtx.GetAll(page, pageSize)
	} else {
		res, err = Ethtx.GetAll(page, pageSize,"create_time >= ? and create_time <= ?",fromTime,toTime)
	}
	return
}

//获取所有的hot交易列表
func HotTxList(fromAddr,toAddr string,scheduleStatus,chainStatus ,fromTime,toTime, page, pageSize int) (res response.EthTransactionPage, err error) {
	Ethtx := model.EthTransaction{}
	if fromTime == 0 && toTime == 0 && scheduleStatus == 0 && chainStatus == 0 && fromAddr == "" && toAddr =="" {
		res, err = Ethtx.GetAll(page, pageSize)
	} else {
		res, err = Ethtx.GetAll(page, pageSize,"create_time >= ? and create_time <= ? or from = ? or to = ? or schedule_status = ? or chain_status = ?",fromTime,toTime,fromAddr,toAddr,scheduleStatus,chainStatus)
	}
	return
}

//创建eth