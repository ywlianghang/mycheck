package InspectionItem

import (
	pub "DepthInspection/api/PublicClass"
	"fmt"
	"strings"
)

//配置参数检查功能
func DatabaseConfigCheck(confParameterList map[string]string)  {
	pub.Loggs.Info("Begin to check that the database configuration parameters are properly configured")
	var configVariablesName ,configValue string
	for i := range confParameterList{
		configVariablesName = i
		configValue = confParameterList[i]
		pub.Loggs.Debug("Start checking database parameters ",configVariablesName)
		a,ok := pub.GlobalVariables[configVariablesName]
		if !ok {
			pub.Loggs.Error("The current data configuration parameter does not exist. Please check if it is incorrectly typed")
		}
		d  := make(map[string]string)
		d["configVariableName"] = configVariablesName
		d["configVariable"] = a   //当前值
		d["configValue"] = configValue //建议值
		d["checkStatus"] = "normal"    //正常
		d["checkType"] = "configParameter"
		d["errorCode"] = "CF1-01"
		if !strings.EqualFold(a,configValue) {
			d["checkStatus"] = "abnormal"    //异常
			d["threshold"] = configValue
			d["currentValue"] = fmt.Sprintf("%s=%s",configVariablesName,a)
			errorStrinfo := fmt.Sprintf(" CF1-01 The current database configuration is \"%s\" Not meeting reservation requirements! The current value of \"%s\" You are advised to set it to \"%s\"",configVariablesName,a,configValue)
			pub.Loggs.Error(errorStrinfo)
		}
		pub.InspectionResult.DatabaseConfigCheck.ConfigParameter = append(pub.InspectionResult.DatabaseConfigCheck.ConfigParameter,d)
	}

	pub.Loggs.Info("The check database configuration parameters are complete")
}



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

