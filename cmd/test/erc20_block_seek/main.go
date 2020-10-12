// 检索erc20到账情况
package main

import (
	"j2pay-server/heth"
	"j2pay-server/xenv"
)

func main() {

	xenv.EnvCreate()
	defer xenv.EnvDestroy()
	heth.CheckErc20BlockSeek()
}
