// @title y2pay
package main

import (
	"flag"
	"fmt"
	"j2pay-server/cron"
	"j2pay-server/ethclient"
	"j2pay-server/model"
	"j2pay-server/pkg/logger"
	"j2pay-server/pkg/setting"
	"j2pay-server/routers"
)

func main() {

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
	//定时器
	go cron.Cron()
	//网关
	router := routers.InitRouter()
	//启动
	panic(router.Run(fmt.Sprintf("%s:%d", setting.ApplicationConf.Host, setting.ApplicationConf.Port)))
}
