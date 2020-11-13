package service

import (
	"j2pay-server/heth"
	"j2pay-server/model"
	"log"
)

type ETHService struct {

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
	gas,limit:=  ethHelp.GetGas()
	ethHelp.GasPrice=gas
	ethHelp.GasLimit=limit
	if err!=nil {
		return
	}
	chainID, err := ethHelp.GetchainID()
	ethHelp.ChainID=chainID
	ethHelp.SendBalance=quantitySent
	success,err=ethHelp.ETHTransaction(ethHelp)
	if err!=nil {
		log.Panicf("调用合约发起交易出错，err: [%T] %s", err, err.Error())
		return
	}
	return
}