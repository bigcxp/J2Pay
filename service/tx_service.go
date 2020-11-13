package service

import (
	"j2pay-server/heth"
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/pkg/casbin"
	"j2pay-server/pkg/util"
	"time"
)

//获取所有的eth交易列表
func EthTxList(fromTime, toTime, page, pageSize int) (res response.EthTransactionPage, err error) {
	Ethtx := model.EthTransaction{}
	if fromTime == 0 && toTime == 0 {
		res, err = Ethtx.GetAll(page, pageSize)
	} else {
		res, err = Ethtx.GetAll(page, pageSize, "create_time >= ? and create_time <= ?", fromTime, toTime)
	}
	return
}

//获取所有的hot交易列表
func HotTxList(fromAddr, toAddr string, scheduleStatus, chainStatus, fromTime, toTime, page, pageSize int) (res response.EthTransactionPage, err error) {
	Ethtx := model.EthTransaction{}
	if fromTime == 0 && toTime == 0 && scheduleStatus == 0 && chainStatus == 0 && fromAddr == "" && toAddr == "" {
		res, err = Ethtx.GetAll(page, pageSize)
	} else {
		res, err = Ethtx.GetAll(page, pageSize, "create_time >= ? and create_time <= ? or from = ? or to = ? or schedule_status = ? or chain_status = ?", fromTime, toTime, fromAddr, toAddr, scheduleStatus, chainStatus)
	}
	return
}

//创建eth钱包交易 储存
func CreateEthTx(eth request.EthTxAdd) error {
	defer casbin.ClearEnforcer()
	nowTime := time.Now().Unix()
	//获取eth钱包地址
	ethAddress := model.GetAddressByWhere("handle_status = ? and use_tag = ?", 1, -2)
	e := model.EthTransaction{
		From:           ethAddress.UserAddress,
		To:             eth.To,
		Balance:        eth.Balance,
		ScheduleStatus: 1,
		TXID:           "",
		ChainStatus:    1,
		CreateTime:     nowTime,
	}
	return e.AddEthTx()

}

//创建hot钱包交易 排程结账 手动结账 代发
func CreateHotTx(hot request.HotTxAdd) error {
	defer casbin.ClearEnforcer()
	nowTime := time.Now().Unix()
	//获取eth钱包地址
	hotAddress := model.GetAddressByWhere("handle_status = ? and use_tag = ?", 1, -1)
	//获取商户钱包地址 商户充币地址归集到主钱包
	adminUser,_ := model.GetUserByWhere("id = ?", hotAddress.UseTag)
	//获取手续费 先检测gaslimit
	heth.CheckGasPrice()
	//获取gasPrice
	gasPriceValue:= model.SQLGetTAppStatusIntValueByK("to_user_gas_price")
	gasLimitValue := model.SQLGetTAppStatusIntValueByK( "gas_limit")
	feeValue := *gasPriceValue * *gasLimitValue
	//判断是代发还是提领还是结账 结账=》排程结账 手动结账
	//1:代发,2:排程结账,3:手动结账
	if hot.Type == 1 {
		h := model.HotTransaction{
			SystemCode:     util.RandString(20),
			From:           hotAddress.UserAddress,
			To:             hot.To,
			Balance:        hot.Balance,
			Type:           1,
			GasFee:         feeValue,
			ScheduleStatus: 1,
			TXID:           "",
			ChainStatus:    1,
			CreateTime:     nowTime,
		}
		return h.AddHotTx()
		//排程结账
	} else if hot.Type == 2{
		h := model.HotTransaction{
			SystemCode:     util.RandString(20),
			From:           hot.From,
			To:             adminUser.Address,
			Balance:        hot.Balance,
			Type:           2,
			GasFee:         feeValue,
			ScheduleStatus: 1,
			TXID:           "",
			ChainStatus:    1,
			CreateTime:     nowTime,
		}
		return h.AddHotTx()
	}else {
		h := model.HotTransaction{
			SystemCode:     util.RandString(20),
			From:           hot.From,
			To:             adminUser.Address,
			Balance:        hot.Balance,
			Type:           3,
			GasFee:         feeValue,
			ScheduleStatus: 1,
			TXID:           "",
			ChainStatus:    1,
			CreateTime:     nowTime,
		}
		return h.AddHotTx()
	}

}
