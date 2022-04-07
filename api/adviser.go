package api

import (
	"fmt"
	"net/http"
	"soulight/model"
	"soulight/serialization"
	"soulight/utils"
	"soulight/utils/errmsg"

	"github.com/gin-gonic/gin"
)

//顾问注册(登录)接口
func AdviserRegister(c *gin.Context) {
	var adviser model.Adviser
	var msg string
	var validCode int
	var code int
	//1.绑定参数
	if err := c.ShouldBind(&adviser); err != nil {
		code = errmsg.INVALID_PARAMS
		c.JSON(
			http.StatusOK, serialization.NewResponse(code),
		)
		return
	}
	//2.校验字段
	msg, validCode = utils.Validate(&adviser)
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
	ad, _ := model.GetOneAdviser(model.Db, map[string]interface{}{"adviser_name": adviser.AdviserName})
	if ad != nil {
		//如果密码验证正确,则返回token
		if utils.CheckPassword(adviser.Password, ad.Password) {
			//分发token
			token, err := utils.GenerateToken(ad.ID, ad.AdviserName)
			if err != nil {
				fmt.Println(err)
				code = errmsg.ERROR_GENARATE_TOKEN
				c.JSON(
					http.StatusOK, serialization.NewResponse(code),
				)
				return
			}
			code = errmsg.SUCCSE
			c.JSON(
				http.StatusOK, serialization.NewResponseWithToken(code, ad, token),
			)

		} else {
			code = errmsg.ERROR_PASSWORD_WRONG
			c.JSON(
				http.StatusOK, serialization.NewResponse(code),
			)
		}
		return
	}
	//4.加密密码并写入数据库
	passwordDigest, _ := utils.SetPassword(adviser.Password)
	if _, err := model.InsertAdviser(model.Db, []map[string]interface{}{{"adviser_name": adviser.AdviserName, "password": passwordDigest}}); err != nil {
		code = errmsg.ERROR_DATABASE
		c.JSON(
			http.StatusOK, serialization.NewResponse(code),
		)
	} else {
		adviser, _ := model.GetOneAdviser(model.Db, map[string]interface{}{"adviser_name": adviser.AdviserName})
		code = errmsg.SUCCSE
		c.JSON(
			http.StatusOK, serialization.NewResponseWithData(code, adviser),
		)
	}
}

//顾问修改信息接口
func AdviserEdit(c *gin.Context) {
	var edit_adviser model.EditAdviser
	var code int
	//1.参数绑定
	id := c.GetInt("id")
	if err := c.ShouldBind(&edit_adviser); err != nil {
		code = errmsg.INVALID_PARAMS
		c.JSON(
			http.StatusOK, serialization.NewResponse(code),
		)
		return
	}
	//2.判断要修改的用户名是否存在
	ad, _ := model.GetOneAdviser(model.Db, map[string]interface{}{"adviser_name": edit_adviser.AdviserName})
	if ad != nil && ad.ID != id {
		code = errmsg.ERROR_USERNAME_USED
		c.JSON(
			http.StatusOK, serialization.NewResponse(code),
		)
		return
	}
	//3.更新数据库
	update_map := map[string]interface{}{
		"adviser_name": edit_adviser.AdviserName,
		"bio":          edit_adviser.Bio,
		"work_exp":     edit_adviser.WorkExp,
		"about":        edit_adviser.About,
	}
	if _, err := model.UpdateAdviser(model.Db, map[string]interface{}{"id": id}, update_map); err != nil {
		fmt.Println(err)
		code = errmsg.ERROR_DATABASE
		c.JSON(
			http.StatusOK, serialization.NewResponse(code),
		)
	} else {
		adviser, _ := model.GetOneAdviser(model.Db, map[string]interface{}{"id": id})
		code = errmsg.SUCCSE
		c.JSON(
			http.StatusOK, serialization.NewResponseWithData(code, adviser),
		)
	}
}

//顾问主页接口
func AdviserInfo(c *gin.Context) {

}
