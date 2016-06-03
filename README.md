##LogerWraper一个更好的golang日志系统
       是golang 的日志库 ，是对golang内置log的封装，基于donnie4w/go-logger开发
#donnie4w/go-logger存在的问题：
###1、对于多参数log使用[ ]层层包装，很乱
###2、不支持多实例创建，全局使用一个log实例很有局限性分类日志（多文件），所有log写到同一个文件，这个无法接受。
#LogerWraper改进
###1、对于多参数去掉了[ ]层层包装，只保留最外层的[ ]
###2、支持根据业务类型分别创建LogerWraper实例将日志写入到不同文件
###3、支持Debug模式，这个模式可以输出详细的日志调用 函数 文件名 和行号 快速定位问题
###4、注意问题：
           必须先调用 SetRollingDaily 或 SetRollingDaily 先创建实例才能再对日志进行设置各项Set操作
###4、example:
* 配置文件方式（推荐）
<pre><code>
需要按照log.cfg模板进行配置

logpath  日志存储路径

其中debug  取值 0  不打印调用代码信息  非零 打印调用代码信息
console  0  不在控制台打印 非零在控制台打印 日志信息
type 1  按天分割  0 按文件大小分割
level （0-6）的整数{

	ALL LEVEL = iota //0
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	OFF //6
}

注意：所有值必须填写
</code></pre>

const (
	LogConfigFile = "log.cfg"
)

var loger LogerWraper.LogerWraper

func InitLog() {
	logConfigFile := "logpath" + LogConfigFile
	loger.InitLog(logConfigFile)
	loger.Info("Loger ready")

}

* 原生方式麻烦（不推荐）
{
		
        var mysqlloger LogerWraper.LogerWraper
		mysqlloger.SetRollingDaily("/usr/local/test/mysqlloger", "mysqlloger.log") //必须首先调用
		//mysqlloger.SetRollingFile //必须首先调用
		mysqlloger.SetLevel(LogerWraper.ALL)
		mysqlloger.SetDebug(true)
		mysqlloger.SetConsole(true)
		mysqlloger.Info(i, 123, "hello LogerWraper", "LogerWraper easy")


		var MongoDBlloger LogerWraper.LogerWraper

		MongoDBlloger.SetRollingDaily("/usr/local/test/MongoDBlloger", "MongoDBlloger.log") //必须首先调用
		//mysqlloger.SetRollingFile //必须首先调用
		MongoDBlloger.SetLevel(LogerWraper.ALL)
		MongoDBlloger.SetDebug(true)
		MongoDBlloger.SetConsole(true)
		MongoDBlloger.Info(i, 123, "hello LogerWraper", "LogerWraper easy")
	}
###5、说明：
       用法类似java日志工具包log4j
打印日志有5个方法 Debug，Info，Warn, Error ,Fatal  日志级别由低到高
设置日志级别的方法为：LogerWraper.SetLevel() 如：LogerWraper.SetLevel(logger.WARN)
则：LogerWraper.Debug(....),LogerWraper.Info(...) 日志不会打出，而 
 LogerWraper.Warn(...),LogerWraper.Error(...),loLogerWraper.Fatal(...)日志会打出。
设置日志级别的参数有7个，分别为：ALL，DEBUG，INFO，WARN，ERROR，FATAL，OFF
其中 ALL表示所有调用打印日志的方法都会打出，而OFF则表示都不会打出。


日志文件切割有两种类型：1为按日期切分。2为按日志大小切分。
按日期切分时：每天一个备份日志文件，后缀为 .yyyy-MM-dd 
过0点是生成前一天备份文件

按大小切分是需要3个参数，1为文件大小，2为单位，3为文件数量
文件增长到指定限值时，生成备份文件，结尾为依次递增的自然数。
文件数量增长到指定限制时，新生成的日志文件将覆盖前面生成的同名的备份日志文件。

示例：

	var mysqlloger LogerWraper.LogerWraper

	//指定日志文件备份方式为文件大小的方式
	//第一个参数为日志文件存放目录
	//第二个参数为日志文件命名
	//第三个参数为备份文件最大数量
	//第四个参数为备份文件大小
	//第五个参数为文件大小的单位 KB，MB，GB TB
	//mysqlloger.SetRollingFile("d:/logtest", "test.log", 10, 5, logger.KB)

    //或
	//指定日志文件备份方式为日期的方式
	//第一个参数为日志文件存放目录
	//第二个参数为日志文件命名
	//mysqlloger.SetRollingDaily("d:/logtest", "test.log")

	//指定日志级别  ALL，DEBUG，INFO，WARN，ERROR，FATAL，OFF 级别由低到高
	//一般习惯是测试阶段为debug，生成环境为info以上
	//mysqlloger.SetLevel(LogerWraper.DEBUG)

		mysqlloger.SetRollingDaily("/usr/local/test/golangLogertest", "logtest.log") //必须首先调用
        //指定是否控制台打印，默认为true
		//mysqlloger.SetRollingFile //必须首先调用
		mysqlloger.SetLevel(LogerWraper.ALL)
		mysqlloger.SetDebug(true)
		mysqlloger.SetConsole(true)
		mysqlloger.Info(i, 123, "######———", "++++++++++++")


