package service

import (
	"j2pay-server/hcommon"
	"j2pay-server/heth"
	"j2pay-server/model"
	"log"
	"math"
	"strconv"
	"time"
)

type ETHService struct {

}

// 主要逻辑就是先乘，trunc之后再除回去，就达到了保留N位小数的效果
func FormatFloat(num float64, decimal int) string {
	// 默认乘1
	d := float64(1)
	if decimal > 0 {
		// 10的N次方
		d = math.Pow10(decimal)
	}
	// math.trunc作用就是返回浮点数的整数部分
	// 再除回去，小数点后无效的0也就不存在了
	return strconv.FormatFloat(math.Trunc(num*d)/d, 'f', -1, 64)
}


//保存发送交易记录
func (e * ETHService) saveTransaction(ethHelp *heth.ETHHelp,withdrawID int64,transactionType int64,tokenId int64)error{
	_, err := model.SQLCreateTSend(
		&model.TSend{
			RelatedType:  transactionType,
			RelatedID:    withdrawID,
			TokenID:		tokenId,
			TxID:         ethHelp.TxHash,
			FromAddress:  ethHelp.FromAddress,
			ToAddress:    ethHelp.ToAddress,
			BalanceReal:  FormatFloat(ethHelp.SendBalance,  64),
			GasLimit:          ethHelp.SendGasLimit,
			GasPrice:     ethHelp.SendGasPrice,
			Nonce:        ethHelp.Nonce,
			Hex:          ethHelp.RawTxHex,
			CreateTime:time.Now().Unix(),
			HandleStatus: hcommon.SendStatusInit,
			HandleMsg:    "init",
		},
	)
	if err != nil {
		return err
	}
	return nil
}

//交易一致处于padding中，造成交易卡死，重新发起新的交易覆盖之前交易，需要提高gas费用
//param txId  交易发起时的txid
//maxGasJudgment 是否需要判断系统配置最大gas费限制
func (e * ETHService) ERC20PaddingHander(txid string,maxGasJudgment bool)(success bool,err error){
	success=false
	var tSend=model.TSend{}
	tSendMode:= tSend.SQLGetTSendByTXID(txid)
	tokenConfigSql:=model.TAppConfigToken{}
	tokenConfig,err:=  tokenConfigSql.SQLSelectByID(tSendMode.TokenID)
	if err!=nil || tokenConfig.TokenAddress==""{
		log.Print("查询usdt绑定的合约出错，err: [%T] %s", err, err.Error())
		return
	}
	var ethHelp =heth.ETHHelp{}
	ethHelp.ToAddress= tSendMode.ToAddress
	ethHelp.FromAddress = tokenConfig.HotAddress
	ethHelp.ContractAddress=tokenConfig.TokenAddress
	ethHelp.Places= tokenConfig.TokenDecimals
	ethHelp.Nonce=tSendMode.Nonce
	pri,err:= ethHelp.GetPrivateKey(ethHelp.FromAddress)
	if err!=nil {
		return
	}
	ethHelp.PrivateKey=pri
	if err!=nil {
		return
	}
	gasData:= heth.GasData{}
	gasDataSql:= ethHelp.GetGas()
	//交易失败，增加原有gasprice比例20%
	gasData.FastGasPrice=int64(float64( tSendMode.GasPrice)*1.2)
	//本次交易允许的最大值
	var maximumAllowed =gasDataSql.FastGasPrice
	//开关-----判断是否最大
	if maxGasJudgment {
		//如果从eth上获取的gasPrice快速交易参数，大于了系统设置的最大值，就使用系统最大值
		if gasDataSql.FastGasPrice > gasDataSql.MaxGasPrice {
			maximumAllowed = gasDataSql.MaxGasPrice
		}
		//如果gas最大值已经小于当前提交的gas值，那么说明当前交易不合理
		if maximumAllowed <= gasData.FastGasPrice {
		//TODO 取消交易逻辑
		log.Print("当前交易的gas费用不合理，不适合交易")
		return
		}
	}
	//增加20%后还是低于，本次交易允许的最大值
	if gasData.FastGasPrice<maximumAllowed{
		gasData.FastGasPrice=maximumAllowed
	}

	gasData.FastGasLimit=tSendMode.GasLimit
	ethHelp.GasData=&gasData
	if err!=nil {
		return
	}
	chainID, err := ethHelp.GetchainID()
	ethHelp.ChainID=chainID
	ethHelp.SendBalance,err=strconv.ParseFloat(  tSendMode.BalanceReal,64)
	var ethHelpRet,err2=ethHelp.ERC20Transaction(ethHelp)
	if err2!=nil {
		log.Print("调用合约发起交易出错，err: [%T] %s", err2, err2.Error())
		return
	}
	if ethHelpRet.TxHash!=""{
		err=e.saveTransaction(&ethHelpRet,0,hcommon.SendRelationTypeWithdraw,tokenConfig.ID);
		if err!=nil {
			log.Print("保存交易信息出现异常！，err: [%T] %s", err, err.Error())
		}
		success=true
		return
	}
	return
}


//发起代币（erc20）交易
//param toAddress 发送目标地址
//param 发送数量
//maxGasJudgment 是否需要判断系统配置最大gas费限制
func (e * ETHService) ERC20Transaction(toAddress string,quantitySent float64,maxGasJudgment bool)(success bool,err error){
	success=false
	tokenConfigSql:=model.TAppConfigToken{}
	tokenConfig,err:=  tokenConfigSql.SQLSelectBySymbol("usdt")
	if err!=nil{
		log.Printf("查询usdt绑定的合约出错，err: [%T] %s", err, err.Error())
		return
	}
	var ethHelp =heth.ETHHelp{}
	ethHelp.ToAddress= toAddress
	ethHelp.FromAddress = tokenConfig.HotAddress
	ethHelp.ContractAddress=tokenConfig.TokenAddress
	ethHelp.Places= tokenConfig.TokenDecimals
	nonce, err :=	ethHelp.GetNONCE(ethHelp.FromAddress)
	if err!=nil {
		return
	}
	println("nonce====",nonce)
	ethHelp.Nonce=nonce
	pri,err:= ethHelp.GetPrivateKey(ethHelp.FromAddress)
	if err!=nil {
		return
	}
	ethHelp.PrivateKey=pri
	if err!=nil {
		return
	}
	gasData:=  ethHelp.GetGas()
	ethHelp.GasData=gasData
	if err!=nil {
		return
	}
	//开关------   判断是否最大
	if maxGasJudgment {
		//如果从eth上获取的gasPrice快速交易参数，大于了系统设置的最大值，就使用系统最大值
		if gasData.FastGasPrice > gasData.MaxGasPrice {
			ethHelp.GasData.FastGasPrice = gasData.MaxGasPrice
		}
	}
	chainID, err := ethHelp.GetchainID()
	ethHelp.ChainID=chainID
	ethHelp.SendBalance=quantitySent
	var ethHelpRet,err2=ethHelp.ERC20Transaction(ethHelp)
	if err2!=nil {
		log.Print("调用合约发起交易出错，err: [%T] %s", err2, err2.Error())
		return
	}
	if ethHelpRet.TxHash!=""{
		err=e.saveTransaction(&ethHelpRet,0,hcommon.SendRelationTypeWithdraw,tokenConfig.ID);
		if err!=nil {
			log.Print("保存交易信息出现异常！，err: [%T] %s", err, err.Error())
		}
		success=true
		return
	}
	return
}