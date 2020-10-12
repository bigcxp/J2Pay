// 检测eth到账
package main

import (
	"j2pay-server/heth"
)

func main() {
	xenv.EnvCreate()
	defer xenv.EnvDestroy()

	heth.CheckBlockSeek()
}
