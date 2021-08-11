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
		if !strings.EqualFold(a,configValue) {
			d["checkStatus"] = "abnormal"    //异常
			d["checkType"] = "configParameter"
			errorStrinfo := fmt.Sprintf("The current database configuration is \"%s\" Not meeting reservation requirements! The current value of \"%s\" You are advised to set it to \"%s\"",configVariablesName,a,configValue)
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
	for i := range pub.InformationSchemaTablesData{
		var d = make(map[string]string)
		if pub.InformationSchemaTablesData[i]["TABLE_COLLATION"] != nil{
			if _,ok := tmpCharsetCorrespondingTable[pub.InformationSchemaTablesData[i]["TABLE_COLLATION"].(string)];ok {
				tableCharset = tmpCharsetCorrespondingTable[pub.InformationSchemaTablesData[i]["TABLE_COLLATION"].(string)]
			}
		}
		//表字符集检查
		d["database"] = pub.InformationSchemaTablesData[i]["TABLE_SCHEMA"].(string)
		d["tableName"] = pub.InformationSchemaTablesData[i]["TABLE_NAME"].(string)
		if !strings.Contains(tableCharset, "utf8"){
			d["charset"] = tableCharset
			d["checkStatus"] = "abnormal"    //异常
			d["checkType"] = "tableCharset"
			pub.InspectionResult.BaselineCheckTablesDesign.TableCharset = append(pub.InspectionResult.BaselineCheckTablesDesign.TableCharset,d)
			pub.Loggs.Error(fmt.Sprintf("The current table character set is not UTF8 or UTF8MB4 character. error info: Database is \"%s\" table is \"%s\" table charset is \"%s\" ",pub.InformationSchemaTablesData[i]["TABLE_SCHEMA"],pub.InformationSchemaTablesData[i]["TABLE_NAME"],tableCharset))
		}else{
			d["charset"] = tableCharset
			d["checkStatus"] = "normal"    //异常
			d["checkType"] = "tableCharset"
			pub.InspectionResult.BaselineCheckTablesDesign.TableCharset = append(pub.InspectionResult.BaselineCheckTablesDesign.TableCharset,d)
		}
		//检查引擎不是innodb的
		m := newMap(d)
		if  pub.InformationSchemaTablesData[i]["ENGINE"] != nil && !strings.EqualFold(pub.InformationSchemaTablesData[i]["ENGINE"].(string),"innodb"){
			m["checkStatus"] = "abnormal"
			m["checkType"] = "tableEngine"
			pub.InspectionResult.BaselineCheckTablesDesign.TableEngine = append(pub.InspectionResult.BaselineCheckTablesDesign.TableEngine,m)
			pub.Loggs.Error(fmt.Sprintf("The current table engine set is not innodb engine. error info: Database is \"%s\" table is \"%s\" table engine is \"%s\" ",pub.InformationSchemaTablesData[i]["TABLE_SCHEMA"],pub.InformationSchemaTablesData[i]["TABLE_NAME"],pub.InformationSchemaTablesData[i]["ENGINE"]))
		}
		if pub.InformationSchemaTablesData[i]["ENGINE"] != nil && strings.EqualFold(pub.InformationSchemaTablesData[i]["ENGINE"].(string),"innodb"){
			m["checkType"] = "tableEngine"
			m["checkStatus"] = "normal"
			pub.InspectionResult.BaselineCheckTablesDesign.TableEngine = append(pub.InspectionResult.BaselineCheckTablesDesign.TableEngine,m)
		}
	}

	//检查表是否使用外键
	for i := range pub.InformationSchemaKeyColumnUsage {
		var d = make(map[string]string)
		d["database"] = pub.InformationSchemaKeyColumnUsage[i]["databaseName"].(string)
		d["tableName"] = pub.InformationSchemaKeyColumnUsage[i]["tableName"].(string)
		d["checkType"] = "tableForeign"
		d["checkStatus"] = "normal"    //正常
		if pub.InformationSchemaKeyColumnUsage[i]["REFERENCED_TABLE_NAME"] != nil && pub.InformationSchemaKeyColumnUsage[i]["REFERENCED_COLUMN_NAME"] != nil {
			d["checkStatus"] = "abnormal"    //异常
			d["columnName"] = pub.InformationSchemaKeyColumnUsage[i]["columnName"].(string)
			d["constraintName"] = pub.InformationSchemaKeyColumnUsage[i]["CONSTRAINT_NAME"].(string)
			d["referencedTableName"] = pub.InformationSchemaKeyColumnUsage[i]["REFERENCED_TABLE_NAME"].(string)
			d["referencedColumnName"] = pub.InformationSchemaKeyColumnUsage[i]["REFERENCED_COLUMN_NAME"].(string)
			pub.Loggs.Error(fmt.Sprintf("The current table uses a foreign key constraint. The information is as follows: database: \"%s\" " +
				"tableName: \"%s\" column: \"%s\" Foreign key constraint name: \"%s\" Foreign key constraints table: \"%s\"" +
				"Foreign key constraints columns: \"%s\"",pub.InformationSchemaKeyColumnUsage[i]["databaseName"],pub.InformationSchemaKeyColumnUsage[i]["tableName"],pub.InformationSchemaKeyColumnUsage[i]["columnName"],pub.InformationSchemaKeyColumnUsage[i]["CONSTRAINT_NAME"],pub.InformationSchemaKeyColumnUsage[i]["REFERENCED_TABLE_NAME"],pub.InformationSchemaKeyColumnUsage[i]["REFERENCED_COLUMN_NAME"]))
		}
		pub.InspectionResult.BaselineCheckTablesDesign.TableForeign = append(pub.InspectionResult.BaselineCheckTablesDesign.TableForeign,d)
	}
	//检查没有主键的表
	var ke,vl string
	for v :=range pub.InformationSchemaColumnsData{
		var dd = make(map[string]string)
		dd["database"] = pub.InformationSchemaColumnsData[v]["TABLE_SCHEMA"].(string)
		dd["tableName"] = pub.InformationSchemaColumnsData[v]["TABLE_NAME"].(string)
		dd["checkType"] = "tableNoPrimaryKey"
		if pub.InformationSchemaColumnsData[v]["COLUMN_KEY"] == "PRI" {
			dd["checkStatus"] = "normal"    //正常
			pub.InspectionResult.BaselineCheckTablesDesign.TableNoPrimaryKey = append(pub.InspectionResult.BaselineCheckTablesDesign.TableNoPrimaryKey,dd)
		}else{
			if ke != pub.InformationSchemaColumnsData[v]["TABLE_SCHEMA"].(string) || vl != pub.InformationSchemaColumnsData[v]["TABLE_NAME"].(string){
				dd["checkStatus"] = "abnormal"    //异常
				pub.InspectionResult.BaselineCheckTablesDesign.TableNoPrimaryKey = append(pub.InspectionResult.BaselineCheckTablesDesign.TableNoPrimaryKey,dd)
				pub.Loggs.Error(fmt.Sprintf("The current table has no primary key. error info: Database is \"%s\" table is \"%s\"",pub.InformationSchemaColumnsData[v]["TABLE_SCHEMA"],pub.InformationSchemaColumnsData[v]["TABLE_NAME"]))
			}
		}
		ke = pub.InformationSchemaColumnsData[v]["TABLE_SCHEMA"].(string)
		vl = pub.InformationSchemaColumnsData[v]["TABLE_NAME"].(string)
	}
}

