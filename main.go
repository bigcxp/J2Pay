// @title y2pay
package main

import (
	"flag"
	"fmt"
	"j2pay-server/model"
	"j2pay-server/pkg/logger"
	"j2pay-server/pkg/setting"
	"j2pay-server/routers"
)

func main() {
	flag.Parse()

	// 初始化操作 (因为 init 方法无法保证我们想要的顺序)
	setting.Setup()
	logger.Setup()
	model.Setup()

	router := routers.InitRouter()
	panic(router.Run(fmt.Sprintf("%s:%d", setting.ApplicationConf.Host, setting.ApplicationConf.Port)))
}