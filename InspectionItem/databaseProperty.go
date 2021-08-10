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
		d["checkType"] = "binlogDiskUsageRate"
		d["checkStatus"] = "abnormal"    //异常
		d["binlogDiskUsageRate"] = strconv.Itoa(binlogDiskUsageRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.BinlogDiskUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.BinlogDiskUsageRate,d)
		PublicClass.Loggs.Error("The current database binlog is using too many disk writes. It is recommended to modify the binlog_cache_size parameter")
	}else{
		var d = make(map[string]string)
		d["checkType"] = "binlogDiskUsageRate"
		d["checkStatus"] = "normal"    //正常
		d["binlogDiskUsageRate"] = strconv.Itoa(binlogDiskUsageRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.BinlogDiskUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.BinlogDiskUsageRate,d)
	}
	
	//统计历史连接数最大使用率，使用创建过
	historyConnectionMaxUsageRate, err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["Threads_created"],PublicClass.GlobalVariables["max_connections"])
	if historyConnectionMaxUsageRate > 80 && err == nil {
		var d = make(map[string]string)
		d["checkType"] = "historyConnectionMaxUsageRate"
		d["checkStatus"] = "abnormal"    //异常
		d["binlogDiskUsageRate"] = strconv.Itoa(historyConnectionMaxUsageRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.HistoryConnectionMaxUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.HistoryConnectionMaxUsageRate,d)
		PublicClass.Loggs.Error("If the maximum usage of historical database connections exceeds 80%, change the max_connections value and check services")
	}else{
		var d = make(map[string]string)
		d["checkType"] = "historyConnectionMaxUsageRate"
		d["checkStatus"] = "normal"    //正常
		d["binlogDiskUsageRate"] = strconv.Itoa(historyConnectionMaxUsageRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.HistoryConnectionMaxUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.HistoryConnectionMaxUsageRate,d)
	}
	
	//统计数据库使用中使用磁盘临时表占使用内存临时表的占用比例Created_tmp_disk_tables/Created_tmp_tables *100% <=25%
	tmpDiskTableUsageRate, err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["Created_tmp_disk_tables"], PublicClass.GlobalStatus["Created_tmp_tables"])
	if tmpDiskTableUsageRate > 25 && err == nil {
		var d = make(map[string]string)
		d["checkType"] = "tmpDiskTableUsageRate"
		d["checkStatus"] = "abnormal"    //异常
		d["binlogDiskUsageRate"] = strconv.Itoa(tmpDiskTableUsageRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.TmpDiskTableUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.TmpDiskTableUsageRate,d)
		PublicClass.Loggs.Error("Too many disk temporary tables are being used. Check the slow SQL log or parameters tmp_table_size and max_heap_table_size")
	}else{
		var d = make(map[string]string)
		d["checkType"] = "tmpDiskTableUsageRate"
		d["checkStatus"] = "normal"    //正常
		d["binlogDiskUsageRate"] = strconv.Itoa(tmpDiskTableUsageRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.TmpDiskTableUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.TmpDiskTableUsageRate,d)
	}

	//统计数据库使用中使用磁盘临时表占使用内存临时表的占用比例Created_tmp_disk_tables/Created_tmp_tables *100% <=25%
	tmpDiskfileUsageRate, err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["Created_tmp_files"], PublicClass.GlobalStatus["Created_tmp_tables"])
	if tmpDiskfileUsageRate > 10 && err == nil {
		var d = make(map[string]string)
		d["checkType"] = "tmpDiskfileUsageRate"
		d["checkStatus"] = "abnormal"    //异常
		d["binlogDiskUsageRate"] = strconv.Itoa(tmpDiskfileUsageRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.TmpDiskfileUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.TmpDiskfileUsageRate,d)
		PublicClass.Loggs.Error("Too many disk temporary file are being used. Check the slow SQL log or parameters tmp_table_size and max_heap_table_size")
	}else{
		var d = make(map[string]string)
		d["checkType"] = "tmpDiskfileUsageRate"
		d["checkStatus"] = "normal"    //正常
		d["binlogDiskUsageRate"] = strconv.Itoa(tmpDiskfileUsageRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.TmpDiskfileUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.TmpDiskfileUsageRate,d)
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
		d["checkType"] = "innodbBufferPoolUsageRate"
		d["checkStatus"] = "abnormal"    //异常
		d["binlogDiskUsageRate"] = strconv.Itoa(innodbBufferPoolUsageRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.InnodbBufferPoolUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.InnodbBufferPoolUsageRate,d)
		PublicClass.Loggs.Warn("The InnoDB buffer pool usage is lower than 80%")
	}else{
		var d = make(map[string]string)
		d["checkType"] = "innodbBufferPoolUsageRate"
		d["checkStatus"] = "normal"    //正常
		d["binlogDiskUsageRate"] = strconv.Itoa(innodbBufferPoolUsageRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.InnodbBufferPoolUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.InnodbBufferPoolUsageRate,d)
	}
	//统计数据库Innodb buffer pool 的脏页率Innodb_buffer_pool_pages_dirty * 100 / Innodb_buffer_pool_pages_total
	innodbBufferPoolDirtyPagesRate, err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["Innodb_buffer_pool_pages_dirty"], PublicClass.GlobalStatus["Innodb_buffer_pool_pages_total"])
	if innodbBufferPoolDirtyPagesRate > 50 && err == nil {
		var d = make(map[string]string)
		d["checkType"] = "innodbBufferPoolDirtyPagesRate"
		d["checkStatus"] = "abnormal"    //异常
		d["binlogDiskUsageRate"] = strconv.Itoa(innodbBufferPoolDirtyPagesRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.InnodbBufferPoolDirtyPagesRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.InnodbBufferPoolDirtyPagesRate,d)
		PublicClass.Loggs.Warn("The proportion of dirty pages in the MySQL InnoDB buffer pool exceeds 50%")
	}else{
		var d = make(map[string]string)
		d["checkType"] = "innodbBufferPoolDirtyPagesRate"
		d["checkStatus"] = "normal"    //正常
		d["binlogDiskUsageRate"] = strconv.Itoa(innodbBufferPoolDirtyPagesRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.InnodbBufferPoolDirtyPagesRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.InnodbBufferPoolDirtyPagesRate,d)
	}
	//统计数据库Innodb buffer pool的命中率Innodb_buffer_pool_reads *100 /Innodb_buffer_pool_read_requests
	innodbBufferPoolHitRate, err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["Innodb_buffer_pool_reads"], PublicClass.GlobalStatus["Innodb_buffer_pool_read_requests"])
	if 100-innodbBufferPoolHitRate < 99 && err == nil {
		var d = make(map[string]string)
		d["checkType"] = "innodbBufferPoolHitRate"
		d["checkStatus"] = "abnormal"    //异常
		d["binlogDiskUsageRate"] = strconv.Itoa(innodbBufferPoolHitRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.InnodbBufferPoolHitRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.InnodbBufferPoolHitRate,d)
		PublicClass.Loggs.Warn("The cache hit ratio of MySQL InnoDB buffer pool is too low. You are advised to increase the size of innoDB buffer pool")
	}else{
		var d = make(map[string]string)
		d["checkType"] = "innodbBufferPoolHitRate"
		d["checkStatus"] = "normal"    //正常
		d["binlogDiskUsageRate"] = strconv.Itoa(innodbBufferPoolHitRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.InnodbBufferPoolHitRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.InnodbBufferPoolHitRate,d)
	}
	
	//统计数据库文件句柄使用率open_files / open_files_limit * 100% <= 75％
	openFileUsageRate, err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["open_files"], PublicClass.GlobalVariables["open_files_limit"])
	if openFileUsageRate > 75 && err == nil {
		var d = make(map[string]string)
		d["checkType"] = "openFileUsageRate"
		d["checkStatus"] = "abnormal"    //异常
		d["binlogDiskUsageRate"] = strconv.Itoa(openFileUsageRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.OpenFileUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.OpenFileUsageRate,d)
		PublicClass.Loggs.Warn("If the database file handle usage reaches 75%, you are advised to adjust the open_files_LIMIT parameter")
	}else{
		var d = make(map[string]string)
		d["checkType"] = "openFileUsageRate"
		d["checkStatus"] = "normal"    //异常
		d["binlogDiskUsageRate"] = strconv.Itoa(openFileUsageRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.OpenFileUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.OpenFileUsageRate,d)
	}

	//统计数据库表打开缓存率Open_tables *100/table_open_cache
	openTableCacheUsageRate, err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["open_files"], PublicClass.GlobalVariables["open_files_limit"])
	if openTableCacheUsageRate > 80 && err == nil {
		var d = make(map[string]string)
		d["checkType"] = "openTableCacheUsageRate"
		d["checkStatus"] = "abnormal"    //异常
		d["binlogDiskUsageRate"] = strconv.Itoa(openTableCacheUsageRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.OpenTableCacheUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.OpenTableCacheUsageRate,d)
		PublicClass.Loggs.Warn("Database open table cache usage exceeds 80%, you are advised to adjust the table_open_cache parameter")
	}else{
		var d = make(map[string]string)
		d["checkType"] = "openTableCacheUsageRate"
		d["checkStatus"] = "normal"    //正常
		d["binlogDiskUsageRate"] = strconv.Itoa(openTableCacheUsageRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.OpenTableCacheUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.OpenTableCacheUsageRate,d)
	}

	//统计数据库表缓存溢出使用率Table_open_cache_overflows *100 /(Table_open_cache_hits+Table_open_cache_misses)
	openTableTotal,err :=  PublicClass.Strea.Add(PublicClass.GlobalStatus["Table_open_cache_hits"], PublicClass.GlobalStatus["Table_open_cache_misses"])
	openTableCacheOverflowsUsageRate,err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["Table_open_cache_overflows"],openTableTotal)
	if openTableCacheOverflowsUsageRate > 10 && err == nil {
		var d = make(map[string]string)
		d["checkType"] = "openTableCacheOverflowsUsageRate"
		d["checkStatus"] = "abnormal"    //异常
		d["binlogDiskUsageRate"] = strconv.Itoa(openTableCacheOverflowsUsageRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.OpenTableCacheOverflowsUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.OpenTableCacheOverflowsUsageRate,d)
		PublicClass.Loggs.Warn("If the tablespace cache overflow usage is greater than 10%, you are advised to adjust the table_open_cache parameter")
	}else{
		var d = make(map[string]string)
		d["checkType"] = "openTableCacheOverflowsUsageRate"
		d["checkStatus"] = "normal"    //正常
		d["binlogDiskUsageRate"] = strconv.Itoa(openTableCacheOverflowsUsageRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.OpenTableCacheOverflowsUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.OpenTableCacheOverflowsUsageRate,d)
	}

	//统计数据库全表扫描的占比率Select_scan *100 /Queries
	selectScanUsageRate,err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["Select_scan"],PublicClass.GlobalStatus["Queries"])
	if selectScanUsageRate > 10 && err == nil {
		var d = make(map[string]string)
		d["checkType"] = "selectScanUsageRate"
		d["checkStatus"] = "abnormal"    //异常
		d["binlogDiskUsageRate"] = strconv.Itoa(selectScanUsageRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.SelectScanUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.SelectScanUsageRate,d)
		PublicClass.Loggs.Warn("The database does not use indexes. If the full table scan usage exceeds 10%, You are advised to check the slow SQL")
	}else{
		var d = make(map[string]string)
		d["checkType"] = "selectScanUsageRate"
		d["checkStatus"] = "normal"    //正常
		d["binlogDiskUsageRate"] = strconv.Itoa(selectScanUsageRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.SelectScanUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.SelectScanUsageRate,d)
	}
	//统计数据库join语句发生全表扫描占比率Select_full_join *100 /Queries
	selectfullJoinScanUsageRate,err := PublicClass.Strea.Percentage(PublicClass.GlobalStatus["Select_full_join"],PublicClass.GlobalStatus["Queries"])
	if selectfullJoinScanUsageRate > 10 && err == nil {
		var d = make(map[string]string)
		d["checkType"] = "selectfullJoinScanUsageRate"
		d["checkStatus"] = "abnormal"    //异常
		d["binlogDiskUsageRate"] = strconv.Itoa(selectfullJoinScanUsageRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.SelectfullJoinScanUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.SelectfullJoinScanUsageRate,d)
		PublicClass.Loggs.Warn("The database uses the JOIN statement and the non-driver table does not use the index. The full table scan usage is greater than 10%. You are advised to check for slow SQL")
	}else{
		var d = make(map[string]string)
		d["checkType"] = "selectfullJoinScanUsageRate"
		d["checkStatus"] = "normal"    //正常
		d["binlogDiskUsageRate"] = strconv.Itoa(selectfullJoinScanUsageRate)
		PublicClass.InspectionResult.DatabasePerformanceCheck.SelectfullJoinScanUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceCheck.SelectfullJoinScanUsageRate,d)
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
		databaseTableName := fmt.Sprintf("%s@%s",v["TABLE_SCHEMA"],v["TABLE_NAME"])
		if _,ok := tmpTableColumnMap[databaseTableName];ok{
			tmpColumnInfoSliect := strings.Split(tmpTableColumnMap[databaseTableName],"@")
			tmpColumnName := tmpColumnInfoSliect[0]
			tmpColumnType := tmpColumnInfoSliect[2]
			d["columnName"] = tmpColumnName
			d["columnType"] = tmpColumnType
			d["checkType"] = "tableAutoPrimaryKeyUsageRate"
			d["autoIncrement"] = strconv.Itoa(int(v["AUTO_INCREMENT"].(int64)))
			//检查是否存在自增主键溢出风险。统计数据库自增id列快要溢出的表
			if strings.Contains(tmpColumnType,"unsigned"){
				if unsignedIntUsageRate,err := PublicClass.Strea.Percentage(v["AUTO_INCREMENT"],MunsignedInt);err ==nil && unsignedIntUsageRate >=1{
					d["checkStatus"] = "abnormal"
					PublicClass.Loggs.Warn(fmt.Sprintf("The self-value-added usage of tables in the database exceeds 85%%, causing data type overflow risks. The details are as follows: Database: \"%v\", table name: \"%v\", increment column name: \"%v\", increment column data type: \"%v\", current increment column: \"%v\"",v["TABLE_SCHEMA"],v["TABLE_NAME"],v["column_name"],v["COLUMN_TYPE"],v["AUTO_INCREMENT"]))
				} else{
					d["checkStatus"] = "normal"
				}
				PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.TableAutoPrimaryKeyUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.TableAutoPrimaryKeyUsageRate,d)
			}else{
				if intUsageRate,err := PublicClass.Strea.Percentage(v["AUTO_INCREMENT"],Mint);err ==nil && intUsageRate >=1{
					d["checkStatus"] = "abnormal"
					PublicClass.Loggs.Warn(fmt.Sprintf("The self-value-added usage of tables in the database exceeds 85%%, causing data type overflow risks. The details are as follows: Database: \"%v\", table name: \"%v\", increment column name: \"%v\", increment column data type: \"%v\", current increment column: \"%v\"",v["TABLE_SCHEMA"],v["TABLE_NAME"],v["column_name"],v["COLUMN_TYPE"],v["AUTO_INCREMENT"]))
				}else{
					d["checkStatus"] = "normal"
				}
				PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.TableAutoPrimaryKeyUsageRate = append(PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.TableAutoPrimaryKeyUsageRate,d)
			}
		}
		m := newMap(d)
		//单表行数大于500w，且平均行长大于10KB。
		tableRows ,_ := strconv.Atoi(fmt.Sprintf("%s",v["TABLE_ROWS"]))
		avgRowLength ,_ := strconv.Atoi(fmt.Sprintf("%s",v["AVG_ROW_LENGTH"]))
		if tableRows > 5000000 && avgRowLength/1024 > 10 {
			m["checkType"] = "tableRows"
			m["tableRows"] = strconv.Itoa(tableRows)
			m["avgRowLength"] = strconv.Itoa(avgRowLength)
			m["checkStatus"] = "abnormal"
			PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.TableRows = append(PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.TableRows, m)
		} else {
			m["checkStatus"] = "normal"
			m["checkType"] = "tableRows"
			m["tableRows"] = strconv.Itoa(tableRows)
			m["avgRowLength"] = strconv.Itoa(avgRowLength)
			PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.TableRows = append(PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.TableRows, m)
		}
		// 单表大于6G，并且碎片率大于30%。
		n := newMap(d)
		var dataLength,indexLength,dataFree int
		if v["DATA_FREE"] != nil{ dataLength = int(v["DATA_FREE"].(int64)) }
		if v["INDEX_LENGTH"] != nil{ indexLength = int(v["INDEX_LENGTH"].(int64))}
		if v["DATA_FREE"] != nil { dataFree = int(v["DATA_FREE"].(int64)) }
		dataLengthTotal := dataLength + indexLength //表空间
		if diskFragmentationRate, err := PublicClass.Strea.Percentage(dataFree, dataLengthTotal); diskFragmentationRate > 1 && err == nil && dataLengthTotal/1024/1024/1024 > 1 {
			n["checkType"] = "diskFragmentationRate"
			n["dataLengthTotal"] = strconv.Itoa(int(dataLengthTotal))
			n["diskFragmentationRate"] = strconv.Itoa(diskFragmentationRate)
			n["checkStatus"] = "abnormal"
			PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.DiskFragmentationRate = append(PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.DiskFragmentationRate, n)
		} else {
			n["checkType"] = "diskFragmentationRate"
			n["dataLengthTotal"] = strconv.Itoa(int(dataLengthTotal))
			n["diskFragmentationRate"] = strconv.Itoa(diskFragmentationRate)
			n["checkStatus"] = "normal"
			PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.DiskFragmentationRate = append(PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.DiskFragmentationRate, n)
		}

		//单表行数大于1000W，且表空间大于30G
		o := newMap(d)
		if dataLengthTotal/1024/1024/1024 > 30 && tableRows > 10000000 {
			o["checkType"] = "bigTable"
			o["dataLengthTotal"] = strconv.Itoa(dataLengthTotal / 1024 / 1024 / 1024)
			o["tableRows"] = strconv.Itoa(tableRows)
			o["checkStatus"] = "abnormal"
			PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.BigTable = append(PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.BigTable, o)
		} else {
			o["checkType"] = "bigTable"
			o["dataLengthTotal"] = strconv.Itoa(dataLengthTotal / 1024 / 1024 / 1024)
			o["tableRows"] = strconv.Itoa(tableRows)
			o["checkStatus"] = "normal"
			PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.BigTable = append(PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.BigTable, o)
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
		if arrDay > 7 {
			p["checkType"] = "coldTable"
			p["checkStatus"] = "abnormal"
			PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.ColdTable = append(PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.ColdTable, p)
			PublicClass.Loggs.Error("有大于7天的表")
		} else {
			p["checkType"] = "coldTable"
			p["checkStatus"] = "normal"
			PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.ColdTable = append(PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.ColdTable, p)
		}
	}
}

