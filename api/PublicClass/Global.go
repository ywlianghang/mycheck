package PublicClass

import (
	"DepthInspection/api/Stream"
	"DepthInspection/api/loggs"
	"fmt"
	"os"
	"runtime"
	"strings"
)

//顶级目录
type InspectionResultsStruct struct {
	DatabaseConfigCheck DatabaseConfigCheckResultStruct
	DatabaseBaselineCheck DatabaseBaselineCheckStruct
	DatabasePerformance DatabasePerformanceCheckStruct
	DatabaseSecurity  DatabaseSecurityCheckStruct
}

//一级结构体，数据库配置参数检查
type DatabaseConfigCheckResultStruct struct{
	ConfigParameter []map[string]string
}
//一级结构体，数据库性能检查
type DatabasePerformanceCheckStruct struct {
	PerformanceStatus DatabasePerformanceStatusCheckResultStruct  //检查状态
	PerformanceTableIndex DatabasePerformanceTableIndexCheckResultStruct  //检查表和索引
}
//一级结构体，数据库基线检查
type DatabaseBaselineCheckStruct struct{
	TableDesign BaselineCheckTablesDesignResultStruct
	ColumnDesign BaselineCheckColumnDesignResultStruct
	IndexColumnsDesign BaselineCheckIndexColumnDesignResultStruct
	ProcedureTriggerDesign BaselineCheckProcedureTriggerDesignResultStruct
}
//一级结构体，数据库安全检查
type DatabaseSecurityCheckStruct struct {
	UserPriDesign BaselineCheckUserPriDesignResultStruct
	PortDesign BaselineCheckPortDesignResultStruct
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
	IndexColumnType []map[string]string
	IndexColumnIsNull []map[string]string
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


//配置文件相关结构体初始化
var Logconfig *loggs.LogStruct
var Info  *loggs.BaseInfo
var Dbconfig *DatabaseExecStruct
var Loggs  loggs.LogOutputInterface
var DBexecInter DatabaseOperation
var Strea  *Stream.StreamStruct
var ResultOutput *loggs.ResultOutputFileEntity
var InspectionConfSwitch *loggs.InspectionConfSwitchFileEntity
var InspectionConfInput *loggs.InspectionConfInputFileEntity

//巡检结果结构体初始化
//一级结构体初始化
var InspectionResult = &InspectionResultsStruct{}
var DatabaseBaselineCheckResult = &DatabaseBaselineCheckStruct{}
var DatabasePerformance = &DatabasePerformanceCheckStruct{}
var DatabaseConfigCheckResult = &DatabaseConfigCheckResultStruct{}
var DatabaseSecurityResult = &DatabaseSecurityCheckStruct{}

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

func PathExists(path string){
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if strings.Contains(path,".\\") || strings.Contains(path,"./"){
			if pathdir,err2 := os.Getwd();err2 ==nil{
				tmppath := pathdir + path[1:]
				sysType := runtime.GOOS
				var aa int
				if sysType == "linux"{
					aa = strings.LastIndex(tmppath,"/")
				}
				if sysType == "windows"{
					aa = strings.LastIndex(tmppath,"\\")
				}
				path2 :=  tmppath[:aa]
				err1 := os.MkdirAll(path2, os.ModePerm)
				if err1 != nil{
					fmt.Println(err1)
					os.Exit(1)
				}
			}
		}
	}
}

