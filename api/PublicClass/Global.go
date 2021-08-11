package PublicClass

import (
	"DepthInspection/api/Stream"
	"DepthInspection/api/loggs"
	"fmt"
)

//一级结构体，数据库配置参数检查
type DatabaseConfigCheckResultStruct struct{
	ConfigParameter []map[string]string
}
//一级结构体，数据库性能检查
type DatabasePerformanceCheckStruct struct {
	PerformanceStatus DatabasePerformanceStatusCheckResultStruct  //检查状态
	PerformanceTableIndex DatabasePerformanceTableIndexCheckResultStruct  //检查表和索引
}

type BaselineCheckTablesDesignResultStruct struct{
	TableCharset []map[string]string
	TableEngine []map[string]string
	TableForeign []map[string]string
	TableNoPrimaryKey []map[string]string
}
type BaselineCheckColumnDesignResultStruct struct{
	TableAutoIncrement []map[string]string
	TableBigColumns []map[string]string
}
type BaselineCheckIndexColumnDesignResultStruct struct{
	IndexColumnIsNull []map[string]string
	IndexColumnIsEnumSet []map[string]string
	IndexColumnIsBlobText []map[string]string
	IndexColumnIsRepeatIndex []map[string]string
}
type BaselineCheckProcedureTriggerDesignResultStruct struct{
	TableProcedure []map[string]string
	TableTrigger []map[string]string
	TableFunc []map[string]string
}
type BaselineCheckUserPriDesignResultStruct struct{
	AnonymousUsers []map[string]string
	EmptyPasswordUser []map[string]string
	RootUserRemoteLogin []map[string]string
	NormalUserConnectionUnlimited []map[string]string
	UserPasswordSame []map[string]string
	NormalUserDatabaseAllPrivilages []map[string]string
	NormalUserSuperPrivilages []map[string]string
}
type BaselineCheckPortDesignResultStruct struct{
	DatabasePort []map[string]string
}
type DatabasePerformanceStatusCheckResultStruct struct{
	BinlogDiskUsageRate []map[string]string
	HistoryConnectionMaxUsageRate []map[string]string
	TmpDiskTableUsageRate []map[string]string
	TmpDiskfileUsageRate []map[string]string
	InnodbBufferPoolUsageRate []map[string]string
	InnodbBufferPoolDirtyPagesRate []map[string]string
	InnodbBufferPoolHitRate []map[string]string
	OpenFileUsageRate []map[string]string
	OpenTableCacheUsageRate []map[string]string
	OpenTableCacheOverflowsUsageRate []map[string]string
	SelectScanUsageRate []map[string]string
	SelectfullJoinScanUsageRate []map[string]string
}
type DatabasePerformanceTableIndexCheckResultStruct struct{
	TableAutoPrimaryKeyUsageRate []map[string]string
	TableRows []map[string]string
	DiskFragmentationRate []map[string]string
	BigTable []map[string]string
	ColdTable []map[string]string
}
type InspectionResultsStruct struct {
	DatabaseConfigCheck DatabaseConfigCheckResultStruct
	BaselineCheckTablesDesign BaselineCheckTablesDesignResultStruct
	BaselineCheckColumnsDesign BaselineCheckColumnDesignResultStruct
	BaselineCheckIndexColumnDesign BaselineCheckIndexColumnDesignResultStruct
	BaselineCheckProcedureTriggerDesign BaselineCheckProcedureTriggerDesignResultStruct
	BaselineCheckUserPriDesign BaselineCheckUserPriDesignResultStruct
	BaselineCheckPortDesign BaselineCheckPortDesignResultStruct
	DatabasePerformance DatabasePerformanceCheckStruct
	//DatabasePerformanceCheck DatabasePerformanceStatusCheckResultStruct
	//DatabasePerformanceTableIndexCheck DatabasePerformanceTableIndexCheckResultStruct
}

//配置文件相关结构体初始化
var Logconfig *loggs.LogStruct
var Info  *loggs.BaseInfo
var Dbconfig *DatabaseExecStruct
var Loggs  loggs.LogOutputInterface
var DBexecInter DatabaseOperation
var Strea  *Stream.StreamStruct
var ResultOutput *loggs.ResultOutputFileEntity

//巡检结果结构体初始化
//一级结构体初始化
var InspectionResult = &InspectionResultsStruct{}
var DatabasePerformance = &DatabasePerformanceCheckStruct{}
var DatabaseConfigCheckResult = &DatabaseConfigCheckResultStruct{}


var BaselineCheckTablesDesignResult = &BaselineCheckTablesDesignResultStruct{}
var BaselineCheckColumnDesignResult = &BaselineCheckColumnDesignResultStruct{}
var BaselineCheckIndexColumnDesignResult = &BaselineCheckIndexColumnDesignResultStruct{}
var BaselineCheckProcedureTriggerDesignResult = &BaselineCheckProcedureTriggerDesignResultStruct{}
var BaselineCheckPortDesignDesignResult = &BaselineCheckPortDesignResultStruct{}
var DatabasePerformanceStatusCheckResult = &DatabasePerformanceStatusCheckResultStruct{}
var DatabasePerformanceTableIndexCheckResult = &DatabasePerformanceTableIndexCheckResultStruct{}