//列设计合规性
func (baselineCheck *DatabaseBaselineCheckStruct) BaselineCheckColumnsDesign(){
	pub.Loggs.Info("Begin a baseline check to check database columns design compliance")
	for i := range pub.InformationSchemaColumnsData {
		var d = make(map[string]string)
		d["database"] = pub.InformationSchemaColumnsData[i]["TABLE_SCHEMA"].(string)
		d["tableName"] = pub.InformationSchemaColumnsData[i]["TABLE_NAME"].(string)
		d["columnName"] = pub.InformationSchemaColumnsData[i]["COLUMN_NAME"].(string)
		d["columnType"] = pub.InformationSchemaColumnsData[i]["COLUMN_TYPE"].(string)
		//主键自增列是否为bigint
		if pub.InformationSchemaColumnsData[i]["EXTRA"] == "auto_increment"{
			if  !strings.Contains(pub.InformationSchemaColumnsData[i]["COLUMN_TYPE"].(string),"bigint"){
				d["checkStatus"] = "abnormal"    //异常
				d["checkType"] = "tableAutoIncrement"
				pub.InspectionResult.BaselineCheckColumnsDesign.TableAutoIncrement = append(pub.InspectionResult.BaselineCheckColumnsDesign.TableAutoIncrement,d)
				pub.Loggs.Error(fmt.Sprintf("The primary key column is not of type Bigint. The information is as follows: database: \"%s\" tableName: \"%s\" columnsName: \"%s\" columnType: \"%s\".", pub.InformationSchemaColumnsData[i]["TABLE_SCHEMA"],pub.InformationSchemaColumnsData[i]["TABLE_NAME"],pub.InformationSchemaColumnsData[i]["COLUMN_NAME"],pub.InformationSchemaColumnsData[i]["COLUMN_TYPE"]))
			}else {
				d["checkStatus"] = "normal"    //异常
				d["checkType"] = "tableAutoIncrement"
				pub.InspectionResult.BaselineCheckColumnsDesign.TableAutoIncrement = append(pub.InspectionResult.BaselineCheckColumnsDesign.TableAutoIncrement,d)
			}
		}else {
			//表中是否存在大字段blob、text、varchar(8099)、timestamp数据类型
			m := newMap(d)
			if pub.InformationSchemaColumnsData[i]["COLUMN_TYPE"] != nil {
				ce := pub.InformationSchemaColumnsData[i]["COLUMN_TYPE"].(string)
				if ce == "blob" || strings.Contains(ce, "text") || ce == "timestamp" {
					m["checkStatus"] = "abnormal" //异常
					m["checkType"] = "tableBigColumns"
					pub.InspectionResult.BaselineCheckColumnsDesign.TableBigColumns = append(pub.InspectionResult.BaselineCheckColumnsDesign.TableBigColumns, m)
					pub.Loggs.Error(fmt.Sprintf("The column data types of the current table in the database exist BLOB, TEXT, TIMESTAMP. The information is as follows: database: \"%s\" tableName: \"%s\" columnsName: \"%s\" columnType: \"%s\".", pub.InformationSchemaColumnsData[i]["TABLE_SCHEMA"], pub.InformationSchemaColumnsData[i]["TABLE_NAME"], pub.InformationSchemaColumnsData[i]["COLUMN_NAME"], pub.InformationSchemaColumnsData[i]["COLUMN_TYPE"]))
				} else {
					m["checkStatus"] = "normal" //正常
					m["checkType"] = "tableBigColumns"
					pub.InspectionResult.BaselineCheckColumnsDesign.TableBigColumns = append(pub.InspectionResult.BaselineCheckColumnsDesign.TableBigColumns, m)
				}
			}
		}
		//var dd = make(map[string]string)
		//表列数是否大于255
	}
}

