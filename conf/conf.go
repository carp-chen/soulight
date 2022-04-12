package conf

import (
	"fmt"
	"soulight/model"

	"gopkg.in/ini.v1"
)

var (
	AppMode    string
	HttpPort   string
	JwtKey     string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string
	RedisHost  string
	RedisPort  string
)

func Init() {
	file, err := ini.Load("./conf/conf.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
		panic(err)
	}
	LoadServer(file)
	LoadMysqlData(file)
	LoadRedis(file)
	model.Init(DbHost, DbPort, DbUser, DbPassWord, DbName, RedisHost, RedisPort)
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").String()
	HttpPort = file.Section("server").Key("HttpPort").String()
	JwtKey = file.Section("server").Key("JwtKey").String()
}

func LoadMysqlData(file *ini.File) {
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
}

func LoadRedis(file *ini.File) {
	RedisHost = file.Section("redis").Key("RedisHost").String()
	RedisPort = file.Section("redis").Key("RedisPort").String()
}
