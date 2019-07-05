package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"log"
)
//MySql数据库连接
var mysqlConn	 	*sql.DB

type DB struct {
	tables		string
	selects 	string
	wheres		string
	limits		string
	joins		string
	orders		string
	dataMapping	[]interface{}		//Data mapping , 序列化传参
}

func connect() *DB {
	if mysqlConn == nil {
		mysqlConn, _ = sql.Open("mysql", "")
		mysqlConn.SetMaxOpenConns(100)		//最大连接数
		mysqlConn.SetMaxIdleConns(50)		//空闲连接数
	}
	//MysqlConn.Ping()
	return &DB{
		selects: "*",
	}
}

//table
func Table(table string) (db *DB) {
	db = connect()
	db.tables = table
	return db
}

//条件
func (db *DB) Where(where string, args ...interface{}) *DB {
	if where != "" {
		db.wheres = " WHERE " + where
		db.dataMapping = append(db.dataMapping, args)
	}
	return db
}

//查询字段
func (db *DB) Select(selects string) *DB {
	db.selects = selects
	return db
}

//limit
func (db *DB) Limit(limit string) *DB {
	if limit != "" {
		db.limits = " LIMIT " + limit
	}
	return db
}

func (db *DB) Join(join string) *DB {
	db.joins = join
	return db
}

func (db *DB) Order(order string) *DB {
	if order != "" {
		db.orders = " ORDER BY "+order;
	}
	return db
}

//查询
func (db *DB) Get() ([]map[string]string, error) {
	select_sql := "SELECT "+db.selects+" FROM "+db.tables
	if db.joins != "" {
		select_sql += " "+db.joins
	}
	if db.wheres != "" {
		select_sql += db.wheres
	}
	if db.orders != "" {
		select_sql += db.orders
	}
	if db.limits != "" {
		select_sql += db.limits
	}
	var data []map[string]string
	//查询多条
	select_rows,err := mysqlConn.Query(select_sql, db.dataMapping ...)
	if err != nil {
		panic(err.Error())
	}
	for select_rows.Next() {
		columns, _ := select_rows.Columns()

		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))

		for i := range values {
			scanArgs[i] = &values[i]
		}

		//将数据保存到 record 字典
		err = select_rows.Scan(scanArgs...)
		record := make(map[string]string)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		data = append(data, record)
	}
	select_rows.Close()
	return data,err
}

//查询单条
func (DB *DB) First(selects string) (data map[string]string) {
	select_sql := "SELECT "+selects+" FROM "+DB.tables
	if DB.wheres != "" {
		select_sql += DB.wheres
	}

	columns := strings.Split(selects,",")
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for i := range values {
		scanArgs[i] = &values[i]
	}
	select_err := mysqlConn.QueryRow(select_sql).Scan(scanArgs...)
	//将数据保存到 record 字典
	record := make(map[string]string)
	for i, col := range values {
		if col != nil {
			record[columns[i]] = string(col.([]byte))
		}
	}

	if select_err != nil { //如果没有查询到任何数据就进入if中err：no rows in result set
		log.Println(select_err)
		return record
	}

	//log.Println(data)
	return record
}

//删除
func delete() {}

//更新
func update() {}

//添加单条
func (DB *DB) Insert(data map[string]string) {

}

//添加多条
func (DB *DB) InsertAll() {

}

//count
func (DB *DB) Count() (int,error) {
	sql := "SELECT count(*) FROM `"+DB.tables+"`"
	if DB.wheres != "" {
		sql += DB.wheres
	}
	var count int
	err := mysqlConn.QueryRow(sql).Scan(&count)
	if err != nil {
		panic(err.Error())
		return 0,err
	}
	//MysqlConn.Close()	 | 不需要关闭
	return count,err
}

//原始sql
func (DB *DB) _DoExec()  {

}


/**
 * 启动事务
 * @return void
 */
func startTrans() {}

/**
* 用于非自动提交状态下面的查询提交
* @return boolen
 */
func commit() {}

/**
 * 事务回滚
 * @return boolen
 */
func rollback() {}