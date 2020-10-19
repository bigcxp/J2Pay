// @title y2pay
package main

import (
	"flag"
	"fmt"
	_ "github.com/ethereum/go-ethereum/accounts/keystore"
	"j2pay-server/ethclient"
	"j2pay-server/heth"
	"j2pay-server/pkg/setting"
	"j2pay-server/routers"
)

func main() {
	//把用户传递的命令行参数解析为对应变量的值
	flag.Parse()
	// 初始化操作 (因为 init 方法无法保证我们想要的顺序)
	setting.Setup()
	//日志
//	logger.Setup()
	//初始化以太坊节点
	ethclient.InitClient(fmt.Sprintf("%s", setting.EthConf.Url))
	//生成热钱包地址
	//address, err := heth.CreateHotAddress(1)
	//if err != nil {
	//	return
	//}
	//fmt.Println(address)
	//检测充币
	//heth.CheckBlockSeek()
	//发送交易
	heth.CheckRawTxSend()
	//确认tx是否打包
	heth.CheckRawTxConfirm()
	//零钱整理到冷钱包
	//heth.CheckAddressOrg()

	//生成备用地址
	//free, err := heth.CheckAddressFree()
	//if err != nil {
	//	return
	//}
	//fmt.Println(free)
	//网关
	router := routers.InitRouter()
	//启动
	panic(router.Run(fmt.Sprintf("%s:%d", setting.ApplicationConf.Host, setting.ApplicationConf.Port)))
}
