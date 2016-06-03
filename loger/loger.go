package loger

import (
	"clearvanish/config"

	"github.com/shanhai2015/LogerWraper"
)

const (
	LogConfigFile = "log.cfg"
)

var Loger LogerWraper.LogerWraper

func InitLog() {
	logConfigFile := config.GetCurrentPath() + LogConfigFile
	Loger.InitLog(logConfigFile)
	Loger.Info("Loger ready")

}
