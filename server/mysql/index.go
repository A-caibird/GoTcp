package mysqlDB

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

const dataSourceName = "root:775028@tcp(127.0.0.1:3306)/GoTcp?charset=utf8mb4&parseTime=True&loc=Local"

func InitDB() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", dataSourceName)
	db.SetMaxIdleConns(10)
	// 最大连接数
	db.SetConnMaxLifetime(100)
	if err != nil {
		panic(err)
	}
	return
}
func DBError(err error) error {
	if err != nil {
		switch {
		case errors.Is(err, errors.New("数据库连接异常")):
			fmt.Println("数据库连接异常!")
		case errors.Is(err, errors.New("数据库准备 SQL 语句异常")):
			fmt.Println("数据库准备 SQL 语句异常!")
		case errors.Is(err, errors.New("数据库执行 SQL 语句异常")):
			fmt.Println("数据库执行 SQL 语句异常!")
		}
		return err
	}
	return nil
}
