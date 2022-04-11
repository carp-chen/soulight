package middleware

import (
	"soulight/model"
	"soulight/response"
	"soulight/utils"
	"soulight/utils/errmsg"

	"github.com/gin-gonic/gin"
)

func Identify() gin.HandlerFunc {
	return func(c *gin.Context) {
		cla, _ := c.Get("claims")
		claims, _ := cla.(*utils.Claims)
		if claims.Identity == "user" {
			u, _ := model.GetOneUser(model.Db, map[string]interface{}{"id": claims.Id})
			if u == nil {
				response.SendResponse(c, errmsg.ERROR_USER_NOT_EXIST)
				c.Abort()
				return
			}
			c.Set("indentity", u)
		} else if claims.Identity == "adviser" {
			a, _ := model.GetOneAdviser(model.Db, map[string]interface{}{"id": claims.Id})
			if a == nil {
				response.SendResponse(c, errmsg.ERROR_USER_NOT_EXIST)
				c.Abort()
				return
			}
			c.Set("indentity", a)
		}
	}

}
