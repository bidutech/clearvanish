package LogerWraper

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

///////////////////////////////////////

/**************调用顺序很重要，必须先创建实例
mysqlloger.SetRollingDaily("/usr/local/test/golangLogertest", "test.log")
mysqlloger.SetLevel(LogerWraper.ALL)
mysqlloger.SetDebug(true)
mysqlloger.SetConsole(true)
mysqlloger.Info(i, 123, "+++","end")
*/

type LogConfig struct {
	LogFileName  string `json:"logname"`
	LogFilePath  string `json:"logpath"`
	Debug        int    `json:"debug"`
	ConsolePrint int    `json:"console"`
	Level        LEVEL  `json:"level"`
	MaxNumber    int32  `json:"maxnumber"`
	MaxSize      int64  `json:"maxsize"` //默认单位MB
	Type         int    `json:"type"`    // 1 dailylog   2 filelog
}

func (l *LogerWraper) InitLog(logConfigName string) {
	l.realnitLogFromconfigfile(logConfigName)
}

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

func (l *LogerWraper) realnitLogFromconfigfile(logConfigName string) {
	l.initLoger()

	var LogConf LogConfig
	l.initLogConfig(logConfigName, &LogConf)
	fmt.Println(LogConf)
	level := LogConf.Level
	if level < ALL {
		level = ALL
	} else if level > OFF {
		level = OFF
	}
	var debug bool
	if LogConf.Debug == 0 {
		debug = false
	} else {
		debug = true
	}

	var console bool
	if LogConf.ConsolePrint == 0 {
		console = false
	} else {
		console = true
	}

	l.mkdir(LogConf.LogFilePath)
	var daily bool
	daily = false
	if LogConf.Type == 1 {
		daily = true
		l.initRollingDaily(LogConf.LogFileName, LogConf.LogFilePath, debug, console, level)
	} else if LogConf.Type == 2 {
		l.initRollingFile(LogConf.LogFileName, LogConf.LogFilePath, debug, console, level,
			LogConf.MaxNumber, LogConf.MaxSize, MB)
	} else {
		fmt.Println("ERR Log Type")
	}

	l.lcfg.logLevel = level
	l.lcfg.dailyRolling = daily
	l.lcfg.consoleAppender = console
	l.lcfg.RollingFile = daily
	l.dodebug = debug

}

func (l *LogerWraper) isDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

func (l *LogerWraper) mkdir(path string) {
	if !l.isDirExists(path) {
		err := os.MkdirAll(path, 0777)
		if err != nil {
			fmt.Printf("ERR InitLogInfo:%s ,config path:%s", err, path)
			os.Exit(-1)

		} else {
			fmt.Print("Create Directory OK!")
		}
	}
}

func (l *LogerWraper) initLogConfig(configpath string, conf *LogConfig) {
	f, err := os.Open(configpath)
	defer f.Close()
	if nil == err {
		buff := bufio.NewReader(f)
		for {
			line, err := buff.ReadBytes('\n')
			if err != nil || io.EOF == err {
				return
			}
			errjson := json.Unmarshal(line, conf)
			if errjson != nil {
				fmt.Printf("ERR InitConfig Configsr.DataPath:%s,line:%s\n", configpath, line)
				os.Exit(-1)
			}
			break
		}
	} else {
		fmt.Printf("read config error-2")
		os.Exit(-1)
	}
}

func (l *LogerWraper) initRollingDaily(logFileName, logFilePath string, debug bool, consolePrint bool, _level LEVEL) {

	l.SetRollingDaily(logFilePath, logFileName)
	l.SetLevel(_level)
	//指定是否控制台打印，默认为true
	l.SetConsole(consolePrint)
	l.SetDebug(debug)
}

func (l *LogerWraper) initRollingFile(logFileName, logFilePath string,
	debug bool,
	consolePrint bool, _level LEVEL,
	maxNumber int32, maxSize int64, _unit UNIT /*logger.KB GB*/) {

	l.SetRollingFile(logFilePath, logFileName, maxNumber, maxSize, _unit)
	//指定是否控制台打印，默认为true
	l.SetConsole(consolePrint)
	l.SetDebug(debug)
	l.SetLevel(_level)
}