func ConfigInit() {
	//读取配置文件
	getConf := Info.GetConf()
	dbConnInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",getConf.DBinfo.Username,getConf.DBinfo.Password,
		getConf.DBinfo.Host,getConf.DBinfo.Port,getConf.DBinfo.Database,getConf.DBinfo.Charset)
	PathExists(getConf.Logs.OutputFile.Logfile)
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
	InspectionConfSwitch = &loggs.InspectionConfSwitchFileEntity{
		ConfigSwitch: getConf.InspectionConfSwitch.ConfigSwitch,
	}
	InspectionConfInput = &loggs.InspectionConfInputFileEntity{
		DatabaseEnvironment: getConf.InspectionConfInput.DatabaseEnvironment,
		DatabaseConfiguration: getConf.InspectionConfInput.DatabaseConfiguration,
		DatabasePerformance: getConf.InspectionConfInput.DatabasePerformance,
		DatabaseBaseline: getConf.InspectionConfInput.DatabaseBaseline,
		DatabaseSecurity: getConf.InspectionConfInput.DatabaseSecurity,
		DatabaseSpace: getConf.InspectionConfInput.DatabaseSpace,
	}
	Loggs = Logconfig
	DBexecInter = Dbconfig
	Strea = &Stream.StreamStruct{}
}
var (
	ConfigurationCanCheck = make(map[string]string)
	PerformanceCanCheck = make(map[string]string)
	BaselineCanCheck = make(map[string]string)
	SecurityCanCheck = make(map[string]string)
)
func YamlconfigInit(){
	tmpStatusValCheck := "true false"
	//检测配置开关参数是否设置合理，填写的true或false是否正确
	for k,v := range InspectionConfSwitch.ConfigSwitch{
		if !strings.Contains(tmpStatusValCheck, strings.ToLower(v)) {
			fmt.Println(fmt.Sprintf("当前配置参数 %s 选项输入有误，只能输入是true或false", k))
			os.Exit(1)
		}
	}
	for k,v := range InspectionConfSwitch.ConfigSwitch{
		//初始化配置参数，生成需要检测的列表
		if strings.EqualFold(k,"databaseConfigurationSwitch") && strings.EqualFold(v,"true"){
			for i := range InspectionConfInput.DatabaseConfiguration {
				cc := InspectionConfInput.DatabaseConfiguration[i]
				if !strings.Contains(tmpStatusValCheck, strings.ToLower(cc["subCheckSwitch"])) {
					fmt.Println("当前配置参数subCheckSwitchch 选项输入有误，只能输入是true或false")
					os.Exit(1)
				}
				if strings.EqualFold(cc["subCheckSwitch"], "true") {
					key := strings.ToLower(cc["checkconfiguration"])
					val := strings.ToLower(cc["checkThreshold"])
					ConfigurationCanCheck[key] = val
				}
			}
		}
		//初始化性能列表，生成需要检测的列表
		if strings.EqualFold(k,"databasePerformanceSwitch") && strings.EqualFold(v,"true"){
			for i := range InspectionConfInput.DatabasePerformance {
				cc := InspectionConfInput.DatabasePerformance[i]
				if !strings.Contains(tmpStatusValCheck, strings.ToLower(cc["subCheckSwitch"])) {
					fmt.Println("当前配置参数subCheckSwitchch 选项输入有误，只能输入是true或false")
					os.Exit(1)
				}
				if strings.EqualFold(cc["subCheckSwitch"], "true") {
					key := strings.ToLower(cc["checkconfiguration"])
					val := strings.ToLower(cc["checkThreshold"])
					PerformanceCanCheck[key] = val
				}
			}
		}
		//初始化基线检查列表，生成需要检测的列表
		if strings.EqualFold(k,"databaseBaselineSwitch") && strings.EqualFold(v,"true"){
			for i := range InspectionConfInput.DatabaseBaseline {
				cc := InspectionConfInput.DatabaseBaseline[i]
				if !strings.Contains(tmpStatusValCheck, strings.ToLower(cc["subCheckSwitch"])) {
					fmt.Println("当前配置参数subCheckSwitchch 选项输入有误，只能输入是true或false")
					os.Exit(1)
				}
				if strings.EqualFold(cc["subCheckSwitch"], "true") {
					key := strings.ToLower(cc["checkconfiguration"])
					val := strings.ToLower(cc["checkThreshold"])
					BaselineCanCheck[key] = val
				}
			}
		}
		//初始化安全检查列表，生成需要检测的列表
		if strings.EqualFold(k,"databaseSecuritySwitch") && strings.EqualFold(v,"true"){
			for i := range InspectionConfInput.DatabaseSecurity {
				cc := InspectionConfInput.DatabaseSecurity[i]
				if !strings.Contains(tmpStatusValCheck, strings.ToLower(cc["subCheckSwitch"])) {
					fmt.Println("当前配置参数subCheckSwitchch 选项输入有误，只能输入是true或false")
					os.Exit(1)
				}
				if strings.EqualFold(cc["subCheckSwitch"], "true") {
					key := strings.ToLower(cc["checkconfiguration"])
					val := strings.ToLower(cc["checkThreshold"])
					SecurityCanCheck[key] = val
				}
			}
		}
	}
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
