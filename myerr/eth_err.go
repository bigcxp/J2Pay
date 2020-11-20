package myerr

import "strings"

//eth错误类
type EthError struct {
	s    string
	Code int
}

//判断错误类型
func (e *EthError) ErrorType(str string) (ethErr *EthError) {
	ethErr.s = str
	//nonce已经存在
	if strings.ContainsAny(str, "replacement transaction underpriced") {
		//表示nonce已存在
		ethErr.Code = 1
	}
	return
}

func (e *EthError) Error() string {
	return e.s
}