func newMap(source map[string]string) map[string]string {
	var n = make(map[string]string)
	for k,v := range source {
		n[k]=v
	}
	return n
}
//数据库的基线检查功能--检查表设计合规性
//检查表字符集是否为utf8
func (baselineCheck *DatabaseBaselineCheckStruct) BaselineCheckTablesDesign() {
	//表字符集检查 ~ 	表引擎检查
	pub.Loggs.Info("Begin a baseline check to check database table design compliance")
	var tableCharset string
	//字符集处理,生成字符集对应表
	var tmpCharsetCorrespondingTable = make(map[string]string)   //字符集对应表
	informationSchemaCollationsData := pub.InformationSchemaCollationsData
	for k := range informationSchemaCollationsData {
		ac := informationSchemaCollationsData[k]["COLLATION_NAME"].(string)
		ad := informationSchemaCollationsData[k]["CHARACTER_SET_NAME"].(string)
		tmpCharsetCorrespondingTable[ac] = ad
	}
	for ki,vi := range pub.BaselineCanCheck {
		if strings.EqualFold(ki, "tableCharset") || strings.EqualFold(ki, "tableEngine"){
			for i := range pub.InformationSchemaTablesData {
				ist := pub.InformationSchemaTablesData[i]
				var d = make(map[string]string)
				if ist["TABLE_COLLATION"] != nil {
					if _, ok := tmpCharsetCorrespondingTable[ist["TABLE_COLLATION"].(string)]; ok {
						tableCharset = tmpCharsetCorrespondingTable[ist["TABLE_COLLATION"].(string)]
					}
				}
				//表字符集检查
				d["database"] = ist["TABLE_SCHEMA"].(string)
				d["tableName"] = ist["TABLE_NAME"].(string)
				if strings.EqualFold(ki, "tableCharset") {
					if !strings.Contains(vi, tableCharset) {
						d["charset"] = tableCharset
						d["checkStatus"] = "abnormal" //异常
						d["checkType"] = "tableCharset"
						d["threshold"] = fmt.Sprintf("非%s", vi)
						d["errorCode"] = "BL1-TC01"
						d["currentValue"] = fmt.Sprintf("%s.%s", ist["TABLE_SCHEMA"], ist["TABLE_NAME"])
						pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableCharset = append(pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableCharset, d)
						pub.Loggs.Error(fmt.Sprintf(" BL1-TC01 The current table character set is not UTF8 or UTF8MB4 character. error info: Database is \"%s\" table is \"%s\" table charset is \"%s\" ", ist["TABLE_SCHEMA"], ist["TABLE_NAME"], tableCharset))
					} else {
						d["charset"] = tableCharset
						d["checkStatus"] = "normal" //异常
						d["checkType"] = "tableCharset"
						pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableCharset = append(pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableCharset, d)
					}
				}
				//检查引擎不是innodb的
				if strings.EqualFold(ki, "tableEngine") {
					m := newMap(d)
					m["checkType"] = "tableEngine"
					if pub.InformationSchemaTablesData[i]["ENGINE"] != nil && !strings.EqualFold(pub.InformationSchemaTablesData[i]["ENGINE"].(string), "innodb") {
						m["checkStatus"] = "abnormal"
						m["threshold"] = fmt.Sprintf("非%s", vi)
						m["errorCode"] = "BL1-TC02"
						m["currentValue"] = fmt.Sprintf("%s.%s", ist["TABLE_SCHEMA"], ist["TABLE_NAME"])
						pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableEngine = append(pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableEngine, m)
						pub.Loggs.Error(fmt.Sprintf(" BL1-TC02 The current table engine set is not innodb engine. error info: Database is \"%s\" table is \"%s\" table engine is \"%s\" ", ist["TABLE_SCHEMA"], ist["TABLE_NAME"], ist["ENGINE"]))
					}
					if ist["ENGINE"] != nil && strings.EqualFold(ist["ENGINE"].(string), "innodb") {
						m["checkStatus"] = "normal"
						m["currentValue"] = fmt.Sprintf("%s.%s", ist["TABLE_SCHEMA"], ist["TABLE_NAME"])
						pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableEngine = append(pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableEngine, m)
					}
				}
			}
		}
		//检查表是否使用外键
		if strings.EqualFold(ki, "tableForeign") {
			for i := range pub.InformationSchemaKeyColumnUsage {
				iskl := pub.InformationSchemaKeyColumnUsage[i]
				var d = make(map[string]string)
				d["database"] = iskl["databaseName"].(string)
				d["tableName"] = iskl["tableName"].(string)
				d["checkType"] = "tableForeign"
				if iskl["REFERENCED_TABLE_NAME"] != nil && iskl["REFERENCED_COLUMN_NAME"] != nil {
					d["checkStatus"] = "abnormal" //异常
					d["threshold"] = fmt.Sprintf("%s",vi)
					d["errorCode"] = "BL1-TC03"
					d["currentValue"] = fmt.Sprintf("%s.%s", iskl["databaseName"], iskl["tableName"])
					d["columnName"] = iskl["columnName"].(string)
					d["constraintName"] = iskl["CONSTRAINT_NAME"].(string)
					d["referencedTableName"] = iskl["REFERENCED_TABLE_NAME"].(string)
					d["referencedColumnName"] = iskl["REFERENCED_COLUMN_NAME"].(string)
					pub.Loggs.Error(fmt.Sprintf(" BL1-TC03 The current table uses a foreign key constraint. The information is as follows: database: \"%s\" "+
						"tableName: \"%s\" column: \"%s\" Foreign key constraint name: \"%s\" Foreign key constraints table: \"%s\""+
						"Foreign key constraints columns: \"%s\"", iskl["databaseName"], iskl["tableName"], iskl["columnName"], iskl["CONSTRAINT_NAME"], iskl["REFERENCED_TABLE_NAME"], iskl["REFERENCED_COLUMN_NAME"]))
					pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableForeign = append(pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableForeign, d)
				} else {
					d["checkStatus"] = "normal" //正常
					pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableForeign = append(pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableForeign, d)
				}
			}
		}
		//检查没有主键的表
		if strings.EqualFold(ki, "tableNoPrimaryKey") {
			var ke, vl string
			for v := range pub.InformationSchemaColumnsData {
				icd := pub.InformationSchemaColumnsData[v]
				var dd = make(map[string]string)
				dd["database"] = icd["TABLE_SCHEMA"].(string)
				dd["tableName"] = icd["TABLE_NAME"].(string)
				dd["checkType"] = "tableNoPrimaryKey"
				if icd["COLUMN_KEY"] == "PRI" {
					dd["checkStatus"] = "normal" //正常
					pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableNoPrimaryKey = append(pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableNoPrimaryKey, dd)
				} else {
					if ke != icd["TABLE_SCHEMA"].(string) || vl != icd["TABLE_NAME"].(string) {
						dd["checkStatus"] = "abnormal" //异常
						dd["threshold"] = fmt.Sprintf("%s",vi)
						dd["errorCode"] = "BL1-TC04"
						dd["currentValue"] = fmt.Sprintf("%s.%s", icd["TABLE_SCHEMA"], icd["TABLE_NAME"])
						pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableNoPrimaryKey = append(pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableNoPrimaryKey, dd)
						pub.Loggs.Error(fmt.Sprintf(" BL1-TC04 The current table has no primary key. error info: Database is \"%s\" table is \"%s\"", icd["TABLE_SCHEMA"], icd["TABLE_NAME"]))
					}
				}
				ke = icd["TABLE_SCHEMA"].(string)
				vl = icd["TABLE_NAME"].(string)
			}
		}
	}
}

