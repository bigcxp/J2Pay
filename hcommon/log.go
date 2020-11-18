package hcommon
import (
	"log"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

// Log 对外服务的日志对象
var Log LoggerAble

// ZapLog zap日志对象
var ZapLog *zap.Logger

// conf zap日志配置
var conf zap.Config

// LoggerAble 日志对象接口
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
	conf = zap.NewDevelopmentConfig()
	conf.DisableStacktrace = true
	conf.Encoding = "console"
	ZapLog, err = conf.Build()
	if err != nil {
		log.Fatalf("build logger error: [%T] %s", err, err.Error())
	}
	Log = ZapLog.Sugar()
}

// LogSetToProd 设置为生产环境
func LogSetToProd() error {
	var err error
	conf.Development = false
	conf.Encoding = "json"
	conf.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	ZapLog, err = conf.Build()
	if err != nil {
		return err
	}
	Log = ZapLog.Sugar()
	return nil
}

// LogSetLevel 设置日志等级
func LogSetLevel(level zapcore.Level) error {
	var err error
	conf.Level = zap.NewAtomicLevelAt(level)
	ZapLog, err = conf.Build()
	if err != nil {
		return err
	}
	Log = ZapLog.Sugar()
	return nil
}

