package LogerWraper

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	_VER string = "1.0.0"
)

type LEVEL int32

const DATEFORMAT = "2006-01-02"

type UNIT int64

const (
	_       = iota
	KB UNIT = 1 << (iota * 10)
	MB
	GB
	TB
)

const (
	ALL LEVEL = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

type LogerConfig struct {
	logLevel        LEVEL //=1
	maxFileSize     int64
	maxFileCount    int32
	dailyRolling    bool //true
	consoleAppender bool //true
	RollingFile     bool //false
}
type _FILE struct {
	dir      string
	filename string
	_suffix  int
	isCover  bool
	_date    *time.Time
	mu       *sync.RWMutex
	logfile  *os.File
	lg       *log.Logger
	lcfg     LogerConfig
	logObj   *_FILE
}

type LogerWraper struct {
	lcfg    LogerConfig
	logObj  *_FILE
	dodebug bool
}

func (l *LogerWraper) initLoger() {
	l.lcfg.logLevel = 1
	l.lcfg.dailyRolling = true
	l.lcfg.consoleAppender = true
	l.lcfg.RollingFile = false
	l.dodebug = false
}

func (l *LogerWraper) console(s ...interface{}) {
	if l.lcfg.consoleAppender {
		_, file, line, _ := runtime.Caller(2)
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		log.Println(file+":"+strconv.Itoa(line), s)
	}
}

func (l *LogerWraper) consoleBetter(s string) {
	if l.lcfg.consoleAppender {
		_, file, line, _ := runtime.Caller(2)
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		msg := fmt.Sprintf("%s:%s\t%s", file, strconv.Itoa(line), s)
		log.Println(msg)
	}
}

func catchError() {
	if err := recover(); err != nil {
		log.Println("err", err)
	}
}

func (l *LogerWraper) debug(v string) {
	if l.lcfg.dailyRolling {
		l.fileCheck()
	}
	defer catchError()
	l.logObj.mu.RLock()
	defer l.logObj.mu.RUnlock()

	if l.lcfg.logLevel <= DEBUG {
		debug := fmt.Sprintf("%-6s", "DEBUG")
		l.logObj.lg.Output(2, debug+v)
		l.consoleBetter(string(debug + v))
	}
}

func (l *LogerWraper) info(v string) {
	if l.lcfg.dailyRolling {
		l.fileCheck()
	}
	defer catchError()
	l.logObj.mu.RLock()
	defer l.logObj.mu.RUnlock()
	if l.lcfg.logLevel <= INFO {
		info := fmt.Sprintf("%-6s", "INFO")
		l.logObj.lg.Output(2, info+v)
		l.consoleBetter(string(info + v))
	}
}
func (l *LogerWraper) warn(v string) {
	if l.lcfg.dailyRolling {
		l.fileCheck()
	}
	defer catchError()
	l.logObj.mu.RLock()
	defer l.logObj.mu.RUnlock()
	if l.lcfg.logLevel <= WARN {
		warn := fmt.Sprintf("%-6s", "WARN")
		l.logObj.lg.Output(2, warn+v)
		l.consoleBetter(string(warn + v))
	}
}
func (l *LogerWraper) error(v string) {
	if l.lcfg.dailyRolling {
		l.fileCheck()
	}
	defer catchError()
	l.logObj.mu.RLock()
	defer l.logObj.mu.RUnlock()
	if l.lcfg.logLevel <= ERROR {
		error := fmt.Sprintf("%-6s", "ERROR")
		l.logObj.lg.Output(2, error+v)
		l.consoleBetter(string(error + v))
	}
}

func (l *LogerWraper) fatal(v string) {
	if l.lcfg.dailyRolling {
		l.fileCheck()
	}
	defer catchError()
	l.logObj.mu.RLock()
	defer l.logObj.mu.RUnlock()
	if l.lcfg.logLevel <= FATAL {
		fatal := fmt.Sprintf("%-6s", "FATAL")
		l.logObj.lg.Output(2, fatal+v)
		l.consoleBetter(string(fatal + v))
	}
}

func (f *_FILE) isMustRename() bool {
	if f.lcfg.dailyRolling {
		t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
		if t.After(*f._date) {
			return true
		}
	} else {
		if f.lcfg.maxFileCount > 1 {
			if fileSize(f.dir+"/"+f.filename) >= f.lcfg.maxFileSize {
				return true
			}
		}
	}
	return false
}

func (f *_FILE) rename() {
	if f.lcfg.dailyRolling {
		fn := f.dir + "/" + f.filename + "." + f._date.Format(DATEFORMAT)
		if !isExist(fn) && f.isMustRename() {
			if f.logfile != nil {
				f.logfile.Close()
			}
			err := os.Rename(f.dir+"/"+f.filename, fn)
			if err != nil {
				f.lg.Println("rename err", err.Error())
			}
			t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
			f._date = &t
			f.logObj.logfile, _ = os.Create(f.dir + "/" + f.filename)
			//f.lg = log.New(f.logObj.logfile, "\n", log.Ldate|log.Ltime|log.Lshortfile)
			f.lg = log.New(f.logObj.logfile, "\n", log.Ldate|log.Ltime)
		}
	} else {
		f.coverNextOne()
	}
}

func (f *_FILE) nextSuffix() int {
	return int(f._suffix%int(f.lcfg.maxFileCount) + 1)
}

func (f *_FILE) coverNextOne() {
	f._suffix = f.nextSuffix()
	if f.logfile != nil {
		f.logfile.Close()
	}
	if isExist(f.dir + "/" + f.filename + "." + strconv.Itoa(int(f._suffix))) {
		os.Remove(f.dir + "/" + f.filename + "." + strconv.Itoa(int(f._suffix)))
	}
	os.Rename(f.dir+"/"+f.filename, f.dir+"/"+f.filename+"."+strconv.Itoa(int(f._suffix)))
	f.logObj.logfile, _ = os.Create(f.dir + "/" + f.filename)
	//f.lg = log.New(f.logObj.logfile, "\n", log.Ldate|log.Ltime|log.Lshortfile)
	f.lg = log.New(f.logObj.logfile, "\n", log.Ldate|log.Ltime)
}

func fileSize(file string) int64 {
	fmt.Println("fileSize", file)
	f, e := os.Stat(file)
	if e != nil {
		fmt.Println(e.Error())
		return 0
	}
	return f.Size()
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func (l *LogerWraper) fileMonitor() {
	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer.C:
			l.fileCheck()
		}
	}
}

func (l *LogerWraper) fileCheck() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	if l.logObj != nil && l.logObj.isMustRename() {
		l.logObj.mu.Lock()
		defer l.logObj.mu.Unlock()
		l.logObj.rename()
	}
}
func (l *LogerWraper) callerInfo() string {
	callerInfos := ""
	if l.dodebug {
		funcname, filename, lines, ok := callerName(2) //设置为2 因为这是第二层调用
		if ok {
			callerInfos = fmt.Sprintf("%s %s %d ", funcname, filename, lines)
		}
	}
	return callerInfos
}
func callerName(skip int) (name, file string, line int, ok bool) { //skip=0
	var pc uintptr
	if pc, file, line, ok = runtime.Caller(skip + 1); !ok { //其中在执行 runtime.Caller 调用时, 参数 skip + 1 用于抵消 CallerName 函数自身的调用.
		return
	}
	name = runtime.FuncForPC(pc).Name()
	return
}
