package routes

import (
	"soulight/api"
	"soulight/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	//r.Use(middleware.NewLogger(),middleware.Cors())
	r.GET("/", api.Hello)
	v1 := r.Group("api/v1")
	{
		// 用户及顾问注册(登陆)
		v1.POST("user/register", api.UserRegister)
		v1.POST("adviser/register", api.AdviserRegister)
		authed := v1.Group("/")                                  //需要登陆保护
		authed.Use(middleware.JwtToken(), middleware.Identify()) //jwt验证，身份验证
		{
			//用户模块
			authed.POST("user/edit", api.UserEdit)
			authed.GET("user/advisers", api.AdviserList)
			authed.GET("user/adviser", api.AdviserInfoForUser)
			//顾问模块
			authed.POST("adviser/edit", api.AdviserEdit)
			authed.POST("adviser/status", api.AdviserStatus)
			authed.GET("adviser/info", api.AdviserInfo)
			authed.POST("adviser/service", api.AdviserService)
			//订单模块
			authed.POST("order/create", api.OrderCreate)
			authed.GET("order/list", api.OrderList)
			authed.GET("order/info", api.OrderInfo)
			authed.POST("order/reply", api.OrderReply)
			authed.POST("order/urgent", api.OrderUrgent)
			//评论模块
			authed.POST("order/review", api.OrderReview)
		}
	}
	return r
}
