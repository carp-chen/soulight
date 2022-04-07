package api

import (
	"fmt"
	"net/http"
	"soulight/model"
	"soulight/serialization"
	"soulight/utils/bcrypt"
	"soulight/utils/errmsg"
	"soulight/utils/jwt"
	"soulight/utils/validator"
	"strconv"

	"github.com/fatih/structs"

	"github.com/gin-gonic/gin"
)

func Hello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello!"})
}

func UserRegister(c *gin.Context) {
	var user model.User
	var msg string
	var validCode int
	var code int
	//1.绑定参数
	if err := c.ShouldBind(&user); err != nil {
		code = 400
		c.JSON(
			http.StatusOK, serialization.NewResponse(code),
		)
		return
	}
	//2.校验字段
	msg, validCode = validator.Validate(&user)
	if validCode != errmsg.SUCCSE {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  validCode,
				"message": msg,
			},
		)
		return
	}
	//3.用户名如已存在,即为登陆
	u, _ := model.GetOneUser(model.Db, map[string]interface{}{"username": user.Username})
	if u != nil {
		//如果密码验证正确,则返回token
		if bcrypt.CheckPassword(user.Password, u.Password) {
			//分发token
			token, err := jwt.GenerateToken(uint(u.ID), u.Username)
			if err != nil {
				fmt.Println(err)
				code = 1009
				c.JSON(
					http.StatusOK, serialization.NewResponse(code),
				)
				return
			}
			code = 200
			c.JSON(
				http.StatusOK, serialization.NewResponseWithToken(code, u, token),
			)

		} else {
			code = 1002
			c.JSON(
				http.StatusOK, serialization.NewResponse(code),
			)
		}
		return
	}
	//4.加密密码并写入数据库
	passwordDigest, _ := bcrypt.SetPassword(user.Password)
	//	if _, err := model.Db.Exec("insert into user(username,password) values(?,?)", user.Username, passwordDigest)
	if _, err := model.InsertUser(model.Db, []map[string]interface{}{{"username": user.Username, "password": passwordDigest}}); err != nil {
		code = 3000
		c.JSON(
			http.StatusOK, serialization.NewResponse(code),
		)
	} else {
		user, _ := model.GetOneUser(model.Db, map[string]interface{}{"username": user.Username})
		code = 200
		c.JSON(
			http.StatusOK, serialization.NewResponseWithData(code, user),
		)
	}
}

func EditUser(c *gin.Context) {
	var user model.User
	var code int
	//1.参数绑定
	id, _ := strconv.Atoi(c.Param("id"))
	if err := c.ShouldBind(&user); err != nil {
		code = 500
		c.JSON(
			http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			},
		)
		return
	}
	//2.判断要修改的用户名是否存在
	u, _ := model.GetOneUser(model.Db, map[string]interface{}{"username": user.Username})
	if u != nil && u.ID != id {
		code = 1001
		c.JSON(
			http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			},
		)
		return
	}
	//3.更新数据库
	user_map := structs.Map(&user)
	if result, _ := model.UpdateUser(model.Db, map[string]interface{}{"id": id}, user_map); result != 0 {
		code = 500
		c.JSON(
			http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			},
		)
	} else {
		code = 200
		c.JSON(
			http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			},
		)
	}

}
