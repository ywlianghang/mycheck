package InspectionItem

import (
	"DepthInspection/api/PublicClass"
	"fmt"
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

func (baselineCheck *DatabaseBaselineCheckStruct) DatabaseBinlogdesign(aa *PublicClass.ConfigInfo) {
	strSql1 := fmt.Sprintf("show global status")
	cc1 := aa.DatabaseExecInterf.DBQueryDateMap(aa, strSql1)
	strSql2 := fmt.Sprintf("show global variables")
	dd1 := aa.DatabaseExecInterf.DBQueryDateMap(aa, strSql2)
	//统计使用磁盘的binlog写入占使用内存buffer的binlog写入的百分比，大于100%则需要增加binlog_cache_size
	binlogDiskUsageRate, err := aa.Streamm.Percentage(cc1["Binlog_cache_disk_use"], cc1["Binlog_cache_use"])
	if binlogDiskUsageRate > 100 && err == nil {
		aa.Loggs.Error("The current database binlog is using too many disk writes. It is recommended to modify the binlog_cache_size parameter")
	}
	//统计历史连接数最大使用率，使用创建过
	historyConnectionMaxUsageRate, err := aa.Streamm.Percentage(cc1["Threads_created"], dd1["max_connections"])
	if historyConnectionMaxUsageRate > 80 && err == nil {
		aa.Loggs.Error("If the maximum usage of historical database connections exceeds 80%, change the max_connections value and check services")
	}
	//统计数据库使用中使用磁盘临时表占使用内存临时表的占用比例Created_tmp_disk_tables/Created_tmp_tables *100% <=25%
	tmpDiskTableUsageRate, err := aa.Streamm.Percentage(cc1["Created_tmp_disk_tables"], cc1["Created_tmp_tables"])
	if tmpDiskTableUsageRate > 25 && err == nil {
		aa.Loggs.Error("Too many disk temporary tables are being used. Check the slow SQL log or parameters tmp_table_size and max_heap_table_size")
	}
	//统计数据库使用中使用磁盘临时表占使用内存临时表的占用比例Created_tmp_disk_tables/Created_tmp_tables *100% <=25%
	tmpDiskfileUsageRate, err := aa.Streamm.Percentage(cc1["Created_tmp_files"], cc1["Created_tmp_tables"])
	if tmpDiskfileUsageRate > 10 && err == nil {
		aa.Loggs.Error("Too many disk temporary file are being used. Check the slow SQL log or parameters tmp_table_size and max_heap_table_size")
	}
	//统计数据库表扫描率  handler_read_rnd_next/com_select *100
	//	tableScanUsageRate,err := aa.Streamm.Percentage(cc1["handler_read_rnd_next"],cc1["com_select"])
	//	if tableScanUsageRate > 10 && err == nil{
	//		aa.Loggs.Error("Too many disk temporary file are being used. Check the slow SQL log or parameters tmp_table_size and max_heap_table_size")
	//	}
	//统计数据库Innodb buffer pool 使用率 100 - (Innodb_buffer_pool_pages_free * 100 / Innodb_buffer_pool_pages_total) # 单位为%
	innodbBufferPoolUsageRate, err := aa.Streamm.Percentage(cc1["Innodb_buffer_pool_pages_free"], cc1["Innodb_buffer_pool_pages_total"])
	if 100-innodbBufferPoolUsageRate < 80 && err == nil {
		aa.Loggs.Warn("The InnoDB buffer pool usage is lower than 80%")
	}
	//统计数据库Innodb buffer pool 的脏页率Innodb_buffer_pool_pages_dirty * 100 / Innodb_buffer_pool_pages_total
	innodbBufferPoolDirtyPagesRate, err := aa.Streamm.Percentage(cc1["Innodb_buffer_pool_pages_dirty"], cc1["Innodb_buffer_pool_pages_total"])
	if innodbBufferPoolDirtyPagesRate > 50 && err == nil {
		aa.Loggs.Warn("The proportion of dirty pages in the MySQL InnoDB buffer pool exceeds 50%")
	}
	//统计数据库Innodb buffer pool的命中率Innodb_buffer_pool_reads *100 /Innodb_buffer_pool_read_requests
	innodbBufferPoolHitRate, err := aa.Streamm.Percentage(cc1["Innodb_buffer_pool_reads"], cc1["Innodb_buffer_pool_read_requests"])
	if 100-innodbBufferPoolHitRate < 99 && err == nil {
		aa.Loggs.Warn("The cache hit ratio of MySQL InnoDB buffer pool is too low. You are advised to increase the size of innoDB buffer pool")
	}
	//统计数据库文件句柄使用率open_files / open_files_limit * 100% <= 75％
	openFileUsageRate, err := aa.Streamm.Percentage(cc1["open_files"], dd1["open_files_limit"])
	if openFileUsageRate > 75 && err == nil {
		aa.Loggs.Warn("If the database file handle usage reaches 75%, you are advised to adjust the open_files_LIMIT parameter")
	}
	//统计数据库表打开缓存率Open_tables *100/table_open_cache
	openTableCacheUsageRate, err := aa.Streamm.Percentage(cc1["open_files"], dd1["open_files_limit"])
	if openTableCacheUsageRate > 80 && err == nil {
		aa.Loggs.Warn("Database open table cache usage exceeds 80%, you are advised to adjust the table_open_cache parameter")
	}
	//统计数据库表缓存溢出使用率Table_open_cache_overflows *100 /(Table_open_cache_hits+Table_open_cache_misses)
	openTableTotal,err :=  aa.Streamm.Add(cc1["Table_open_cache_hits"], cc1["Table_open_cache_misses"])
	openTableCacheOverflowsUsageRate,err := aa.Streamm.Percentage(cc1["Table_open_cache_overflows"],openTableTotal)
	if openTableCacheOverflowsUsageRate > 10 && err == nil {
		aa.Loggs.Warn("If the tablespace cache overflow usage is greater than 10%, you are advised to adjust the table_open_cache parameter")
	}
	//统计数据库全表扫描的占比率Select_scan *100 /Queries
	selectScanUsageRate,err := aa.Streamm.Percentage(cc1["Select_scan"],cc1["Queries"])
	if selectScanUsageRate > 10 && err == nil {
		aa.Loggs.Warn("The database does not use indexes. If the full table scan usage exceeds 10%, You are advised to check the slow SQL")
	}
	//统计数据库join语句发生全表扫描占比率Select_full_join *100 /Queries
	selectfullJoinScanUsageRate,err := aa.Streamm.Percentage(cc1["Select_full_join"],cc1["Queries"])
	if selectfullJoinScanUsageRate > 10 && err == nil {
		aa.Loggs.Warn("The database uses the JOIN statement and the non-driver table does not use the index. The full table scan usage is greater than 10%. You are advised to check for slow SQL")
	}

}
func (baselineCheck *DatabaseBaselineCheckStruct) DatabaseTableIndexCheck(aa *PublicClass.ConfigInfo) {
	ignoreTableSchema := "'mysql','information_schema','performance_schema','sys'"
	strSql := fmt.Sprintf("select c.TABLE_SCHEMA,c.TABLE_NAME,t.TABLE_ROWS,t.AVG_ROW_LENGTH,t.DATA_LENGTH,t.MAX_DATA_LENGTH,t.INDEX_LENGTH,t.DATA_FREE,t.AUTO_INCREMENT,t.CREATE_TIME,t.UPDATE_TIME,c.column_name,c.DATA_TYPE,c.COLUMN_TYPE,c.EXTRA from information_schema.tables t join information_schema.columns c on t.table_name=c.table_name where c.EXTRA='auto_increment' and t.table_schema not in (%s);",ignoreTableSchema)
	cc := aa.DatabaseExecInterf.DBQueryDateJson(aa,strSql)
	nowDateTime := aa.DatabaseExecInterf.DBQueryDateString(aa,"select now() as datetime")
	for _,v := range cc{
		//检查是否存在自增主键溢出风险。统计数据库自增id列快要溢出的表
		if v["EXTRA"] == "auto_increment" && v["DATA_TYPE"].(string) == "int" {
			if strings.Contains(v["COLUMN_TYPE"].(string),"unsigned"){
				if unsignedIntUsageRate,err := aa.Streamm.Percentage(v["AUTO_INCREMENT"],MunsignedInt);err ==nil && unsignedIntUsageRate >=1{
					aa.Loggs.Warn(fmt.Sprintf("The self-value-added usage of tables in the database exceeds 85%%, causing data type overflow risks. The details are as follows: Database: \"%v\", table name: \"%v\", increment column name: \"%v\", increment column data type: \"%v\", current increment column: \"%v\"",v["TABLE_SCHEMA"],v["TABLE_NAME"],v["column_name"],v["COLUMN_TYPE"],v["AUTO_INCREMENT"]))
				}
			}else{
				if intUsageRate,err := aa.Streamm.Percentage(v["AUTO_INCREMENT"],Mint);err ==nil && intUsageRate >=1{
					aa.Loggs.Warn(fmt.Sprintf("The self-value-added usage of tables in the database exceeds 85%%, causing data type overflow risks. The details are as follows: Database: \"%v\", table name: \"%v\", increment column name: \"%v\", increment column data type: \"%v\", current increment column: \"%v\"",v["TABLE_SCHEMA"],v["TABLE_NAME"],v["column_name"],v["COLUMN_TYPE"],v["AUTO_INCREMENT"]))
				}
			}
		}
		//单表行数大于500w，且平均行长大于10KB。
		if v["TABLE_ROWS"].(int64) >5000000 && v["AVG_ROW_LENGTH"].(int64)/1024 >10{
			aa.Loggs.Warn("aaa")
		}
		// 单表大于6G，并且碎片率大于30%。
		dataLengthTotal := v["DATA_LENGTH"].(int64)+v["INDEX_LENGTH"].(int64)   //表空间
		if diskFragmentationRate,err := aa.Streamm.Percentage(v["DATA_FREE"],dataLengthTotal); diskFragmentationRate >1 && err == nil && dataLengthTotal/1024/1024/1024 >1 {
			fmt.Println("bbb")
		}
		//单表行数大于1000W，且表空间大于30G
		if dataLengthTotal/1024/1024/1024 >30 && v["TABLE_ROWS"].(int64) >10000000{
			fmt.Println("ccc")
		}
		//检查一个星期内未更新的表
		var tableUpdateTime string
		if v["v[UPDATE_TIME"] == nil {
			tableUpdateTime = v["CREATE_TIME"].(string)
		}else {
			tableUpdateTime = v["UPDATE_TIME"].(string)
		}
		arrDay,err := aa.Streamm.GetTimeDayArr(tableUpdateTime,nowDateTime)
		if err != nil{
			fmt.Println(err)
		}
		if arrDay >7 {
			fmt.Println("有大于7天的表")
		}

	}



}

