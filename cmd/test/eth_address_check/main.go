// 检测eth剩余可用地址是否满足需求，
// 如果不足则创建地址
package main

import (
	"j2pay-server/heth"
)

func main() {

	heth.CheckAddressFree()
}
