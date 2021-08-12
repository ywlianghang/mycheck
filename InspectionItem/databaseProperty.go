package InspectionItem

import (
	"DepthInspection/api/PublicClass"
	"fmt"
	"strconv"
	"strings"
)

const (
	MtinyInt = 127
	MunsignedTinyInt = 255
	MsmallInt = 32767
	MunsignedSmallInt = 65535
	MmediumInt = 8388607
	MunsignedMediumInt = 16777215
    Mint = 2147483647
    MunsignedInt = 4294967295
    MbigInt = 9223372036854775807
    MunsignedBigint uint64 = 18446744073709551615
)

func (baselineCheck *DatabaseBaselineCheckStruct) DatabasePerformanceStatusCheck() {
	//var c []map[string]string
	//统计使用磁盘的binlog写入占使用内存buffer的binlog写入的百分比，大于100%则需要增加binlog_cache_size
	binlogDiskUsageRate, err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["Binlog_cache_disk_use"], PublicClass.GlobalStatus["Binlog_cache_use"])
	if binlogDiskUsageRate > 100 && err == nil {
		var d = make(map[string]string)
		d["threshold"] = ">100%"
		d["errorCode"] = "PF1-01"
		d["checkType"] = "binlogDiskUsageRate"
		d["checkStatus"] = "abnormal"    //异常
		d["currentValue"] = fmt.Sprintf("%s=%s","binlogDiskUsageRate",strconv.Itoa(binlogDiskUsageRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.BinlogDiskUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.BinlogDiskUsageRate,d)
		PublicClass.Loggs.Error(" PF1-01 The current database binlog is using too many disk writes. It is recommended to modify the \"binlog_cache_size\" parameter")
	}else{
		var d = make(map[string]string)
		d["threshold"] = ">100%"
		d["errorCode"] = "PF1-01"
		d["checkType"] = "binlogDiskUsageRate"
		d["checkStatus"] = "normal"    //正常
		d["binlogDiskUsageRate"] = strconv.Itoa(binlogDiskUsageRate)
		d["currentValue"] = fmt.Sprintf("%s=%s","binlogDiskUsageRate",strconv.Itoa(binlogDiskUsageRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.BinlogDiskUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.BinlogDiskUsageRate,d)
	}
	
	//统计历史连接数最大使用率，使用创建过
	historyConnectionMaxUsageRate, err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["Threads_created"],PublicClass.GlobalVariables["max_connections"])
	if historyConnectionMaxUsageRate > 80 && err == nil {
		var d = make(map[string]string)
		d["threshold"] = ">80%"
		d["errorCode"] = "PF1-02"
		d["checkType"] = "historyConnectionMaxUsageRate"
		d["checkStatus"] = "abnormal"    //异常
		d["currentValue"] = fmt.Sprintf("%s=%s","historyConnectionMaxUsageRate",strconv.Itoa(historyConnectionMaxUsageRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.HistoryConnectionMaxUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.HistoryConnectionMaxUsageRate,d)
		PublicClass.Loggs.Error(" PF1-02 If the maximum usage of historical database connections exceeds \"80%\", change the \"max_connections\" value and check services")
	}else{
		var d = make(map[string]string)
		d["threshold"] = ">80%"
		d["errorCode"] = "PF1-02"
		d["checkType"] = "historyConnectionMaxUsageRate"
		d["checkStatus"] = "normal"    //正常
		d["currentValue"] = fmt.Sprintf("%s=%s","historyConnectionMaxUsageRate",strconv.Itoa(historyConnectionMaxUsageRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.HistoryConnectionMaxUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.HistoryConnectionMaxUsageRate,d)
	}
	
	//统计数据库使用中使用磁盘临时表占使用内存临时表的占用比例Created_tmp_disk_tables/Created_tmp_tables *100% <=25%
	tmpDiskTableUsageRate, err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["Created_tmp_disk_tables"], PublicClass.GlobalStatus["Created_tmp_tables"])
	if tmpDiskTableUsageRate > 25 && err == nil {
		var d = make(map[string]string)
		d["threshold"] = ">25%"
		d["errorCode"] = "PF1-03"
		d["checkType"] = "tmpDiskTableUsageRate"
		d["checkStatus"] = "abnormal"    //异常
		d["binlogDiskUsageRate"] = strconv.Itoa(tmpDiskTableUsageRate)
		d["currentValue"] = fmt.Sprintf("%s=%s","tmpDiskTableUsageRate",strconv.Itoa(tmpDiskTableUsageRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.TmpDiskTableUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.TmpDiskTableUsageRate,d)
		PublicClass.Loggs.Error(" PF1-03 Too many disk temporary tables are being used. Check the slow SQL log or parameters \"tmp_table_size\" and \"max_heap_table_size\"")
	}else{
		var d = make(map[string]string)
		d["threshold"] = ">25%"
		d["errorCode"] = "PF1-03"
		d["checkType"] = "tmpDiskTableUsageRate"
		d["checkStatus"] = "normal"    //正常
		d["currentValue"] = fmt.Sprintf("%s=%s","tmpDiskTableUsageRate",strconv.Itoa(tmpDiskTableUsageRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.TmpDiskTableUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.TmpDiskTableUsageRate,d)
	}

	//统计数据库使用中使用磁盘临时表占使用内存临时表的占用比例Created_tmp_disk_tables/Created_tmp_tables *100% <=25%
	tmpDiskfileUsageRate, err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["Created_tmp_files"], PublicClass.GlobalStatus["Created_tmp_tables"])
	if tmpDiskfileUsageRate > 10 && err == nil {
		var d = make(map[string]string)
		d["threshold"] = ">10%"
		d["errorCode"] = "PF1-04"
		d["checkType"] = "tmpDiskfileUsageRate"
		d["checkStatus"] = "abnormal"    //异常
		d["currentValue"] = fmt.Sprintf("%s=%s","tmpDiskfileUsageRate",strconv.Itoa(tmpDiskfileUsageRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.TmpDiskfileUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.TmpDiskfileUsageRate,d)
		PublicClass.Loggs.Error(" PF1-04 Too many disk temporary file are being used. Check the slow SQL log or parameters \"tmp_table_size\" and \"max_heap_table_size\"")
	}else{
		var d = make(map[string]string)
		d["threshold"] = ">10%"
		d["errorCode"] = "PF1-04"
		d["checkType"] = "tmpDiskfileUsageRate"
		d["checkStatus"] = "normal"    //正常
		d["currentValue"] = fmt.Sprintf("%s=%s","tmpDiskfileUsageRate",strconv.Itoa(tmpDiskfileUsageRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.TmpDiskfileUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.TmpDiskfileUsageRate,d)
	}
	//统计数据库表扫描率  handler_read_rnd_next/com_select *100
	//	tableScanUsageRate,err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["handler_read_rnd_next"],PublicClass.GlobalStatus["com_select"])
	//	if tableScanUsageRate > 10 && err == nil{
	//		PublicClass.Loggs.Error("Too many disk temporary file are being used. Check the slow SQL log or parameters tmp_table_size and max_heap_table_size")
	//	}

	// 统计数据库Innodb buffer pool 使用率 100 - (Innodb_buffer_pool_pages_free * 100 / Innodb_buffer_pool_pages_total) # 单位为%
	innodbBufferPoolUsageRate, err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["Innodb_buffer_pool_pages_free"], PublicClass.GlobalStatus["Innodb_buffer_pool_pages_total"])
	if 100-innodbBufferPoolUsageRate < 80 && err == nil {
		var d = make(map[string]string)
		d["threshold"] = "<80%"
		d["errorCode"] = "PF1-05"
		d["checkType"] = "innodbBufferPoolUsageRate"
		d["checkStatus"] = "abnormal"    //异常
		d["currentValue"] = fmt.Sprintf("%s=%s","innodbBufferPoolUsageRate",strconv.Itoa(innodbBufferPoolUsageRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolUsageRate,d)
		PublicClass.Loggs.Warn(" PF1-05 The InnoDB buffer pool usage is lower than 80%")
	}else{
		var d = make(map[string]string)
		d["threshold"] = "<80%"
		d["errorCode"] = "PF1-05"
		d["checkType"] = "innodbBufferPoolUsageRate"
		d["checkStatus"] = "normal"    //正常
		d["currentValue"] = fmt.Sprintf("%s=%d","innodbBufferPoolUsageRate",strconv.Itoa(innodbBufferPoolUsageRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolUsageRate,d)
	}
	//统计数据库Innodb buffer pool 的脏页率Innodb_buffer_pool_pages_dirty * 100 / Innodb_buffer_pool_pages_total
	innodbBufferPoolDirtyPagesRate, err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["Innodb_buffer_pool_pages_dirty"], PublicClass.GlobalStatus["Innodb_buffer_pool_pages_total"])
	if innodbBufferPoolDirtyPagesRate > 50 && err == nil {
		var d = make(map[string]string)
		d["threshold"] = ">50%"
		d["errorCode"] = "PF1-06"
		d["checkType"] = "innodbBufferPoolDirtyPagesRate"
		d["checkStatus"] = "abnormal"    //异常
		d["currentValue"] = fmt.Sprintf("%s=%s","innodbBufferPoolDirtyPagesRate",strconv.Itoa(innodbBufferPoolDirtyPagesRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolDirtyPagesRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolDirtyPagesRate,d)
		PublicClass.Loggs.Warn(" PF1-06 The proportion of dirty pages in the MySQL InnoDB buffer pool exceeds \"50%\"")
	}else{
		var d = make(map[string]string)
		d["threshold"] = ">50%"
		d["errorCode"] = "PF1-06"
		d["checkType"] = "innodbBufferPoolDirtyPagesRate"
		d["checkStatus"] = "normal"    //正常
		d["currentValue"] = fmt.Sprintf("%s=%s","innodbBufferPoolDirtyPagesRate",strconv.Itoa(innodbBufferPoolDirtyPagesRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolDirtyPagesRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolDirtyPagesRate,d)
	}
	//统计数据库Innodb buffer pool的命中率Innodb_buffer_pool_reads *100 /Innodb_buffer_pool_read_requests
	innodbBufferPoolHitRate, err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["Innodb_buffer_pool_reads"], PublicClass.GlobalStatus["Innodb_buffer_pool_read_requests"])
	if 100-innodbBufferPoolHitRate < 99 && err == nil {
		var d = make(map[string]string)
		d["threshold"] = "<99%"
		d["errorCode"] = "PF1-07"
		d["checkType"] = "innodbBufferPoolHitRate"
		d["checkStatus"] = "abnormal"    //异常
		d["currentValue"] = fmt.Sprintf("%s=%s","innodbBufferPoolHitRate",strconv.Itoa(100-innodbBufferPoolHitRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolHitRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolHitRate,d)
		PublicClass.Loggs.Warn(" PF1-07 The cache hit ratio of MySQL InnoDB buffer pool is too low. You are advised to increase the size of \"innoDB buffer pool\"")
	}else{
		var d = make(map[string]string)
		d["threshold"] = "<99%"
		d["errorCode"] = "PF1-07"
		d["checkType"] = "innodbBufferPoolHitRate"
		d["checkStatus"] = "normal"    //正常
		d["currentValue"] = fmt.Sprintf("%s=%s","innodbBufferPoolHitRate",strconv.Itoa(100-innodbBufferPoolHitRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolHitRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolHitRate,d)
	}
	
	//统计数据库文件句柄使用率open_files / open_files_limit * 100% <= 75％
	openFileUsageRate, err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["open_files"], PublicClass.GlobalVariables["open_files_limit"])
	if openFileUsageRate > 75 && err == nil {
		var d = make(map[string]string)
		d["threshold"] = ">75%"
		d["errorCode"] = "PF1-08"
		d["checkType"] = "openFileUsageRate"
		d["checkStatus"] = "abnormal"    //异常
		d["currentValue"] = fmt.Sprintf("%s=%s","openFileUsageRate",strconv.Itoa(openFileUsageRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.OpenFileUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.OpenFileUsageRate,d)
		PublicClass.Loggs.Warn(" PF1-08 If the database file handle usage reaches \"75%\", you are advised to adjust the \"open_files_LIMIT\" parameter")
	}else{
		var d = make(map[string]string)
		d["threshold"] = ">75%"
		d["errorCode"] = "PF1-08"
		d["checkType"] = "openFileUsageRate"
		d["checkStatus"] = "normal"    //异常
		d["currentValue"] = fmt.Sprintf("%s=%s","openFileUsageRate",strconv.Itoa(openFileUsageRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.OpenFileUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.OpenFileUsageRate,d)
	}

	//统计数据库表打开缓存率Open_tables *100/table_open_cache
	openTableCacheUsageRate, err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["open_files"], PublicClass.GlobalVariables["open_files_limit"])
	if openTableCacheUsageRate > 80 && err == nil {
		var d = make(map[string]string)
		d["threshold"] = ">80%"
		d["errorCode"] = "PF1-09"
		d["checkType"] = "openTableCacheUsageRate"
		d["checkStatus"] = "abnormal"    //异常
		d["currentValue"] = fmt.Sprintf("%s=%s","openTableCacheUsageRate",strconv.Itoa(openTableCacheUsageRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.OpenTableCacheUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.OpenTableCacheUsageRate,d)
		PublicClass.Loggs.Warn(" PF1-09 Database open table cache usage exceeds \"80%\", you are advised to adjust the \"table_open_cache\" parameter")
	}else{
		var d = make(map[string]string)
		d["threshold"] = ">80%"
		d["errorCode"] = "PF1-09"
		d["checkType"] = "openTableCacheUsageRate"
		d["checkStatus"] = "normal"    //正常
		d["currentValue"] = fmt.Sprintf("%s=%s","openTableCacheUsageRate",strconv.Itoa(openTableCacheUsageRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.OpenTableCacheUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.OpenTableCacheUsageRate,d)
	}

	//统计数据库表缓存溢出使用率Table_open_cache_overflows *100 /(Table_open_cache_hits+Table_open_cache_misses)
	openTableTotal,err :=  PublicClass.Strea.Add(PublicClass.GlobalStatus["Table_open_cache_hits"], PublicClass.GlobalStatus["Table_open_cache_misses"])
	openTableCacheOverflowsUsageRate,err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["Table_open_cache_overflows"],openTableTotal)
	if openTableCacheOverflowsUsageRate > 10 && err == nil {
		var d = make(map[string]string)
		d["threshold"] = ">10%"
		d["errorCode"] = "PF1-10"
		d["checkType"] = "openTableCacheOverflowsUsageRate"
		d["checkStatus"] = "abnormal"    //异常
		d["currentValue"] = fmt.Sprintf("%s=%s","openTableCacheOverflowsUsageRate",strconv.Itoa(openTableCacheOverflowsUsageRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.OpenTableCacheOverflowsUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.OpenTableCacheOverflowsUsageRate,d)
		PublicClass.Loggs.Warn(" PF1-10 If the tablespace cache overflow usage is greater than \"10%\", you are advised to adjust the \"table_open_cache\" parameter")
	}else{
		var d = make(map[string]string)
		d["threshold"] = ">10%"
		d["errorCode"] = "PF1-10"
		d["checkType"] = "openTableCacheOverflowsUsageRate"
		d["checkStatus"] = "normal"    //正常
		d["currentValue"] = fmt.Sprintf("%s=%s","openTableCacheOverflowsUsageRate",strconv.Itoa(openTableCacheOverflowsUsageRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.OpenTableCacheOverflowsUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.OpenTableCacheOverflowsUsageRate,d)
	}

	//统计数据库全表扫描的占比率Select_scan *100 /Queries
	selectScanUsageRate,err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["Select_scan"],PublicClass.GlobalStatus["Queries"])
	if selectScanUsageRate > 10 && err == nil {
		var d = make(map[string]string)
		d["threshold"] = ">10%"
		d["errorCode"] = "PF1-11"
		d["checkType"] = "selectScanUsageRate"
		d["checkStatus"] = "abnormal"    //异常
		d["currentValue"] = fmt.Sprintf("%s=%s","selectScanUsageRate",strconv.Itoa(selectScanUsageRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.SelectScanUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.SelectScanUsageRate,d)
		PublicClass.Loggs.Warn(" PF1-11 The database does not use indexes. If the full table scan usage exceeds \"10%\", You are advised to check the slow SQL")
	}else{
		var d = make(map[string]string)
		d["threshold"] = ">10%"
		d["errorCode"] = "PF1-11"
		d["checkType"] = "selectScanUsageRate"
		d["checkStatus"] = "normal"    //正常
		d["currentValue"] = fmt.Sprintf("%s=%s","selectScanUsageRate",strconv.Itoa(selectScanUsageRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.SelectScanUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.SelectScanUsageRate,d)
	}
	//统计数据库join语句发生全表扫描占比率Select_full_join *100 /Queries
	selectfullJoinScanUsageRate,err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["Select_full_join"],PublicClass.GlobalStatus["Queries"])
	if selectfullJoinScanUsageRate > 10 && err == nil {
		var d = make(map[string]string)
		d["threshold"] = ">10%"
		d["errorCode"] = "PF1-12"
		d["checkType"] = "selectfullJoinScanUsageRate"
		d["checkStatus"] = "abnormal"    //异常
		d["currentValue"] = fmt.Sprintf("%s=%s","selectfullJoinScanUsageRate",strconv.Itoa(selectfullJoinScanUsageRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.SelectfullJoinScanUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.SelectfullJoinScanUsageRate,d)
		PublicClass.Loggs.Warn(" PF1-12 The database uses the JOIN statement and the non-driver table does not use the index. The full table scan usage is greater than \"10%\". You are advised to check for slow SQL")
	}else{
		var d = make(map[string]string)
		d["threshold"] = ">10%"
		d["errorCode"] = "PF1-12"
		d["checkType"] = "selectfullJoinScanUsageRate"
		d["checkStatus"] = "normal"    //正常
		d["currentValue"] = fmt.Sprintf("%s=%s","selectfullJoinScanUsageRate",strconv.Itoa(selectfullJoinScanUsageRate))
		PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.SelectfullJoinScanUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceStatus.SelectfullJoinScanUsageRate,d)
	}

}

func (baselineCheck *DatabaseBaselineCheckStruct) DatabasePerformanceTableIndexCheck() {
	nowDateTime := PublicClass.DBexecInter.DBQueryDateString("select now() as datetime")
	var tmpTableColumnMap = make(map[string]string)
	for i := range PublicClass.InformationSchemaColumnsData{
		c := PublicClass.InformationSchemaColumnsData[i]
		//过滤自增主键为int类型的库表，并输出库、表、数据类型等信息
		if c["EXTRA"] != nil && c["EXTRA"] == "auto_increment" && c["DATA_TYPE"] == "int"{
			keyDT := fmt.Sprintf("%s@%s",c["TABLE_SCHEMA"],c["TABLE_NAME"])
			valCT := fmt.Sprintf("%s@%s@%s",c["COLUMN_NAME"],c["DATA_TYPE"],c["COLUMN_TYPE"])
			tmpTableColumnMap[keyDT] = valCT
		}
	}

	for _,v := range PublicClass.InformationSchemaTablesData {
		var d = make(map[string]string)
		d["database"] = v["TABLE_SCHEMA"].(string)
		d["tableName"] = v["TABLE_NAME"].(string)

		databaseTableName := fmt.Sprintf("%s@%s",v["TABLE_SCHEMA"],v["TABLE_NAME"])
		if _,ok := tmpTableColumnMap[databaseTableName];ok{
			e := newMap(d)
			tmpColumnInfoSliect := strings.Split(tmpTableColumnMap[databaseTableName],"@")
			tmpColumnName := tmpColumnInfoSliect[0]
			tmpColumnType := tmpColumnInfoSliect[2]
			e["columnName"] = tmpColumnName
			e["columnType"] = tmpColumnType
			e["threshold"] = ">85"
			e["errorCode"] = "PF2-01"
			e["checkType"] = "tableAutoPrimaryKeyUsageRate"
			e["autoIncrement"] = strconv.Itoa(int(v["AUTO_INCREMENT"].(int64)))
			//检查是否存在自增主键溢出风险。统计数据库自增id列快要溢出的表
			if strings.Contains(tmpColumnType,"unsigned"){
				if unsignedIntUsageRate,err := PublicClass.Strea.Percentage(v["AUTO_INCREMENT"],MunsignedInt);err ==nil && unsignedIntUsageRate >=85{
					e["checkStatus"] = "abnormal"
					e["currentValue"] = fmt.Sprintf("%s.%s",v["TABLE_SCHEMA"],v["TABLE_NAME"])
					PublicClass.Loggs.Warn(fmt.Sprintf(" PF2-01 The self-value-added usage of tables in the database exceeds \"85%%\", causing data type overflow risks. The details are as follows: Database: \"%v\", table name: \"%v\", increment column name: \"%v\", increment column data type: \"%v\", current increment values: \"%v\"",v["TABLE_SCHEMA"],v["TABLE_NAME"],tmpColumnName,tmpColumnType,v["AUTO_INCREMENT"]))
				} else{
					e["checkStatus"] = "normal"
				}
				PublicClass.InspectionResult.DatabasePerformance.PerformanceTableIndex.TableAutoPrimaryKeyUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceTableIndex.TableAutoPrimaryKeyUsageRate,e)
			}else{
				if intUsageRate,err := PublicClass.Strea.Percentage(v["AUTO_INCREMENT"],Mint);err ==nil && intUsageRate >=85{
					e["checkStatus"] = "abnormal"
					PublicClass.Loggs.Warn(fmt.Sprintf(" PF2-01 The self-value-added usage of tables in the database exceeds 85%%, causing data type overflow risks. The details are as follows: Database: \"%v\", table name: \"%v\", increment column name: \"%v\", increment column data type: \"%v\", current increment values: \"%v\"",v["TABLE_SCHEMA"],v["TABLE_NAME"],tmpColumnName,tmpColumnType,v["AUTO_INCREMENT"]))
				}else{
					e["checkStatus"] = "normal"
				}
				PublicClass.InspectionResult.DatabasePerformance.PerformanceTableIndex.TableAutoPrimaryKeyUsageRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceTableIndex.TableAutoPrimaryKeyUsageRate,e)
			}
		}
		//单表行数大于500w，且平均行长大于10KB。
		m := newMap(d)
		tableRows ,_ := strconv.Atoi(fmt.Sprintf("%s",v["TABLE_ROWS"]))
		avgRowLength ,_ := strconv.Atoi(fmt.Sprintf("%s",v["AVG_ROW_LENGTH"]))
		m["threshold"] = ">500W"
		m["errorCode"] = "PF2-02"
		m["checkType"] = "tableRows"
		m["tableRows"] = strconv.Itoa(tableRows)
		m["avgRowLength"] = strconv.Itoa(avgRowLength)
		if tableRows > 5000000 && avgRowLength/1024 > 10 {
			m["checkStatus"] = "abnormal"
			m["currentValue"] = fmt.Sprintf("%s.%s",v["TABLE_SCHEMA"],v["TABLE_NAME"])
			PublicClass.InspectionResult.DatabasePerformance.PerformanceTableIndex.TableRows = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceTableIndex.TableRows, m)
			PublicClass.Loggs.Warn(fmt.Sprintf(" PF2-02 The current table is a large table if the number of rows is greater than 5 million and the average line length is greater than 10KB. The details are as follows: Database: \"%v\", table name: \"%v\", tableRows: \"%v\", avgRowLength：\"%d\"",v["TABLE_SCHEMA"],v["TABLE_NAME"],tableRows,avgRowLength/1024))
		} else {
			m["checkStatus"] = "normal"
			PublicClass.InspectionResult.DatabasePerformance.PerformanceTableIndex.TableRows = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceTableIndex.TableRows, m)
		}

		// 单表大于6G，并且碎片率大于30%。
		n := newMap(d)
		var dataLength,indexLength,dataFree int
		if v["DATA_FREE"] != nil{ dataLength = int(v["DATA_FREE"].(int64)) }
		if v["INDEX_LENGTH"] != nil{ indexLength = int(v["INDEX_LENGTH"].(int64))}
		if v["DATA_FREE"] != nil { dataFree = int(v["DATA_FREE"].(int64)) }
		dataLengthTotal := dataLength + indexLength //表空间
		n["threshold"] = ">30%"
		n["errorCode"] = "PF2-03"
		n["checkType"] = "diskFragmentationRate"
		if diskFragmentationRate, err := PublicClass.Strea.Percentage(dataFree, dataLengthTotal); diskFragmentationRate > 1 && err == nil && dataLengthTotal/1024/1024/1024 > 1 {
			n["checkStatus"] = "abnormal"
			n["currentValue"] = fmt.Sprintf("%s.%s",v["TABLE_SCHEMA"],v["TABLE_NAME"])
			PublicClass.InspectionResult.DatabasePerformance.PerformanceTableIndex.DiskFragmentationRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceTableIndex.DiskFragmentationRate, n)
			PublicClass.Loggs.Warn(fmt.Sprintf(" PF2-03 If the current tablespace contains more than 6 GB and the disk fragmentation rate is greater than 30%%, you are advised to run THE ALTER command to delete disk fragmentation. The details are as follows: Database: \"%v\", table name: \"%v\", Table space size: \"%dG\", diskFragmentationRate：\"%d\"",v["TABLE_SCHEMA"],v["TABLE_NAME"],dataLengthTotal/1024/1024/1024,diskFragmentationRate))
		} else {
			n["diskFragmentationRate"] = strconv.Itoa(diskFragmentationRate)
			n["checkStatus"] = "normal"
			PublicClass.InspectionResult.DatabasePerformance.PerformanceTableIndex.DiskFragmentationRate = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceTableIndex.DiskFragmentationRate, n)
		}

		//单表行数大于1000W，且表空间大于30G
		o := newMap(d)
		o["threshold"] = ">30G"
		o["errorCode"] = "PF2-04"
		o["checkType"] = "bigTable"
		o["dataLengthTotal"] = strconv.Itoa(dataLengthTotal / 1024 / 1024 / 1024)
		o["tableRows"] = strconv.Itoa(tableRows)
		if dataLengthTotal/1024/1024/1024 > 30 && tableRows > 10000000 {
			o["checkStatus"] = "abnormal"
			o["currentValue"] = fmt.Sprintf("%s.%s",v["TABLE_SCHEMA"],v["TABLE_NAME"])
			PublicClass.InspectionResult.DatabasePerformance.PerformanceTableIndex.BigTable = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceTableIndex.BigTable, o)
			PublicClass.Loggs.Warn(fmt.Sprintf(" PF2-04 If the number of rows in the current table is greater than 1000W and the tablespace is greater than 30G, the table belongs to a large table. Recommended Attention Table. The details are as follows: Database: \"%v\", table name: \"%v\", tableRows：\"%d\", Table space size: \"%dG\"",v["TABLE_SCHEMA"],v["TABLE_NAME"],tableRows,dataLengthTotal/1024/1024/1024))
		} else {
			o["checkStatus"] = "normal"
			PublicClass.InspectionResult.DatabasePerformance.PerformanceTableIndex.BigTable = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceTableIndex.BigTable, o)
		}

		//检查一个星期内未更新的表
		var tableUpdateTime string
		if v["CREATE_TIME"] != nil {
			if v["v[UPDATE_TIME"] == nil {
				tableUpdateTime = v["CREATE_TIME"].(string)
			} else {
				tableUpdateTime = v["UPDATE_TIME"].(string)
			}
		}
		arrDay, _ := PublicClass.Strea.GetTimeDayArr(tableUpdateTime, nowDateTime)
		p := newMap(d)
		p["threshold"] = ">7"
		p["errorCode"] = "PF2-05"
		p["checkType"] = "coldTable"
		if arrDay > 7 {
			p["checkStatus"] = "abnormal"
			p["currentValue"] = fmt.Sprintf("%s.%s",v["TABLE_SCHEMA"],v["TABLE_NAME"])
			PublicClass.InspectionResult.DatabasePerformance.PerformanceTableIndex.ColdTable = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceTableIndex.ColdTable, p)
			PublicClass.Loggs.Error(fmt.Sprintf(" PF2-05 The current table has not been updated for seven days (no DML has occurred against the table). The details are as follows: Database: \"%v\", table name: \"%v\", lasterUpdateTime：\"%v\" ",v["TABLE_SCHEMA"],v["TABLE_NAME"],tableUpdateTime))
		} else {

			p["checkStatus"] = "normal"
			PublicClass.InspectionResult.DatabasePerformance.PerformanceTableIndex.ColdTable = append(PublicClass.InspectionResult.DatabasePerformance.PerformanceTableIndex.ColdTable, p)
		}
	}
}

