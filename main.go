// @title y2pay
package main

import (
	"flag"
	"fmt"
	_ "github.com/ethereum/go-ethereum/accounts/keystore"
	"j2pay-server/ethclient"
	"j2pay-server/pkg/logger"
	"j2pay-server/pkg/setting"
	"j2pay-server/routers"
	"j2pay-server/xenv"
)

func main() {
	flag.Parse()

	// 初始化操作 (因为 init 方法无法保证我们想要的顺序)
	setting.Setup()
	//以太坊节点
	ethclient.InitClient(fmt.Sprintf("%s", setting.EthConf.Url))
	logger.Setup()
	xenv.Setup()
	router := routers.InitRouter()
	//client := util.EthClient()
	////util.GetUserAddress("123456")
	//accounts, err := client.EthAccounts()
	//if err != nil{
	//	return
	//}
	//fmt.Println(accounts)
	//balance, err := client.EthGetBalance("0x3305a26A87bc4Cdb761e7623bf7054EA8376863b", "pending")
	//fmt.Println(balance)
	//
	//count, err := client.NetPeerCount()
	//fmt.Println(count)

	//secret := validate.NewGoogleAuth().GetSecret()
	//code, err := validate.NewGoogleAuth().GetCode(secret)
	//qrcode := validate.NewGoogleAuth().GetQrcode("admin", secret)
	//fmt.Println(qrcode)
	//url := validate.NewGoogleAuth().GetQrcodeUrl("admin", secret)
	//fmt.Println(url)
	//fmt.Println(secret, code, err)
	//verifyCode, err := validate.NewGoogleAuth().VerifyCode("AAXZBXISLVBR65RCCRIMDBLHZEUVIPUR", "997784")
	//fmt.Println(verifyCode)

	panic(router.Run(fmt.Sprintf("%s:%d", setting.ApplicationConf.Host, setting.ApplicationConf.Port)))
}
