package util

import (
	"fmt"
	"time"
)

//时间字符串转时间类型
func TimeStr2Time(timeStr string) (time.Time){

	// 返回的是UTC时间 2020-01-13 22:32:17 +0000 UTC
	utcTimeObj, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err == nil {
		fmt.Println(utcTimeObj, utcTimeObj.Unix())
	}

	// 返回的是当地时间 2020-01-13 22:32:17 +0800 CST
	localTimeObj, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
	if err == nil {
		fmt.Println(localTimeObj)
	}
	return localTimeObj
}





