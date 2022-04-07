package routes

import (
	"soulight/api"

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
		v1.POST("user/edit", api.EditUser)
	}
	return r
}
