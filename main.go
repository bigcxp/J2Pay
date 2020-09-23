// @title y2pay
package main

import (
	"flag"
	"fmt"
	_"github.com/ethereum/go-ethereum/accounts/keystore"
	"j2pay-server/model"
	"j2pay-server/pkg/logger"
	"j2pay-server/pkg/setting"
	"j2pay-server/pkg/util"
	"j2pay-server/routers"
)

func main() {
	flag.Parse()

	// 初始化操作 (因为 init 方法无法保证我们想要的顺序)
	setting.Setup()
	logger.Setup()
	model.Setup()
	router := routers.InitRouter()
	client := util.EthClient()
	//util.GetUserAddress("123456")
	accounts, err := client.EthAccounts()
	if err != nil{
		return
	}
	fmt.Println(accounts)
	balance, err := client.EthGetBalance("0x3305a26A87bc4Cdb761e7623bf7054EA8376863b", "pending")
	fmt.Println(balance)

	count, err := client.NetPeerCount()
	fmt.Println(count)

	panic(router.Run(fmt.Sprintf("%s:%d", setting.ApplicationConf.Host, setting.ApplicationConf.Port)))
}
