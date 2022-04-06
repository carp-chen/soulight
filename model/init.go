package model

import (
	"database/sql"
	"fmt"
	"net/url"
	"strconv"

	"github.com/didi/gendry/manager"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func Database(DbHost string, DbPort string, DbUser string, DbPassWord string, DbName string) {
	dbport, _ := strconv.Atoi(DbPort)
	db, err = manager.New(DbName, DbUser, DbPassWord, DbHost).Set(
		manager.SetCharset("utf8mb4"),
		manager.SetParseTime(true),
		manager.SetLoc(url.QueryEscape("Asia/Shanghai"))).Port(dbport).Open(true)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("数据库连接成功！")
	}
}
