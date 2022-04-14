package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/manager"
	"github.com/didi/gendry/scanner"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron/v3"
)

var Db *sql.DB
var err error
var Cron *cron.Cron

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username" validate:"required,min=4,max=20" label:"用户名"`
	Password  string    `json:"password" validate:"required,min=6,max=30" label:"密码"`
	Img       string    `json:"img"`
	Birth     time.Time `json:"birth"`
	Gender    int8      `json:"gender"`
	Bio       string    `json:"bio"`
	About     string    `json:"about"`
	Coins     int       `json:"coins"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"update_at"`
}

//GetOne gets one record from table user by condition "where"
func GetOneUser(db *sql.DB, where map[string]interface{}) (*User, error) {
	if nil == db {
		return nil, errors.New("sql.DB object couldn't be nil")
	}
	cond, vals, err := builder.BuildSelect("user", where, nil)
	if nil != err {
		return nil, err
	}
	row, err := db.Query(cond, vals...)
	if nil != err || nil == row {
		return nil, err
	}
	defer row.Close()
	var res *User
	err = scanner.Scan(row, &res)
	return res, err
}

type Order struct {
	OrderID      string    `json:"order_id"`
	UserID       int       `json:"user_id"`
	AdviserID    int       `json:"adviser_id"`
	Situation    string    `json:"situation"`
	Question     string    `json:"question"`
	Reply        string    `json:"reply"`
	Cost         int       `json:"cost"`
	Status       int8      `json:"status"`
	ServiceType  int8      `json:"service_type"`
	OrderTime    time.Time `json:"order_time"`
	UrgentTime   time.Time `json:"urgent_time"`
	DeliveryTime time.Time `json:"delivery_time"`
	Rate         int       `json:"rate"`
}

type Transaction struct {
	ID          int       `json:"id"`
	Action      int8      `json:"action"`
	OrderID     string    `json:"order_id"`
	ServiceType int8      `json:"service_type"`
	Coins       int       `json:"coins"`
	Credits     int       `json:"credits"`
	CreateTime  time.Time `json:"create_time"`
}

//GetOne gets one record from table transaction by condition "where"
func GetOneTransaction(db *sql.DB, where map[string]interface{}) (*Transaction, error) {
	if nil == db {
		return nil, errors.New("sql.DB object couldn't be nil")
	}
	cond, vals, err := builder.BuildSelect("transaction_user", where, nil)
	if nil != err {
		return nil, err
	}
	row, err := db.Query(cond, vals...)
	if nil != err || nil == row {
		return nil, err
	}
	defer row.Close()
	var res *Transaction
	err = scanner.Scan(row, &res)
	return res, err
}

func init() {
	Db, err = manager.New("soulight", "root", "123456", "127.0.0.1").Set(
		manager.SetCharset("utf8mb4"),
		manager.SetParseTime(true),
		manager.SetLoc(url.QueryEscape("Asia/Shanghai")),
		manager.SetAllowCleartextPasswords(true),
		manager.SetInterpolateParams(true),
		manager.SetTimeout(1*time.Second),
		manager.SetReadTimeout(1*time.Second)).Port(3306).Open(true)
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
	Cron = cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)), cron.WithLogger(
		cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))
	Cron.Start()
}

func UpdateExpOrders() {
	Cron.AddFunc("@every 5min", func() {
		//1.查找所有过期的订单
		rows, err := Db.Query("select * from orders where status in(0,3) and CURRENT_TIMESTAMP>=ADDTIME(order_time,'24:00:00')")
		if nil != err || nil == rows {
			fmt.Println(err)
			return
		}
		defer rows.Close()
		var res []*Order
		if err := scanner.Scan(rows, &res); err != nil {
			fmt.Println(err)
			return
		}
		//2.遍历过期订单数组，更新订单状态,并且退还用户的金币，并记录流水
		for _, order := range res {
			user, _ := GetOneUser(Db, map[string]interface{}{"id": order.UserID})
			conn, _ := Db.Begin()
			if _, err := conn.Exec("update orders set status=2 where order_id=?", order.OrderID); err != nil {
				fmt.Println(err)
				conn.Rollback()
				return
			}
			if _, err := conn.Exec("update user set coins=coins+? where id=?", order.Cost, order.UserID); err != nil {
				fmt.Println(err)
				conn.Rollback()
				return
			}
			if _, err := conn.Exec(`insert into transaction_user(action,id,order_id,service_type,coins,credits)
		        values(?,?,?,?,?,?)`, 2, order.UserID, order.OrderID, order.ServiceType, user.Coins+order.Cost, order.Cost); err != nil {
				fmt.Println(err)
				conn.Rollback()
				return
			}
			conn.Commit()
		}

	})
}

func UpdateUrgentOrders() {
	Cron.AddFunc("@every 1min", func() {
		//1.查找所有过期的订单
		rows, err := Db.Query("select * from orders where status=3 and CURRENT_TIMESTAMP>=ADDTIME(urgent_time,'01:00:00')")
		if nil != err || nil == rows {
			fmt.Println(err)
			return
		}
		defer rows.Close()
		var res []*Order
		if err := scanner.Scan(rows, &res); err != nil {
			fmt.Println(err)
			return
		}
		//2.遍历加急过期订单数组，更新订单状态,并且退还用户加急的金币，并记录流水
		for _, order := range res {
			user, _ := GetOneUser(Db, map[string]interface{}{"id": order.UserID})
			transaction, _ := GetOneTransaction(Db, map[string]interface{}{"order_id": order.OrderID, "action": 3})
			conn, _ := Db.Begin()
			if _, err := conn.Exec("update orders set status=0,cost=cost-? where order_id=?", -transaction.Credits, order.OrderID); err != nil {
				fmt.Println(err)
				conn.Rollback()
				return
			}
			if _, err := conn.Exec("update user set coins=coins+? where id=?", -transaction.Credits, order.UserID); err != nil {
				fmt.Println(err)
				conn.Rollback()
				return
			}
			if _, err := conn.Exec(`insert into transaction_user(action,id,order_id,service_type,coins,credits)
			    values(?,?,?,?,?,?)`, 4, user.ID, order.OrderID, order.ServiceType, user.Coins-transaction.Credits, -transaction.Credits); err != nil {
				fmt.Println(err)
				conn.Rollback()
				return
			}
			conn.Commit()
		}
	})
}

func main() {
	UpdateExpOrders()
	UpdateUrgentOrders()
	defer Cron.Stop()
	select {}
}
