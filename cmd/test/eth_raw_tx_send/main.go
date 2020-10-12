// 发送交易
package main

import (
	"j2pay-server/heth"
	"j2pay-server/xenv"
)

func main() {

	xenv.EnvCreate()
	defer xenv.EnvDestroy()
	heth.CheckRawTxSend()
}
