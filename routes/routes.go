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
		// 用户操作
		v1.POST("user/register", api.UserRegister)
		authed := v1.Group("/") //需要登陆保护
		authed.Use(middleware.JwtToken())
		{
			authed.POST("user/edit", api.EditUser)
		}
	}
	return r
}