//数据交互变量
var strSql string
var GlobalVariables = make(map[string]string)
var GlobalStatus = make(map[string]string)
var InformationSchemaTablesData, InformationSchemaColumnsData,InformationSchemaCollationsData  []map[string]interface{}
var InformationSchemaKeyColumnUsage,InformationSchemaStatistics,InformationSchemaRoutines []map[string]interface{}
var InformationSchemaTriggers,MysqlUser []map[string]interface{}

//检测耗时相关变量
var CheckBeginTime,CheckEndTime string
var CheckTimeConsuming int64

func ConfigInit() {
	//读取配置文件
	getConf := Info.GetConf()
	dbConnInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",getConf.DBinfo.Username,getConf.DBinfo.Password,
		getConf.DBinfo.Host,getConf.DBinfo.Port,getConf.DBinfo.Database,getConf.DBinfo.Charset)
	Logconfig = &loggs.LogStruct{
		LoggLevel: getConf.Logs.Loglevel,
		Logfile: getConf.Logs.OutputFile.Logfile,
		Skip: getConf.Logs.OutputFile.Skip,
		RotationTime: getConf.Logs.OutputFile.RotationTime,
		IsConsole: getConf.Logs.OutputFile.IsConsole,
		LogMaxAge: getConf.Logs.OutputFile.LogMaxAge,
	}
	Dbconfig = &DatabaseExecStruct{
		MaxIdleConns: getConf.DBinfo.MaxIdleConns,
		DirverName: getConf.DBinfo.DirverName,
		DBconnIdleTime: getConf.DBinfo.DBconnIdleTime,
		ConnInfo: dbConnInfo,
	}
	ResultOutput = &loggs.ResultOutputFileEntity{
		OutputWay: getConf.ResultOutput.OutputWay,
		OutputPath: getConf.ResultOutput.OutputPath,
		OutputFile: getConf.ResultOutput.OutputFile,
		InspectionPersonnel: getConf.ResultOutput.InspectionPersonnel,
		InspectionLevel: getConf.ResultOutput.InspectionLevel,
	}
	Loggs = Logconfig
	DBexecInter = Dbconfig
	Strea = &Stream.StreamStruct{}
}

func QueryDbDateInit(){
	strSql = "show global variables"
	a := DBexecInter.DBQueryDateMap(strSql)
	GlobalVariables = a
	ignoreTableSchema := "'mysql','information_schema','performance_schema','sys'"
	//
	strSql = fmt.Sprintf("show global status")
	b := DBexecInter.DBQueryDateMap(strSql)
	GlobalStatus = b
	strSql = fmt.Sprintf("select * from information_schema.tables where table_schema not in (%s)",ignoreTableSchema)
	InformationSchemaTablesData = DBexecInter.DBQueryDateJson(strSql)
	strSql = fmt.Sprintf("select * from information_schema.columns where table_schema not in (%s)",ignoreTableSchema)
	InformationSchemaColumnsData = DBexecInter.DBQueryDateJson(strSql)
	strSql = fmt.Sprintf("select * from information_schema.COLLATIONS")
	InformationSchemaCollationsData = DBexecInter.DBQueryDateJson(strSql)
	strSql = fmt.Sprintf("select CONSTRAINT_SCHEMA databaseName,TABLE_NAME tableName,COLUMN_NAME columnName,CONSTRAINT_NAME, REFERENCED_TABLE_NAME,REFERENCED_COLUMN_NAME from INFORMATION_SCHEMA.KEY_COLUMN_USAGE where CONSTRAINT_SCHEMA not in (%s)",ignoreTableSchema)
	InformationSchemaKeyColumnUsage = DBexecInter.DBQueryDateJson(strSql)
	strSql = fmt.Sprintf("select * from information_schema.STATISTICS where table_schema not in (%s)",ignoreTableSchema)
	InformationSchemaStatistics = DBexecInter.DBQueryDateJson(strSql)
	strSql = fmt.Sprintf("select ROUTINE_SCHEMA,ROUTINE_NAME,ROUTINE_TYPE,DEFINER,CREATED from information_schema.routines where ROUTINE_SCHEMA not in(%s)",ignoreTableSchema)
	InformationSchemaRoutines = DBexecInter.DBQueryDateJson(strSql)
	strSql = fmt.Sprintf("select TRIGGER_SCHEMA,TRIGGER_NAME,DEFINER,CREATED from information_schema.TRIGGERS where TRIGGER_SCHEMA not in (%s)",ignoreTableSchema)
	InformationSchemaTriggers = DBexecInter.DBQueryDateJson(strSql)
	strSql = fmt.Sprintf("select user,host,authentication_string password from mysql.user")
	MysqlUser = DBexecInter.DBQueryDateJson(strSql)

}
var Ccc = map[string]string {
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