//索引设计合规性
func (baselineCheck *DatabaseBaselineCheckStruct) BaselineCheckIndexColumnDesign(){
	pub.Loggs.Info("Begin by checking that index usage is reasonable and index column creation is standard")
	var tmpMap = make(map[string]string)
	for i := range pub.InformationSchemaStatistics{
		a := pub.InformationSchemaStatistics[i]
		v := fmt.Sprintf("%s_%s_%s",a["TABLE_SCHEMA"],a["TABLE_NAME"],a["COLUMN_NAME"])
		tmpMap[v] = a["INDEX_NAME"].(string)
	}
	for i := range pub.InformationSchemaColumnsData{
		var d = make(map[string]string)
		a := pub.InformationSchemaColumnsData[i]
		v := fmt.Sprintf("%s_%s_%s",a["TABLE_SCHEMA"],a["TABLE_NAME"],a["COLUMN_NAME"])
		d["database"] = a["TABLE_SCHEMA"].(string)
		d["tableName"] = a["TABLE_NAME"].(string)
		d["columnName"] = a["COLUMN_NAME"].(string)
		d["columnType"] = a["COLUMN_TYPE"].(string)
		d["columnIsNull"] = a["IS_NULLABLE"].(string)
		if _,ok := tmpMap[v];ok {
				//判断索引列是否允许为空
			if a["IS_NULLABLE"] == "YES" {
				d["checkStatus"] = "abnormal" //异常
				d["checkType"] = "indexColumnIsNull"
				pub.InspectionResult.BaselineCheckIndexColumnDesign.IndexColumnIsNull = append(pub.InspectionResult.BaselineCheckIndexColumnDesign.IndexColumnIsNull, d)
				pub.Loggs.Error(fmt.Sprintf("An index column is empty.The information is as follows: database: \"%s\"  tablename: \"%s\" indexName: \"%s\" columnName: \"%s\" columnType: \"%s\"", a["TABLE_SCHEMA"], a["TABLE_NAME"], tmpMap[v], a["COLUMN_NAME"], a["COLUMN_TYPE"]))
			} else {
				d["checkStatus"] = "normal" //异常
				d["checkType"] = "indexColumnIsNull"
				pub.InspectionResult.BaselineCheckIndexColumnDesign.IndexColumnIsNull = append(pub.InspectionResult.BaselineCheckIndexColumnDesign.IndexColumnIsNull, d)
			}

			m := newMap(d)
			columnTypeStr := fmt.Sprintf("%s",a["COLUMN_TYPE"])
			//判断索引列是否建立在enum或set类型上面
			if strings.Contains(columnTypeStr, "enum") || strings.Contains(columnTypeStr, "set") {
				m["checkStatus"] = "abnormal" //异常
				m["checkType"] = "indexColumnIsEnumSet"
				pub.InspectionResult.BaselineCheckIndexColumnDesign.IndexColumnIsEnumSet = append(pub.InspectionResult.BaselineCheckIndexColumnDesign.IndexColumnIsEnumSet, m)
				pub.Loggs.Error(fmt.Sprintf("An index column is enum or set type. The information is as follows: The information is as follows: database: \"%s\"  tablename: \"%s\" indexName: \"%s\" columnName: \"%s\" columnType: \"%s\"", a["TABLE_SCHEMA"], a["TABLE_NAME"], tmpMap[v], a["COLUMN_NAME"], a["COLUMN_TYPE"]))
			} else {
				m["checkStatus"] = "normal" //异常
				m["checkType"] = "indexColumnIsEnumSet"
				pub.InspectionResult.BaselineCheckIndexColumnDesign.IndexColumnIsEnumSet = append(pub.InspectionResult.BaselineCheckIndexColumnDesign.IndexColumnIsEnumSet, m)
			}
			n := newMap(d)
			//判断索引列是否建立在enum或set类型上面
			if strings.Contains(columnTypeStr, "blob") || strings.Contains(columnTypeStr, "text") {
				n["checkStatus"] = "abnormal" //异常
				n["checkType"] = "indexColumnIsBlobText"
				pub.InspectionResult.BaselineCheckIndexColumnDesign.IndexColumnIsBlobText = append(pub.InspectionResult.BaselineCheckIndexColumnDesign.IndexColumnIsBlobText, n)
				pub.Loggs.Error(fmt.Sprintf("An index column is blob or text type. The information is as follows: The information is as follows: database: \"%s\"  tablename: \"%s\" indexName: \"%s\" columnName: \"%s\" columnType: \"%s\"", a["databaseName"], a["tableName"], a["indexName"], a["columnName"], a["columnType"]))
			} else {
				n["checkStatus"] = "normal" //正常
				n["checkType"] = "indexColumnIsBlobText"
				pub.InspectionResult.BaselineCheckIndexColumnDesign.IndexColumnIsBlobText = append(pub.InspectionResult.BaselineCheckIndexColumnDesign.IndexColumnIsBlobText, n)
			}
		}
	}

	//利用map合并联合索引列
	var tmpIndexMargeMap = make(map[string]string)
	for k := range pub.InformationSchemaStatistics{
		b := pub.InformationSchemaStatistics
		key := fmt.Sprintf("%s@%s@@%s",b[k]["TABLE_SCHEMA"],b[k]["TABLE_NAME"],b[k]["INDEX_NAME"])
		if val,ok := tmpIndexMargeMap[key];ok && k>1 && b[k-1]["TABLE_SCHEMA"]==b[k]["TABLE_SCHEMA"] && b[k-1]["TABLE_NAME"]==b[k]["TABLE_NAME"]{
			tmpValue := fmt.Sprintf("%s,%s",val,b[k]["COLUMN_NAME"])
			tmpIndexMargeMap[key] = tmpValue
		}else{
			tmpValue := fmt.Sprintf("%s",b[k]["COLUMN_NAME"])
			tmpIndexMargeMap[key] = tmpValue
		}
	}
	//分离出每个库表下包含的索引
	var tmpdatabaseTableIncludeIndexMap = make(map[string]map[string]string)
	for k,v := range tmpIndexMargeMap{
		tmpMapp := make(map[string]string)
		a := strings.Split(k,"@@")  //库表
		if val,ok := tmpdatabaseTableIncludeIndexMap[a[0]];ok{
			for tmpk := range val{
				tmpMapp[tmpk] = val[tmpk]  //旧的key value
			}
			tmpMapp[a[1]] = v  //新的key value
			tmpdatabaseTableIncludeIndexMap[a[0]] = tmpMapp
		}else{
			tmpMapp[a[1]] = v
			tmpdatabaseTableIncludeIndexMap[a[0]] = tmpMapp
		}
	}
	//遍历每一个库表下的索引列，寻找冗余索引
	for k,v := range tmpdatabaseTableIncludeIndexMap{
		var d = make(map[string]string)
		var tmpRedundancyIndexStatus = false
		var tmpDatabase,tmpTablename,tmpIndexRedundancyName,tmpIndexRedundancyColumn,tmpIndexColumnName,tmpIndexIncludeColumn string
		a := strings.Split(k,"@")
		tmpDatabase = a[0]
		tmpTablename = a[1]
		for ki,ui := range v{
			for kii,uii := range v{
				if ui != uii && strings.HasPrefix(uii,ui){
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
		d["redundantIndexes"] = fmt.Sprintf("%s %s,%s %s" ,tmpIndexRedundancyName,tmpIndexRedundancyColumn,tmpIndexColumnName,tmpIndexIncludeColumn)
		d["checkStatus"] = "normal" //正常
		d["checkType"] = "tableIncludeRepeatIndex"
		if tmpRedundancyIndexStatus {
			d["checkStatus"] = "abnormal" //异常
			pub.Loggs.Error(fmt.Sprintf("Redundant index columns appear. The information is as follows: database:\"%s\" tablename: \"%s\" Redundant indexes: (indexName:\"%s\" indexColumns \"%s\"), (indexName: \"%s\" indexColumns: \"%s\")", tmpDatabase,tmpTablename,tmpIndexRedundancyName,tmpIndexRedundancyColumn,tmpIndexColumnName,tmpIndexIncludeColumn))
		}
		pub.InspectionResult.BaselineCheckIndexColumnDesign.IndexColumnIsRepeatIndex = append(pub.InspectionResult.BaselineCheckIndexColumnDesign.IndexColumnIsRepeatIndex,d)
	}
}

//存储过程、存储函数、触发器检查限制
func (baselineCheck *DatabaseBaselineCheckStruct) BaselineCheckProcedureTriggerDesign(){
	pub.Loggs.Info("Begin a baseline check to checking whether the database uses stored procedures, stored functions, or triggers")
	//var c []map[string]string
	for i := range pub.InformationSchemaRoutines{
		cc := pub.InformationSchemaRoutines[i]
		var d = make(map[string]string)
		d["database"] = cc["ROUTINE_SCHEMA"].(string)
		d["definer"] = cc["DEFINER"].(string)
		d["created"] = cc["CREATED"].(string)
		d["routineName"] = cc["ROUTINE_NAME"].(string)
		if  cc["ROUTINE_TYPE"] == "PROCEDURE" {
			d["checkStatus"] = "abnormal"    //异常状态
			d["checkType"] = "tableProcedure"
			pub.InspectionResult.BaselineCheckProcedureTriggerDesign.TableProcedure = append(pub.InspectionResult.BaselineCheckProcedureTriggerDesign.TableProcedure,d)
			pub.Loggs.Error(fmt.Sprintf("The current database uses a storage procedure. The information is as follows: database: \"%s\" routineName: \"%s\" user: \"%s\" create time: \"%s\"" ,cc["ROUTINE_SCHEMA"],cc["ROUTINE_NAME"],cc["DEFINER"],cc["CREATED"]))
		}
		m := newMap(d)
		if cc["ROUTINE_TYPE"] == "FUNCTION" {
			m["checkStatus"] = "abnormal"    //异常状态
			m["checkType"] = "tableFunc"
			pub.InspectionResult.BaselineCheckProcedureTriggerDesign.TableFunc = append(pub.InspectionResult.BaselineCheckProcedureTriggerDesign.TableFunc,m)
			pub.Loggs.Error(fmt.Sprintf("The current database uses a storage function . The information is as follows: database: \"%s\" routineName: \"%s\" user: \"%s\" create time: \"%s\"" ,cc["ROUTINE_SCHEMA"],cc["ROUTINE_NAME"],cc["DEFINER"],cc["CREATED"]))
		}
	}
	// 检查是否使用触发器
	for i := range pub.InformationSchemaTriggers{
		var d = make(map[string]string)
		dd := pub.InformationSchemaTriggers[i]
		d["database"] = dd["TRIGGER_SCHEMA"].(string)
		if dd["TRIGGER_NAME"] != nil{
			d["triggerName"] = dd["TRIGGER_NAME"].(string)
			d["definer"] = dd["DEFINER"].(string)
			d["created"] = dd["CREATED"].(string)
			d["checkStatus"] = "abnormal"    //异常状态
			d["checkType"] = "tableTrigger"
			pub.InspectionResult.BaselineCheckProcedureTriggerDesign.TableTrigger = append(pub.InspectionResult.BaselineCheckProcedureTriggerDesign.TableTrigger,d)
			pub.Loggs.Error(fmt.Sprintf("The current database uses a trigger. The information is as follows: database: \"%s\" triggerName: \"%s\"  user: \"%s\"  create time:\"%s\"" ,dd["TRIGGER_SCHEMA"],dd["TRIGGER_NAME"],dd["DEFINER"],dd["CREATED"]))
		}
	}
	pub.Loggs.Info("Check whether the database is completed using stored programs, stored functions, and stored triggers")
}