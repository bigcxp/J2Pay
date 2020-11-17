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
			Gas:          ethHelp.SendGasLimit,
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

//发起代币（erc20）交易
//@param toAddress 发送目标地址
//@param 发送数量
func (e * ETHService) ERC20Transaction(toAddress string,quantitySent float64)(success bool,err error){
	success=false
	tokenConfigSql:=model.TAppConfigToken{}
	tokenConfig,err:=  tokenConfigSql.SQLSelectBySymbol("usdt")
	if err!=nil{
		log.Panicf("查询usdt绑定的合约出错，err: [%T] %s", err, err.Error())
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
	chainID, err := ethHelp.GetchainID()
	ethHelp.ChainID=chainID
	ethHelp.SendBalance=quantitySent
	var ethHelpRet,err2=ethHelp.ERC20Transaction(ethHelp)
	if err2!=nil {
		log.Panicf("调用合约发起交易出错，err: [%T] %s", err2, err2.Error())
		return
	}
	if ethHelpRet.TxHash!=""{
		err=e.saveTransaction(&ethHelpRet,0,hcommon.SendRelationTypeWithdraw,tokenConfig.ID);
		if err!=nil {
			log.Panicf("保存交易信息出现异常！，err: [%T] %s", err, err.Error())
		}
		success=true
		return
	}
	return
}