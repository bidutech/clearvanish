package loger

import (
	"fmt"

	"github.com/shanhai2015/LogerWraper"
	"github.com/shanhai2015/SHcommon"
)

const (
	LogConfigFile = "log.cfg"
)

var Loger LogerWraper.LogerWraper

func init() {
	logConfigFile := SHcommon.GetCurrentPath() + LogConfigFile
	Loger.InitLog(logConfigFile)
	Loger.Info("Loger ready")
}

func ErrorLog(err error) {
	Loger.Error(fmt.Sprintf("%v", err))
}
