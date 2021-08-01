package PublicClass

import (
	"DepthInspection/api/Stream"
	"DepthInspection/api/loggs"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)
type ConfigInfo struct {
	GetConfig          loggs.BaseInfo
	DatabaseExecInterf DatabaseOperation
	Loggs              loggs.LogOutputInterface
	Streamm            *Stream.StreamStruct
}

type DatabaseOperation interface {
	//连接数据库，返回dbConnSocket信息
	dbConnSocket(aa *ConfigInfo)
	//执行sql语句，返回列及原始数据类型（包含列及值）
	DBQueryDateTmp(aa *ConfigInfo,sqlStr string) *sql.Rows
	//执行sql语句，返回json数据类型（包含列及值）
	DBQueryDateJson(aa *ConfigInfo,sqlStr string) []map[string]interface{}
	//执行sql语句，返回map数据类型（包含列及值）
	DBQueryDateMap(aa *ConfigInfo,sqlStr string) map[string]string
	//执行sql语句，返回value值
	DBQueryDateList(aa *ConfigInfo,sqlStr string) []string
	//执行sql语句，返回value的string值
	DBQueryDateString(aa *ConfigInfo,sqlStr string) string
}

type DatabaseExecStruct struct {
	DirverName string
	ConnInfo   string
	DBconnIdleTime time.Duration
	MaxIdleConns int
	dbConnSocketinfo *sql.DB
}

func (dbStruct *DatabaseExecStruct) dbConnSocket(aa *ConfigInfo){
	db ,_ := sql.Open(dbStruct.DirverName,dbStruct.ConnInfo)
	aa.Loggs.Debug("Initializes the database connection object")
	if err := db.Ping(); err != nil {
		aa.Loggs.Error("Failed to open a database connection and create a session connection. error info: ",err)
		os.Exit(1)
	}
	db.SetConnMaxIdleTime(dbStruct.DBconnIdleTime)
	db.SetMaxIdleConns(dbStruct.MaxIdleConns)
	dbStruct.dbConnSocketinfo = db
}

func (dbStruct *DatabaseExecStruct) DBQueryDateTmp(aa *ConfigInfo,sqlStr string) *sql.Rows{
	dbStruct.dbConnSocket(aa)
	dbconn := dbStruct.dbConnSocketinfo
	stmt,err := dbconn.Prepare(sqlStr)
	aa.Loggs.Debug(fmt.Sprintf("Prepare initialize the SQL statement:%s",sqlStr))
	rows,err := stmt.Query()
	aa.Loggs.Debug(fmt.Sprintf("Execute SQL statement queries: %s",sqlStr))
	if err != nil{
		aa.Loggs.Error(fmt.Sprintf("Execute SQL file ,This is a bad connection. SQL info: %s",sqlStr))
	}
	dbconn.Ping()
	return rows
}
func (dbStruct *DatabaseExecStruct) DBQueryDateJson(aa *ConfigInfo,sqlStr string) []map[string]interface{}{
	rows := dbStruct.DBQueryDateTmp(aa ,sqlStr)
	// 获取列名
	columns,err := rows.Columns()
	if err != nil {
		aa.Loggs.Error("Failed to get database column name. The failure information is as follows:",err)
	}
	// 定义一个切片，长度是字段的个数，切片里面的元素类型是sql.RawBytes
	//values := make([]sql.RawBytes,len(columns))
	//定义一个切片，元素类型是interface{}接口
	//scanArgs := make([]interface{},len(values))
	valuePtrs := make([]interface{}, len(columns))
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, len(columns))
	for rows.Next(){
		for i := 0; i < len(columns); i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {v = string(b)} else {v = val}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	//jsonData, err := json.Marshal(tableData)
	//if err != nil {
	//	aa.Loggs.Panic("Failed to convert database query data to JSON data. error info: ",err)
	//}
	return tableData
}

func (dbStruct *DatabaseExecStruct) DBQueryDateMap(aa *ConfigInfo,sqlStr string) map[string]string{
	rows := dbStruct.DBQueryDateTmp(aa,sqlStr)
	var result map[string]string
	result = make(map[string]string)
	var key_slice,val_slice []string
	var v,c string
	for rows.Next(){
		rows.Scan(&v,&c)
		key_slice = append(key_slice,v)
		val_slice = append(val_slice,c)
	}
	for i := range key_slice{
		result[key_slice[i]] = val_slice[i]
	}
	return result
}

func (dbStruct *DatabaseExecStruct) DBQueryDateList(aa *ConfigInfo,sqlStr string) []string{
	var a []string
	return a
}

func (dbStruct *DatabaseExecStruct) DBQueryDateString(aa *ConfigInfo,sqlStr string) string{
	rows := dbStruct.DBQueryDateTmp(aa,sqlStr)
	var c string
	for rows.Next(){
		rows.Scan(&c)
	}
	return c
}

