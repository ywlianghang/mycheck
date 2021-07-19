package InspectionItem

import (
	"DepthInspection/api/PublicDB"
	"fmt"
	"strings"
)

func DBConfigCheck(aa *PublicDB.ConfigInfo,confParameterList map[string]string) (bool,string) {
	var acd bool
	aa.Loggs.Info("Begin to check that the database configuration parameters are properly configured")
	DBdate := aa.DatabaseExecInterf.DBQueryDateMap(aa,"show global variables")
	var configVariablesName ,configValue string
	for i := range confParameterList{
		configVariablesName = i
		configValue = confParameterList[i]
		aa.Loggs.Debug("Start checking database parameters ",configVariablesName)
		a,ok := DBdate[configVariablesName]
		if !ok {
			aa.Loggs.Error("The current data configuration parameter does not exist. Please check if it is incorrectly typed")
		}
		if !strings.EqualFold(a,configValue) {
			errorStrinfo := fmt.Sprintf("检测当前数据库配置参数为 %s 不符合预定要求! 当前值为 %s 建议设置成 %s",configVariablesName,a,configValue)
			aa.Loggs.Error(errorStrinfo)
		}
	}
	aa.Loggs.Info("The check database configuration parameters are complete")
	return acd,configValue
}

//type databaseBaseLineCheckInterface interface {
//	TableDesignCompliance(aa *PublicDB.ConfigInfo)
//}
type DatabaseBaselineCheckStruct struct {
	strSql string
	ignoreTableSchema string
}
type TableDesignComplianceStruct struct {
	DatabaseName interface{} `json: "databaseName"`
	TableName interface{} `json: "tableName"`
	Engine interface{}  `json: "engine"`
	Charset interface{} `json: "charset"`
}

func (baselineCheck *DatabaseBaselineCheckStruct) BaselineCheckTablesDesign(aa *PublicDB.ConfigInfo) {
	//表字符集检查 ~ 	表引擎检查
	aa.Loggs.Info("Begin a baseline check to check database table design compliance")
	ignoreTableSchema := "'mysql','information_schema','performance_schema','sys'"
	strSql := fmt.Sprintf("select t.table_schema databaseName,t.table_name tableName,lower(engine) engine,lower(c.CHARACTER_SET_NAME) charset from information_schema.tables as t, information_schema.COLLATIONS as c where t.TABLE_COLLATION=c.COLLATION_NAME and t.table_schema not in (%s)",ignoreTableSchema)
	cc := aa.DatabaseExecInterf.DBQueryDateJson(aa,strSql)
	for i := range cc{
		//表字符集检查
		charsetSt := fmt.Sprintf("%v",cc[i]["charset"])
		if !strings.Contains(charsetSt,"utf8"){
			aa.Loggs.Error(fmt.Sprintf("The current table character set is not UTF8 or UTF8MB4 character. error info: Database is %s table is %s table charset is %s ",cc[i]["databaseName"],cc[i]["tableName"],cc[i]["charset"]))
		}
		//表引擎检查
		if cc[i]["engine"] != "innodb"{
			aa.Loggs.Error(fmt.Sprintf("The current table engine set is not innodb engine. error info: Database is %s table is %s table engine is %s ",cc[i]["databaseName"],cc[i]["tableName"],cc[i]["engine"]))
		}
	}
	//检查表是否使用外键
	strSql = fmt.Sprintf("select CONSTRAINT_SCHEMA databaseName,TABLE_NAME tableName,COLUMN_NAME columnName,CONSTRAINT_NAME, REFERENCED_TABLE_NAME,REFERENCED_COLUMN_NAME from INFORMATION_SCHEMA.KEY_COLUMN_USAGE where CONSTRAINT_SCHEMA not in (%s)",ignoreTableSchema)
	dd := aa.DatabaseExecInterf.DBQueryDateJson(aa,strSql)
	for i := range dd {
		if dd[i]["REFERENCED_TABLE_NAME"] != nil && dd[i]["REFERENCED_COLUMN_NAME"] != nil {
			aa.Loggs.Error(fmt.Sprintf("The current table uses a foreign key constraint. The information is as follows: database: %s " +
				"tableName: %s column: %s Foreign key constraint name: %s Foreign key constraints table: %s" +
				"Foreign key constraints columns: %s",dd[i]["databaseName"],dd[i]["tableName"],dd[i]["columnName"],dd[i]["CONSTRAINT_NAME"],dd[i]["REFERENCED_TABLE_NAME"],dd[i]["REFERENCED_COLUMN_NAME"]))
		}
	}
	//检查没有主键的表
	strSql = fmt.Sprintf("select table_schema databaseName, table_name tableName from information_schema.tables where table_name not in (select distinct table_name from information_schema.columns where column_key = 'PRI' ) AND table_schema not in (%s)",ignoreTableSchema)
	ee := aa.DatabaseExecInterf.DBQueryDateJson(aa,strSql)
	for i := range ee{
		if ee[i] != nil{
			aa.Loggs.Error(fmt.Sprintf("The current table has no primary key index. The information is as follows: database: %s tableName: %s" ,ee[i]["databaseName"],ee[i]["tableName"]))
		}
	}
}

//列设计合规性
func (baselineCheck *DatabaseBaselineCheckStruct) BaselineCheckColumnsDesign(aa *PublicDB.ConfigInfo){
	ignoreTableSchema := "'mysql','information_schema','performance_schema','sys'"
	aa.Loggs.Info("Begin a baseline check to check database columns design compliance")

	strSql := fmt.Sprintf("select table_Schema databaseName,table_name tableName,column_name columnName,column_type columnType,COLUMN_KEY columnKey,EXTRA extra from information_schema.columns where table_schema not in(%s)",ignoreTableSchema)
	cc := aa.DatabaseExecInterf.DBQueryDateJson(aa,strSql)
	for i := range cc {
		//主键自增列是否为bigint
		if cc[i]["extra"] == "auto_increment" && cc[i]["columnType"] != "bigint"{
			aa.Loggs.Error(fmt.Sprintf("The primary key column is not of type Bigint. The information is as follows: database: %s tableName: %s columnsName: %s columnType: %s.", cc[i]["databaseName"],cc[i]["tableName"],cc[i]["columnName"],cc[i]["columnType"]))
		}
		//表中是否存在大字段blob、text、varchar(8099)、timestamp数据类型
		ce := fmt.Sprintf("%v",cc[i]["columnType"])
		if cc[i]["columnType"] == "blob" || strings.Contains(ce,"text") || cc[i]["columnType"] == "timestamp"{
			aa.Loggs.Error(fmt.Sprintf("The column data types of the current table in the database exist BLOB, TEXT, TIMESTAMP. The information is as follows: database: %s tableName: %s columnsName: %s columnType: %s.",cc[i]["databaseName"],cc[i]["tableName"],cc[i]["columnName"],cc[i]["columnType"]))
		}
		//表列数是否大于255

	}




	//表中是否存在ENUM列



}
//索引设计合规性
func (baselineCheck *DatabaseBaselineCheckStruct) IndexCompliance(){}
//存储过程及存储函数检查限制
func (baselineCheck *DatabaseBaselineCheckStruct) StoredProcedureTriggerRestrictions(){}
