package utils

import (
	"database/sql"
	"fmt"
	"okc/config"

	_ "github.com/go-sql-driver/mysql"
)

var (
	Db  *sql.DB
	Err error
)

func init() {
	m := config.Config().MYSQL.(map[interface{}]interface{})
	sdn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", m["DATABASE_ROOT"], m["DATABASE_PASSWORD"], m["DATABASE_HOST"], m["DATABASE_PORT"], m["DATABASE_NAME"])
	Db, Err = sql.Open("mysql", sdn)

	// 设置最大连接数
	Db.SetMaxOpenConns(0)
	// 设置空闲连接池
	Db.SetMaxIdleConns(100)

	if Err != nil {
		panic(Err.Error())
	}
}
