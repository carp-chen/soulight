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
			user := authed.Group("user/")
			user.POST("edit", api.UserEdit)
			user.GET("advisers", api.AdviserList)
			user.GET("adviser", api.AdviserInfoForUser)
			user.POST("favorite", api.AddFavorite)
			user.GET("favorites", api.GetFavorites)
			//顾问模块
			adviser := authed.Group("adviser/")
			adviser.POST("edit", api.AdviserEdit)
			adviser.POST("status", api.AdviserStatus)
			adviser.GET("info", api.AdviserInfo)
			adviser.POST("service", api.AdviserService)
			//订单模块
			order := authed.Group("order/")
			order.POST("create", api.OrderCreate)
			order.GET("list", api.OrderList)
			order.GET("info", api.OrderInfo)
			order.POST("reply", api.OrderReply)
			order.POST("urgent", api.OrderUrgent)
			//评论模块
			order.POST("review", api.OrderReview)
			order.POST("reward", api.OrderReward)
		}
	}
	return r
}
