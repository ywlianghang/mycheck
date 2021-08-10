package api
//
//import (
//	log "github.com/sirupsen/logrus"
//	"os"
//)
//
//type ResultOutputInterface interface {
//	//输出结果为log
//	OutputLog()
//	//输出结果为html
//	//输出结果为Exec
//}
//type ResultOutputStruct struct {
//}
//
//func (resultOutput *ResultOutputStruct) OutputLog() {
//	//loggs.SetFormatter(&loggs.JSONFormatter{})
//	log.SetFormatter(&log.TextFormatter{})
//	log.SetOutput(os.Stdout)
//	log.SetLevel(log.InfoLevel)
//	log.Debug("Useful debugging information.")
//	log.Info("Something noteworthy happened!")
//	log.Warn("You should probably take a look at this.")
//	log.Error("Something failed but I'm not quitting.")
//	//loggs.Fatal("Bye.")   //log之后会调用os.Exit(1)
//	//loggs.Panic("I'm bailing.")   //log之后会panic()
//	log.WithFields(log.Fields{
//		"animal": "walrus",
//	}).Info("A walrus appears")
//	entry := log.WithFields(log.Fields{"request_id":32,"user_ip":"127.0.0.1"})
//	entry.Info("something happened on that request")
//	entry.Warn("something not great happened")
//	log.Warn("aaa")
//
//
//}
//
