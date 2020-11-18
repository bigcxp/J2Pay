package util

import "regexp"

// 验证一个输入是不是IP地址
func IsIp(ip string) bool {
	if m, _ := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", ip); !m {
		return false
	}
	return true
}

// 验证一个是不是http请求
func IsHttp(url string) bool {
	if m, _ := regexp.MatchString("^[a-zA-Z]+://(\\w+(-\\w+)*)(\\.(\\w+(-\\w+)*))*(\\?\\s*)?$", url); !m {
		return false
	}
	return true
}
