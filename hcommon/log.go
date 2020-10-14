package hcommon

import (
	"log"

	"go.uber.org/zap"
)

// Log 日志对象
var Log LoggerAble
var ZapLog *zap.Logger

type LoggerAble interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

func init() {
	// 初始化默认的日志对象
	var err error
	conf := zap.NewDevelopmentConfig()
	conf.DisableStacktrace = true
	conf.Encoding = "console"
	ZapLog, err = conf.Build()
	if err != nil {
		log.Fatalf("build logger error: [%T] %s", err, err.Error())
	}
	Log = ZapLog.Sugar()
}
