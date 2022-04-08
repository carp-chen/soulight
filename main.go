package main

import (
	"soulight/conf"
	"soulight/routes"
)

/**
1. 封装返回
2. 中间件获取user或者adviser
**/
func main() {
	conf.Init()
	r := routes.NewRouter()
	_ = r.Run(conf.HttpPort)
}
