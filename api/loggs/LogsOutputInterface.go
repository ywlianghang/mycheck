package loggs

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
	"time"
)

type LogStruct struct {
	Logfile string             //日志文件
	LogMaxAge time.Duration    //日志最大生存时间
	RotationTime time.Duration  //日志轮询时间
	IsConsole bool  //loggs formatter 是否使用颜色 true使用，false不使用
	LoggLevel string   //输出的日志级别
	Skip   int
}

type LogOutputInterface interface {
	configLogOutputformatter() *nested.Formatter
	configLocalFilesystemLogger() logrus.Hook
	logOutputFile(logglevel ,logInfo interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
}

func (logg *LogStruct) Info(args ...interface{}){
	logg.logOutputFile("info",args)
}
func (logg *LogStruct) Debug(args ...interface{}){
	logg.logOutputFile("debug",args)
}
func (logg *LogStruct) Warn(args ...interface{}){
	logg.logOutputFile("warning",args)
}
func (logg *LogStruct) Error(args ...interface{}){
	logg.logOutputFile("error",args)
}
func (logg *LogStruct) Fatal(args ...interface{}){
	logg.logOutputFile("fatal",args)
}
func (logg *LogStruct) Panic(args ...interface{}){
	logg.logOutputFile("panic",args)
}

//日志信息输出到文件，格式为text
func (logg *LogStruct) logOutputFile(logglevel ,logInfo interface{}) {
	log := logrus.New()
	src, err := os.OpenFile(logg.Logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND,os.ModeAppend|os.ModePerm)
	if err != nil {
		fmt.Println("err", err)
	}
	log.Out = src
	hook := logg.configLocalFilesystemLogger()
	log.AddHook(hook)
	log.SetFormatter(logg.configLogOutputformatter())
	log.SetReportCaller(true)
	if logg.LoggLevel == "debug"{
		log.SetLevel(logrus.DebugLevel)
	}else if logg.LoggLevel == "info"{
		log.SetLevel(logrus.InfoLevel)
	} else if logg.LoggLevel == "warning"{
		log.SetLevel(logrus.WarnLevel)
	} else if logg.LoggLevel == "error"{
		log.SetLevel(logrus.ErrorLevel)
	}
	switch logglevel {
		case "debug":
			log.Debug(logInfo)
		case "info":
			log.Info(logInfo)
		case "warning":
			log.Warn(logInfo)
		case "error":
			log.Error(logInfo)
		case "fatal":
			log.Fatal()
		case "panic":
			log.Panic(logInfo)
	}
}

//配置日志输出的格式
func (logg *LogStruct) configLogOutputformatter() *nested.Formatter {
	fmtter := &nested.Formatter{
		HideKeys:        true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerFirst:     true,
		CustomCallerFormatter: func(frame *runtime.Frame) string {
			file := ""
			var line int
			var ok bool
			for i := 0; i < 5; i++ {
				_, file, line, ok = runtime.Caller(logg.Skip + i )
				if !ok {
					return ""
				}
				n := 0
				for i := len(file) - 1; i > 0; i-- {
					if file[i] == '/' {
						n++
						if n >= 2 {
							file = file[i+1:]
							break
						}
					}
				}
				if !strings.HasPrefix(file, "logrus") {
					break
				}
			}
			return fmt.Sprintf(" [%s:%d] ", file, line)
		},
	}
		if logg.IsConsole{fmtter.NoColors = false} else{fmtter.NoColors = true}
		return fmtter
}



//日志信息输出到文件，配置日志轮转方式及格式为text
func (logg *LogStruct) configLocalFilesystemLogger() logrus.Hook {
	writer, err := rotatelogs.New(
	/* 日志轮转相关函数
	`WithLinkName (baseLogPaht）` 为最新的日志建立软连接
	`WithRotationTime` 设置日志分割的时间，隔多久分割一次
	 WithMaxAge 和 WithRotationCount二者只能设置一个
	`WithMaxAge (time.Second*60*3)` 设置文件清理前的最长保存时间
	`WithRotationCount (time.Second*60)` 设置文件清理前最多保存的个数
	*/
	logg.Logfile+".%Y%m%d%H%M",
	rotatelogs.WithMaxAge(logg.LogMaxAge), // 文件最大保存时间
	rotatelogs.WithRotationTime(logg.RotationTime), // 日志切割时间间隔
	)
	if err != nil {
		logrus.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	},logg.configLogOutputformatter())
	return lfHook
}




// config logrus loggs to amqp
//func (logg *LogStruct) LogOutputMq(server, username, password, exchange, exchangeType, virtualHost, routingKey string) {
//	hook := logrus_amqp.NewAMQPHookWithType(server, username, password, exchange, exchangeType, virtualHost, routingKey)
//	logrus.AddHook(hook)
//}

// config logrus loggs to es
//func (logg *LogStruct) LogOutputES(esUrl string, esHOst string, index string) {
//	client, err := elastic.NewClient(elastic.SetURL(esUrl))
//	if err != nil {
//		logrus.Errorf("config es logger error. %+v", errors.WithStack(err))
//	}
//	esHook, err := elogrus.NewElasticHook(client, esHOst, logrus.DebugLevel, index)
//	if err != nil {
//		logrus.Errorf("config es logger error. %+v", errors.WithStack(err))
//	}
//	logrus.AddHook(esHook)
//}

//func init(){
//	var f LogOutputInterface
//	f = &LogStruct{
//		Logfile: "D:\\goProject\\DepthInspection\\aab",
//		LoggLevel: "info",
//		LogMaxAge: 100,
//		IsConsole: false,
//		RotationTime: 100,
//		Skip: 10,
//	}
//}


