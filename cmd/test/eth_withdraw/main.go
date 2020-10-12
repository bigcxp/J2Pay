// 检测提现
package main

import (
	"j2pay-server/heth"
	"j2pay-server/xenv"
)

func main() {
	xenv.EnvCreate()
	defer xenv.EnvDestroy()
	heth.CheckWithdraw()
}
