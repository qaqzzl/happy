package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)
//MySql数据库连接
var GormDB	 	*gorm.DB

func connect() *gorm.DB {
	if GormDB == nil {
		GormDB, _ = gorm.Open("mysql", "blog:mysqlblog336699MM@tcp(qaqzz-com.mysql.rds.aliyuncs.com)/my_blog")
		GormDB.DB().SetMaxOpenConns(100)		//最大连接数
		GormDB.DB().SetMaxIdleConns(50)			//空闲连接数
		defer GormDB.Close()
	}
	return GormDB
}

//table
func Table(table string) (db *gorm.DB) {
	db = connect()
	return db
}