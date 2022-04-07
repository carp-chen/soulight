package middleware

import (
	"net/http"
	"soulight/serialization"
	"soulight/utils"
	"soulight/utils/errmsg"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// JwtToken jwt中间件
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		tokenHeader := c.Request.Header.Get("Authorization")
		//token不存在
		if tokenHeader == "" {
			code = errmsg.ERROR_TOKEN_NOT_EXIST
			c.JSON(http.StatusOK, serialization.NewResponse(code))
			c.Abort()
			return

		}
		//token格式不正确
		checkToken := strings.Split(tokenHeader, " ")
		if len(checkToken) != 2 || checkToken[0] != "Bearer" {
			code = errmsg.ERROR_TOKEN_TYPE_WRONG
			c.JSON(http.StatusOK, serialization.NewResponse(code))
			c.Abort()
			return
		}
		// 解析token
		claims, err := utils.ParseToken(checkToken[1])
		if err != nil {
			if err == utils.TokenExpired {
				code = errmsg.ERROR_TOKEN_TIMEOUT //token过期
			} else {
				code = errmsg.ERROR_TOKEN_WRONG //token不正确
			}
			c.JSON(http.StatusOK, serialization.NewResponse(code))
			c.Abort()
			return
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = errmsg.ERROR_TOKEN_TIMEOUT
			c.JSON(http.StatusOK, serialization.NewResponse(code))
			c.Abort()
			return
		}

		c.Set("id", claims.Id)
		c.Next()
	}
}
