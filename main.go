package main

import (
	"DepthInspection/InspectionItem"
	"DepthInspection/api/PublicDB"
	"DepthInspection/api/loggs"
	"fmt"
)


func ConfigInit() *PublicDB.ConfigInfo {
	info := loggs.BaseInfo{}
	conf := info.GetConf()
	connInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",conf.DBinfo.Username,conf.DBinfo.Password,
		conf.DBinfo.Host,conf.DBinfo.Port,conf.DBinfo.Database,conf.DBinfo.Charset)
	var logconfig = &loggs.LogStruct{
		LoggLevel: conf.Logs.Loglevel,
		Logfile: conf.Logs.OutputFile.Logfile,
		Skip: conf.Logs.OutputFile.Skip,
		RotationTime: conf.Logs.OutputFile.RotationTime,
		IsConsole: conf.Logs.OutputFile.IsConsole,
		LogMaxAge: conf.Logs.OutputFile.LogMaxAge,
	}
	var dbconfig = &PublicDB.DatabaseExecStruct{
		MaxIdleConns: conf.DBinfo.MaxIdleConns,
		DirverName: conf.DBinfo.DirverName,
		DBconnIdleTime: conf.DBinfo.DBconnIdleTime,
		ConnInfo: connInfo,
	}
	var confaa = &PublicDB.ConfigInfo{
		DatabaseExecInterf: dbconfig,
		Loggs: logconfig,
	}
	return confaa

}

func main() {
	aa := ConfigInit()
	var ccc = map[string]string {
		"super_read_only" : "off",
		"read_only" : "off",
		"innodb_read_only": "off",
		"binlog_format" : "row",
		"character_set_server" : "utf8",
		"default_authentication_plugin" : "mysql_native_password",
		"default_storage_engine" : "innodb",
		"default_tmp_storage_engine" : "innodb",
		"innodb_flush_log_at_trx_commit" : "1",
		"innodb_flush_method" : "O_DIRECT",
		"innodb_deadlock_detect" : "on",
		"internal_tmp_disk_storage_engine" : "innodb",
		"query_cache_type" : "off",
		"relay_log_purge" : "on",
		"relay_log_recovery" : "on",
		"sync_binlog" : "1",
		"system_time_zone" : "CST",
		"time_zone" : "system",
		"transaction_isolation" : "READ-COMMITTED",
		"transaction_read_only" : "off",
		"tx_isolation" : "READ-COMMITTED",
		"tx_read_only" : "off",
		"unique_checks" : "on",
	}

	InspectionItem.DBConfigCheck(aa,ccc)
	var c = &InspectionItem.DatabaseBaselineCheckStruct{}
	//c.TableDesignCompliance(aa)

	c.BaselineCheckTablesDesign(aa)
	c.BaselineCheckColumnsDesign(aa)
}

	//var e api.ResultOutputInterface
	//e = &api.ResultOutputStruct{}
	//e.OutputLog()
