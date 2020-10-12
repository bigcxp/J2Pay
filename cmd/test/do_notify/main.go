// 发送通知
package main

import (
	"j2pay-server/app"

)

func main() {
	xenv.EnvCreate()
	defer xenv.EnvDestroy()

	app.CheckDoNotify()
}
