package PublicClass

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)


type DatabaseOperation interface {
	//连接数据库，返回dbConnSocket信息
	dbConnSocket()
	//执行sql语句，返回列及原始数据类型（包含列及值）
	DBQueryDateTmp(sqlStr string) *sql.Rows
	//执行sql语句，返回json数据类型（包含列及值）
	DBQueryDateJson(sqlStr string) []map[string]interface{}
	//执行sql语句，返回map数据类型（包含列及值）
	DBQueryDateMap(sqlStr string) map[string]string
	//执行sql语句，返回value值
	DBQueryDateList(sqlStr string) []string
	//执行sql语句，返回value的string值
	DBQueryDateString(sqlStr string) string
}

type DatabaseExecStruct struct {
	DirverName string
	ConnInfo   string
	DBconnIdleTime time.Duration
	MaxIdleConns int
	dbConnSocketinfo *sql.DB
}

func (dbStruct *DatabaseExecStruct) dbConnSocket(){
	Loggs.Debug("Initializes the database connection object")
	db ,err := sql.Open(dbStruct.DirverName,dbStruct.ConnInfo)
	if err != nil{
		errStr := fmt.Sprintf("unknown driver %q (forgotten import?)", dbStruct.DirverName)
		fmt.Println(errStr)
		Loggs.Error(errStr)
		os.Exit(1)
	}
	Loggs.Debug("Send a ping packet to check the database running status")
	if err := db.Ping(); err != nil {
		errStr := "Failed to open a database connection and create a session connection. pleace Check the database status or network status"
        fmt.Println(errStr)
		Loggs.Error(errStr)
		os.Exit(1)
	}
	db.SetConnMaxIdleTime(dbStruct.DBconnIdleTime)
	db.SetMaxIdleConns(dbStruct.MaxIdleConns)
	dbStruct.dbConnSocketinfo = db
}

func (dbStruct *DatabaseExecStruct) DBQueryDateTmp(sqlStr string) *sql.Rows{
	dbStruct.dbConnSocket()
	dbconn := dbStruct.dbConnSocketinfo
	Loggs.Debug(fmt.Sprintf("Prepare initialize the SQL statement:%s",sqlStr))
	stmt,err := dbconn.Prepare(sqlStr)
	if err != nil {
		Loggs.Error(fmt.Sprintf("Prepare SQL file ,This is a bad connection. SQL info: %s",sqlStr))
	}
	Loggs.Debug(fmt.Sprintf("Execute SQL statement queries: %s",sqlStr))
	rows,err := stmt.Query()
	if err != nil{
		Loggs.Error(fmt.Sprintf("Execute SQL file ,This is a bad connection. SQL info: %s",sqlStr))
	}
	return rows
}

// 查询数据库，返回数据库接口切片，或返回json（包含列名）
func (dbStruct *DatabaseExecStruct) DBQueryDateJson(sqlStr string) []map[string]interface{}{
	rows := dbStruct.DBQueryDateTmp(sqlStr)
	// 获取列名
	columns,err := rows.Columns()
	if err != nil {
		Loggs.Error("Failed to get database column name. The failure information is as follows:",err)
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

//查询数据库，结果集为两列，切第一列的数据是第二列的key，对列数据生成key-value并返回（不包含列名）
func (dbStruct *DatabaseExecStruct) DBQueryDateMap(sqlStr string) map[string]string{
	rows := dbStruct.DBQueryDateTmp(sqlStr)
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

//查询数据库，返回多行多列数据（不包含列名）
func (dbStruct *DatabaseExecStruct) DBQueryDateList(sqlStr string) []string{
	var a []string
	return a
}

//查询数据库，返回单行单列数据（列名）
func (dbStruct *DatabaseExecStruct) DBQueryDateString(sqlStr string) string{
	rows := dbStruct.DBQueryDateTmp(sqlStr)
	var c string
	for rows.Next(){
		rows.Scan(&c)
	}
	return c
}

