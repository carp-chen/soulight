package model

import (
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/didi/gendry/manager"
	"github.com/didi/gendry/scanner"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

// var Cron *cron.Cron
var Pool *redis.Pool
var err error

func Init(DbHost string, Dbport string, Dbuser string, Dbpass string, Dbname string, RedisHost string, RedisPort string) {
	// InitCron()
	InitMysql(DbHost, Dbport, Dbuser, Dbpass, Dbname)
	InitRedis(RedisHost, RedisPort)
}

// func InitCron() {
// 	Cron = cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)), cron.WithLogger(
// 		cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))
// 	Cron.Start()
// }

func InitRedis(Host string, Port string) {
	addresses := Host + ":" + Port
	Pool = &redis.Pool{ //实例化一个连接池
		MaxIdle: 16, //最初的连接数量
		// MaxActive:1000000,    //最大连接数量
		MaxActive:   0,   //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		IdleTimeout: 300, //连接关闭时间 300秒 （300秒不使用自动关闭）
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			return redis.Dial("tcp", addresses)
		},
	}
	c := Pool.Get() //从连接池，取一个链接
	defer c.Close()
	_, err := c.Do("ping")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("redis数据库连接成功！")
	}
}

func InitMysql(DbHost string, DbPort string, DbUser string, DbPassWord string, DbName string) {
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
		fmt.Println("mysql数据库连接成功！")
	}
}