//列设计合规性
func (baselineCheck *DatabaseBaselineCheckStruct) BaselineCheckColumnsDesign(){
	pub.Loggs.Info("Begin a baseline check to check database columns design compliance")
	for i := range pub.InformationSchemaColumnsData {
		icd := pub.InformationSchemaColumnsData[i]
		var d = make(map[string]string)
		d["database"] = icd["TABLE_SCHEMA"].(string)
		d["tableName"] = icd["TABLE_NAME"].(string)
		d["columnName"] = icd["COLUMN_NAME"].(string)
		//主键自增列是否为bigint
		if vi,ok := pub.BaselineCanCheck["tableautoincrement"];ok && pub.InformationSchemaColumnsData[i]["EXTRA"] == "auto_increment" {
				if !strings.Contains(icd["COLUMN_TYPE"].(string), vi) {
					d["checkStatus"] = "abnormal" //异常
					d["checkType"] = "tableAutoIncrement"
					d["threshold"] = "无自增主键"
					d["errorCode"] = "BL2-CC01"
					d["currentValue"] = fmt.Sprintf("%s.%s", icd["TABLE_SCHEMA"], icd["TABLE_NAME"])
					pub.InspectionResult.DatabaseBaselineCheck.ColumnDesign.TableAutoIncrement = append(pub.InspectionResult.DatabaseBaselineCheck.ColumnDesign.TableAutoIncrement, d)
					pub.Loggs.Error(fmt.Sprintf(" BL2-CC01 The primary key column is not of type Bigint. The information is as follows: database: \"%s\" tableName: \"%s\" columnsName: \"%s\" columnType: \"%s\".", icd["TABLE_SCHEMA"], icd["TABLE_NAME"], icd["COLUMN_NAME"], icd["COLUMN_TYPE"]))
				} else {
					d["checkStatus"] = "normal" //异常
					d["checkType"] = "tableAutoIncrement"
					pub.InspectionResult.DatabaseBaselineCheck.ColumnDesign.TableAutoIncrement = append(pub.InspectionResult.DatabaseBaselineCheck.ColumnDesign.TableAutoIncrement, d)
				}
			}
		if vi,ok := pub.BaselineCanCheck["tablebigcolumns"];ok && pub.InformationSchemaColumnsData[i]["EXTRA"] != "auto_increment" {
			//表中是否存在大字段blob、text、varchar(8099)、timestamp数据类型
			m := newMap(d)
			if icd["COLUMN_TYPE"] != nil {
				ce := icd["COLUMN_TYPE"].(string)
				if tmpIndexInt := strings.Index(ce,"(");tmpIndexInt != -1{
					ce = ce[:tmpIndexInt]
				}
				if strings.Contains(vi,ce){
					m["checkStatus"] = "abnormal" //异常
					m["checkType"] = "tableBigColumns"
					m["threshold"] = fmt.Sprintf("%s",vi)
					m["errorCode"] = "BL2-CC02"
					m["currentValue"] = fmt.Sprintf("%s.%s", icd["TABLE_SCHEMA"], icd["TABLE_NAME"])
					pub.InspectionResult.DatabaseBaselineCheck.ColumnDesign.TableBigColumns = append(pub.InspectionResult.DatabaseBaselineCheck.ColumnDesign.TableBigColumns, m)
					pub.Loggs.Error(fmt.Sprintf(" BL2-CC02 The column data types of the current table in the database exist BLOB, TEXT, TIMESTAMP. The information is as follows: database: \"%s\" tableName: \"%s\" columnsName: \"%s\" columnType: \"%s\".", icd["TABLE_SCHEMA"], icd["TABLE_NAME"], icd["COLUMN_NAME"], icd["COLUMN_TYPE"]))
				} else {
					m["checkStatus"] = "normal" //正常
					m["checkType"] = "tableBigColumns"
					pub.InspectionResult.DatabaseBaselineCheck.ColumnDesign.TableBigColumns = append(pub.InspectionResult.DatabaseBaselineCheck.ColumnDesign.TableBigColumns, m)
				}
			}
		}
			//var dd = make(map[string]string)
			//表列数是否大于255
		//}
	}
}

