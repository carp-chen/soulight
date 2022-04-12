package model

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/didi/gendry/manager"
	"github.com/didi/gendry/scanner"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron/v3"
)

var Db *sql.DB
var Cron *cron.Cron
var err error

func init() {
	Cron = cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)), cron.WithLogger(
		cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))
	Cron.Start()
}

func Database(DbHost string, DbPort string, DbUser string, DbPassWord string, DbName string) {
	dbport, _ := strconv.Atoi(DbPort)
	Db, err = manager.New(DbName, DbUser, DbPassWord, DbHost).Set(
		manager.SetCharset("utf8mb4"),
		manager.SetParseTime(true),
		manager.SetLoc(url.QueryEscape("Asia/Shanghai")),
		manager.SetAllowCleartextPasswords(true),
		manager.SetInterpolateParams(true),
		manager.SetTimeout(1*time.Second),
		manager.SetReadTimeout(1*time.Second)).Port(dbport).Open(true)
	scanner.SetTagName("json")
	if err != nil {
		panic(err)
	}
	err = Db.Ping()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("数据库连接成功！")
	}
}
