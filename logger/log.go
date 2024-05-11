package logger

import (
	"github.com/GuoFlight/gerror"
	"github.com/GuoFlight/glog"
	"github.com/sirupsen/logrus"
	"log"
	"openfly/conf"
)

var (
	GLogger *logrus.Logger
)

func InitLog() {
	path := conf.GConf.Log.Path
	logLevel := conf.GConf.Log.Level
	var err error
	GLogger, err = glog.NewLogger(path, logLevel, false, 10)
	if err != nil {
		log.Fatal("日志初始化失败:", err)
	}
	GLogger.Info("日志初始化完成")
}

// PrintErr 输出错误日志
func PrintErr(err *gerror.Gerr, elseInfo map[string]interface{}) *gerror.Gerr {
	if elseInfo == nil {
		elseInfo = make(map[string]interface{})
	}
	elseInfo["ErrFile"] = err.ErrFile
	elseInfo["ErrLine"] = err.ErrLine
	GLogger.WithFields(elseInfo).Error(err.Error())
	return err
}

// PrintWarn 输出Warn日志
func PrintWarn(err *gerror.Gerr, elseInfo map[string]interface{}) *gerror.Gerr {
	if elseInfo == nil {
		elseInfo = make(map[string]interface{})
	}
	elseInfo["ErrFile"] = err.ErrFile
	elseInfo["ErrLine"] = err.ErrLine
	GLogger.WithFields(elseInfo).Warn(err.Error())
	return err
}