//索引设计合规性
func (baselineCheck *DatabaseBaselineCheckStruct) BaselineCheckIndexColumnDesign() {
	pub.Loggs.Info("Begin by checking that index usage is reasonable and index column creation is standard")
	var tmpMap = make(map[string]string)
	for i := range pub.InformationSchemaStatistics {
		iss := pub.InformationSchemaStatistics[i]
		v := fmt.Sprintf("%s_%s_%s", iss["TABLE_SCHEMA"], iss["TABLE_NAME"], iss["COLUMN_NAME"])
		tmpMap[v] = iss["INDEX_NAME"].(string)
	}
	for i := range pub.InformationSchemaColumnsData {
		var d = make(map[string]string)
		icd := pub.InformationSchemaColumnsData[i]
		v := fmt.Sprintf("%s_%s_%s", icd["TABLE_SCHEMA"], icd["TABLE_NAME"], icd["COLUMN_NAME"])
		d["database"] = icd["TABLE_SCHEMA"].(string)
		d["tableName"] = icd["TABLE_NAME"].(string)
		d["columnName"] = icd["COLUMN_NAME"].(string)
		d["columnType"] = icd["COLUMN_TYPE"].(string)
		d["columnIsNull"] = icd["IS_NULLABLE"].(string)
		if _, ok := tmpMap[v]; ok {
			//判断索引列是否允许为空
			if vi, okk := pub.BaselineCanCheck["indexcolumnisnull"]; okk {
				if strings.EqualFold(vi, icd["IS_NULLABLE"].(string)) {
					d["checkStatus"] = "abnormal" //异常
					d["checkType"] = "indexColumnIsNull"
					d["threshold"] = "索引列为空"
					d["errorCode"] = "BL3-IC01"
					d["currentValue"] = fmt.Sprintf("%s.%s", icd["TABLE_SCHEMA"], icd["TABLE_NAME"])
					pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnIsNull = append(pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnIsNull, d)
					pub.Loggs.Error(fmt.Sprintf(" BL3-IC01 An index column is empty.The information is as follows: database: \"%s\"  tablename: \"%s\" indexName: \"%s\" columnName: \"%s\" columnType: \"%s\"", icd["TABLE_SCHEMA"], icd["TABLE_NAME"], tmpMap[v], icd["COLUMN_NAME"], icd["COLUMN_TYPE"]))
				} else {
					d["checkStatus"] = "normal" //异常
					d["checkType"] = "indexColumnIsNull"
					pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnIsNull = append(pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnIsNull, d)
				}
			}
			//判断索引列是否建立在enum、set、blob、text类型上面
			if vi, okk := pub.BaselineCanCheck["indexcolumntype"]; okk {
				tmpColumnType := strings.ToLower(icd["COLUMN_TYPE"].(string))
				if tmpint := strings.Index(tmpColumnType, "("); tmpint > -1 {
					tmpColumnType = tmpColumnType[:tmpint]
				}
				if strings.Contains(strings.ToLower(vi), tmpColumnType) {
					m := newMap(d)
					m["checkStatus"] = "abnormal" //异常
					m["checkType"] = "indexColumnType"
					m["threshold"] = fmt.Sprintf("%s", vi)
					m["errorCode"] = "BL3-IC02"
					m["currentValue"] = fmt.Sprintf("%s.%s", icd["TABLE_SCHEMA"], icd["TABLE_NAME"])
					pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnType = append(pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnType, m)
					pub.Loggs.Error(fmt.Sprintf(" BL3-IC02 An index column is enum or set or blob or text type. The information is as follows: The information is as follows: database: \"%s\"  tablename: \"%s\" indexName: \"%s\" columnName: \"%s\" columnType: \"%s\"", icd["TABLE_SCHEMA"], icd["TABLE_NAME"], tmpMap[v], icd["COLUMN_NAME"], icd["COLUMN_TYPE"]))
				} else {
					m := newMap(d)
					m["checkStatus"] = "normal" //异常
					m["checkType"] = "indexColumnType"
					pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnType = append(pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnType, m)
				}
			}
		}
	}
	if _,okk := pub.BaselineCanCheck["tableincluderepeatindex"];okk {
		//利用map合并联合索引列
		var tmpIndexMargeMap = make(map[string]string)
		for k := range pub.InformationSchemaStatistics {
			b := pub.InformationSchemaStatistics
			key := fmt.Sprintf("%s@%s@@%s", b[k]["TABLE_SCHEMA"], b[k]["TABLE_NAME"], b[k]["INDEX_NAME"])
			if val, ok := tmpIndexMargeMap[key]; ok && k > 1 && b[k-1]["TABLE_SCHEMA"] == b[k]["TABLE_SCHEMA"] && b[k-1]["TABLE_NAME"] == b[k]["TABLE_NAME"] {
				tmpValue := fmt.Sprintf("%s,%s", val, b[k]["COLUMN_NAME"])
				tmpIndexMargeMap[key] = tmpValue
			} else {
				tmpValue := fmt.Sprintf("%s", b[k]["COLUMN_NAME"])
				tmpIndexMargeMap[key] = tmpValue
			}
		}
		//分离出每个库表下包含的索引
		var tmpdatabaseTableIncludeIndexMap = make(map[string]map[string]string)
		for k, v := range tmpIndexMargeMap {
			tmpMapp := make(map[string]string)
			a := strings.Split(k, "@@") //库表
			if val, ok := tmpdatabaseTableIncludeIndexMap[a[0]]; ok {
				for tmpk := range val {
					tmpMapp[tmpk] = val[tmpk] //旧的key value
				}
				tmpMapp[a[1]] = v //新的key value
				tmpdatabaseTableIncludeIndexMap[a[0]] = tmpMapp
			} else {
				tmpMapp[a[1]] = v
				tmpdatabaseTableIncludeIndexMap[a[0]] = tmpMapp
			}
		}
		//遍历每一个库表下的索引列，寻找冗余索引
		for k, v := range tmpdatabaseTableIncludeIndexMap {
			var d = make(map[string]string)
			var tmpRedundancyIndexStatus = false
			var tmpDatabase, tmpTablename, tmpIndexRedundancyName, tmpIndexRedundancyColumn, tmpIndexColumnName, tmpIndexIncludeColumn string
			a := strings.Split(k, "@")
			tmpDatabase = a[0]
			tmpTablename = a[1]
			for ki, ui := range v {
				for kii, uii := range v {
					if ui != uii && strings.HasPrefix(uii, ui) {
						tmpIndexRedundancyColumn = ui
						tmpIndexIncludeColumn = uii
						tmpIndexColumnName = kii
						tmpIndexRedundancyName = ki
						tmpRedundancyIndexStatus = true
					}
				}
			}
			d["database"] = tmpDatabase
			d["tableName"] = tmpTablename
			d["redundantIndexes"] = fmt.Sprintf("%s %s,%s %s", tmpIndexRedundancyName, tmpIndexRedundancyColumn, tmpIndexColumnName, tmpIndexIncludeColumn)
			d["checkStatus"] = "normal" //正常
			d["checkType"] = "tableIncludeRepeatIndex"
			if tmpRedundancyIndexStatus {
				d["checkStatus"] = "abnormal" //异常
				d["threshold"] = "存在重复索引"
				d["errorCode"] = "BL3-IC04"
				d["currentValue"] = fmt.Sprintf("%s.%s", tmpDatabase, tmpTablename)
				pub.Loggs.Error(fmt.Sprintf(" BL3-IC03 Redundant index columns appear. The information is as follows: database:\"%s\" tablename: \"%s\" Redundant indexes: (indexName:\"%s\" indexColumns \"%s\"), (indexName: \"%s\" indexColumns: \"%s\")", tmpDatabase, tmpTablename, tmpIndexRedundancyName, tmpIndexRedundancyColumn, tmpIndexColumnName, tmpIndexIncludeColumn))
			}
			pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnIsRepeatIndex = append(pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnIsRepeatIndex, d)
		}
	}
}
//存储过程、存储函数、触发器、视图、分区表检查限制
func(baselineCheck *DatabaseBaselineCheckStruct) BaselineCheckProcedureTriggerDesign(){
	pub.Loggs.Info("Begin a baseline check to checking whether the database uses stored procedures, stored functions, or triggers")
	if vi,okk := pub.BaselineCanCheck["tableprocedurefunctriggerviews"];okk {
		for i := range pub.InformationSchemaRoutines {
			cc := pub.InformationSchemaRoutines[i]
			var d = make(map[string]string)
			d["database"] = cc["ROUTINE_SCHEMA"].(string)
			if strings.Contains(strings.ToLower(vi),strings.ToLower(cc["ROUTINE_TYPE"].(string))) {
				d["checkStatus"] = "abnormal" //异常状态
				d["checkType"] = "tableProcedureFunc"
				d["threshold"] = fmt.Sprintf("%s",strings.ToLower(cc["ROUTINE_TYPE"].(string)))
				d["errorCode"] = "BL4-PT01"
				d["currentValue"] = fmt.Sprintf("%s.%s", cc["ROUTINE_SCHEMA"], cc["ROUTINE_NAME"])
				pub.InspectionResult.DatabaseBaselineCheck.ProcedureTriggerDesign.TableProcedure = append(pub.InspectionResult.DatabaseBaselineCheck.ProcedureTriggerDesign.TableProcedure, d)
				pub.Loggs.Error(fmt.Sprintf(" BL4-PT01 The current database uses a storage procedure or storage func. The information is as follows: database: \"%s\" routineName: \"%s\" user: \"%s\" create time: \"%s\"", cc["ROUTINE_SCHEMA"], cc["ROUTINE_NAME"], cc["DEFINER"], cc["CREATED"]))
			}
		}
		// 检查是否使用触发器
		for i := range pub.InformationSchemaTriggers {
			var d = make(map[string]string)
			dd := pub.InformationSchemaTriggers[i]
			d["database"] = dd["TRIGGER_SCHEMA"].(string)
			if strings.Contains(strings.ToLower(vi),"trigger") && dd["TRIGGER_NAME"] != nil {
				d["checkStatus"] = "abnormal" //异常状态
				d["checkType"] = "tableTrigger"
				d["threshold"] = fmt.Sprintf("%s","trigger")
				d["errorCode"] = "BL4-PT02"
				d["currentValue"] = fmt.Sprintf("%s.%s", dd["TRIGGER_SCHEMA"], dd["TRIGGER_NAME"])
				pub.InspectionResult.DatabaseBaselineCheck.ProcedureTriggerDesign.TableTrigger = append(pub.InspectionResult.DatabaseBaselineCheck.ProcedureTriggerDesign.TableTrigger, d)
				pub.Loggs.Error(fmt.Sprintf(" BL4-PT02 The current database uses a trigger. The information is as follows: database: \"%s\" triggerName: \"%s\"  user: \"%s\"  create time:\"%s\"", dd["TRIGGER_SCHEMA"], dd["TRIGGER_NAME"], dd["DEFINER"], dd["CREATED"]))
			}
		}
		//检查是否使用试图
		for i := range pub.InformationSchemaViews{
			var d = make(map[string]string)
			dd := pub.InformationSchemaViews[i]
			d["database"] = dd["TABLE_SCHEMA"].(string)
			if strings.Contains(strings.ToLower(vi),"views") && dd["VIEW_DEFINITION"] != nil {
				d["checkStatus"] = "abnormal" //异常状态
				d["checkType"] = "tableViews"
				d["threshold"] = fmt.Sprintf("%s","views")
				d["errorCode"] = "BL4-PT03"
				d["currentValue"] = fmt.Sprintf("%s.%s", dd["TABLE_SCHEMA"], dd["TABLE_NAME"])
				pub.InspectionResult.DatabaseBaselineCheck.ProcedureTriggerDesign.TableView = append(pub.InspectionResult.DatabaseBaselineCheck.ProcedureTriggerDesign.TableView, d)
				pub.Loggs.Error(fmt.Sprintf(" BL4-PT02 The current database uses a trigger. The information is as follows: database: \"%s\" triggerName: \"%s\"  DEFINER: \"%s\"  ", dd["TABLE_SCHEMA"], dd["TABLE_NAME"], dd["DEFINER"]))
			}
		}
	}


	pub.Loggs.Info("Check whether the database is completed using stored programs, stored functions, and stored triggers")
}