package LogerWraper

import (
	"fmt"
	"log"
	"os"
	//"runtime"
	"strconv"
	"sync"
	"time"
)

/**************调用顺序很重要
mysqlloger.SetRollingDaily("/usr/local/test/golangLogertest", "test.log")
mysqlloger.SetLevel(LogerWraper.ALL)
mysqlloger.SetDebug(true)
mysqlloger.SetConsole(true)
mysqlloger.Info(i, 123, "+++","end")
*/

func (l *LogerWraper) SetRollingDaily(fileDir, fileName string) {
	l.initLoger()
	l.lcfg.RollingFile = false
	l.lcfg.dailyRolling = true
	t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
	l.logObj = &_FILE{dir: fileDir, filename: fileName, _date: &t, isCover: false, mu: new(sync.RWMutex)}
	l.logObj.mu.Lock()
	l.logObj.lcfg = l.lcfg
	l.logObj.logObj = l.logObj
	defer l.logObj.mu.Unlock()

	if !l.logObj.isMustRename() {
		l.logObj.logfile, _ = os.OpenFile(fileDir+"/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0)
		//l.logObj.lg = log.New(l.logObj.logfile, "\n", log.Ldate|log.Ltime|log.Lshortfile)
		l.logObj.lg = log.New(l.logObj.logfile, "\n", log.Ldate|log.Ltime)
	} else {
		l.logObj.rename()
	}
}

func (l *LogerWraper) SetRollingFile(fileDir, fileName string, maxNumber int32, maxSize int64, _unit UNIT) {
	l.initLoger()
	l.lcfg.maxFileCount = maxNumber
	l.lcfg.maxFileSize = maxSize * int64(_unit)
	l.lcfg.RollingFile = true
	l.lcfg.dailyRolling = false
	l.logObj = &_FILE{dir: fileDir, filename: fileName, isCover: false, mu: new(sync.RWMutex)}
	l.logObj.mu.Lock()
	l.logObj.lcfg = l.lcfg
	l.logObj.logObj = l.logObj

	defer l.logObj.mu.Unlock()
	for i := 1; i <= int(maxNumber); i++ {
		if isExist(fileDir + "/" + fileName + "." + strconv.Itoa(i)) {
			l.logObj._suffix = i
		} else {
			break
		}
	}
	if !l.logObj.isMustRename() {
		l.logObj.logfile, _ = os.OpenFile(fileDir+"/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0)
		//l.logObj.lg = log.New(l.logObj.logfile, "\n", log.Ldate|log.Ltime|log.Lshortfile)
		l.logObj.lg = log.New(l.logObj.logfile, "\n", log.Ldate|log.Ltime)
	} else {
		l.logObj.rename()
	}
	go l.fileMonitor()
}

func (l *LogerWraper) SetConsole(isConsole bool) {
	l.lcfg.consoleAppender = isConsole
}

func (l *LogerWraper) SetLevel(_level LEVEL) {
	l.lcfg.logLevel = _level
}

func (l *LogerWraper) SetDebug(debug bool) {
	l.dodebug = debug
}
func (l *LogerWraper) Debug(v ...interface{}) {
	msg := l.callerInfo() + fmt.Sprintf("%v", v)
	l.debug(msg)
}

func (l *LogerWraper) Info(v ...interface{}) {

	msg := l.callerInfo() + fmt.Sprintf("%v", v)
	l.info(msg)
}

func (l *LogerWraper) Warn(v ...interface{}) {

	msg := l.callerInfo() + fmt.Sprintf("%v", v)
	l.warn(msg)
}

func (l *LogerWraper) Error(v ...interface{}) {
	msg := l.callerInfo() + fmt.Sprintf("%v", v)
	l.error(msg)
}

func (l *LogerWraper) Fatal(v ...interface{}) {
	msg := l.callerInfo() + fmt.Sprintf("%v", v)
	l.fatal(msg)
}
