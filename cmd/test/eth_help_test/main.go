	package main

	import (
		"flag"
		"fmt"
		"j2pay-server/ethclient"
		"j2pay-server/heth"
		"j2pay-server/model"
		"j2pay-server/pkg/logger"
		"j2pay-server/pkg/setting"
		"j2pay-server/service"
	)

	func init()  {
		//把用户传递的命令行参数解析为对应变量的值
		flag.Parse()
		// 初始化操作 (因为 init 方法无法保证我们想要的顺序)
		setting.Setup()
		//日志
		logger.Setup()
		//初始化数据库
		model.Setup()
		//初始化以太坊节点
		ethclient.InitClient(fmt.Sprintf("%s", setting.EthConf.Url))
	}
func main (){


}

	func check1(){
		var ethHelp =heth.ETHHelp{}
		ethHelp.CheckTransactionByTxHash("0x1968a5ec2226ccc85bc5644b59b3ea1de300342839f1491d838f8afd278f5275")
	}
	func tran1(){
		// 0x64fe95228af4945d81087f81e9382e3623723f21
		//0x99E46a01909078BC8D0A53CA8EFa2a0B7Ec1497c
		var  ToAddress= "0x99E46a01909078BC8D0A53CA8EFa2a0B7Ec1497c"

		var ethService= service.ETHService{}
		var ethxx=heth.ETHHelp{}
		ethxx.GetETHData()
		//覆盖之前交易
		//	ethService.ERC20PaddingHander("0xb27f4b6237227560d8d98fe460c47ad0610830e18f9d25a811ee9ce41550cf7e",true)
		//提交交易
		ethService.ERC20Transaction(ToAddress,0.001511,true)
	}
	func main1(){
		//init()
		//var fromAddress="0x68f2e59b64a1b0b1786712fdc3374786f95aa055"
		//var ethHelp =heth.ETHHelp{}
		//nonce, err :=	ethHelp.GetNONCE(fromAddress)
		//if err!=nil {
		//	return
		//}
		//ethHelp.Nonce=nonce
		//pri,err:= ethHelp.GetPrivateKey(fromAddress)
		//if err!=nil {
		//	return
		//}
		//ethHelp.PrivateKey=pri
		//if err!=nil {
		//	return
		//}
		//gasData:=  ethHelp.GetGas()
		//ethHelp.GasData=gasData
		//if err!=nil {
		//	return
		//}
		//chainID, err := ethHelp.GetchainID()
		//ethHelp.ChainID=chainID
		//heth.AAA(nonce,pri,chainID)




// 0x64fe95228af4945d81087f81e9382e3623723f21
//0x99E46a01909078BC8D0A53CA8EFa2a0B7Ec1497c
		var  ToAddress= "0x99E46a01909078BC8D0A53CA8EFa2a0B7Ec1497c"

		//var ethService= service.ETHService{}
		var ethxx=heth.ETHHelp{}
			ethxx.GetETHData()
		//覆盖之前交易
		//	ethService.ERC20PaddingHander("0xb27f4b6237227560d8d98fe460c47ad0610830e18f9d25a811ee9ce41550cf7e",true)
		//提交交易
		//ethService.ERC20Transaction(ToAddress,0.00155,true)
		ethxx.CheckTransaction([]string{ToAddress})
	}
