package api

import (
	"fmt"
	"net/http"
	"soulight/model"
	"soulight/serialization"
	"soulight/utils"
	"soulight/utils/errmsg"
	"strconv"
	"time"

	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/gin-gonic/gin"
)

func Hello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello!"})
}

//用户注册(登录)接口
func UserRegister(c *gin.Context) {
	var user model.User
	var msg string
	var validCode int
	var code int
	//1.绑定参数
	if err := c.ShouldBind(&user); err != nil {
		code = errmsg.INVALID_PARAMS
		c.JSON(
			http.StatusOK, serialization.NewResponse(code),
		)
		return
	}
	//2.校验字段
	msg, validCode = utils.Validate(&user)
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
		if utils.CheckPassword(user.Password, u.Password) {
			//分发token
			token, err := utils.GenerateToken(u.ID, u.Username)
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
				http.StatusOK, serialization.NewResponseWithToken(code, u, token),
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
	passwordDigest, _ := utils.SetPassword(user.Password)
	//	if _, err := model.Db.Exec("insert into user(username,password) values(?,?)", user.Username, passwordDigest)
	if _, err := model.InsertUser(model.Db, []map[string]interface{}{{"username": user.Username, "password": passwordDigest}}); err != nil {
		code = errmsg.ERROR_DATABASE
		c.JSON(
			http.StatusOK, serialization.NewResponse(code),
		)
	} else {
		user, _ := model.GetOneUser(model.Db, map[string]interface{}{"username": user.Username})
		code = errmsg.SUCCSE
		c.JSON(
			http.StatusOK, serialization.NewResponseWithData(code, user),
		)
	}
}

//用户修改信息接口
func UserEdit(c *gin.Context) {
	var edit_user model.EditUser
	var code int
	//1.参数绑定
	id := c.GetInt("id")
	if err := c.ShouldBind(&edit_user); err != nil {
		code = errmsg.INVALID_PARAMS
		c.JSON(
			http.StatusOK, serialization.NewResponse(code),
		)
		return
	}
	//2.判断要修改的用户名是否存在
	u, _ := model.GetOneUser(model.Db, map[string]interface{}{"username": edit_user.Username})
	if u != nil && u.ID != id {
		code = errmsg.ERROR_USERNAME_USED
		c.JSON(
			http.StatusOK, serialization.NewResponse(code),
		)
		return
	}
	//3.更新数据库
	update_map := map[string]interface{}{
		"username": edit_user.Username,
		"img":      edit_user.Img,
		"birth":    time.Unix(edit_user.Birth, 0),
		"gender":   edit_user.Gender,
		"bio":      edit_user.Bio,
		"about":    edit_user.About,
	}
	if _, err := model.UpdateUser(model.Db, map[string]interface{}{"id": id}, update_map); err != nil {
		fmt.Println(err)
		code = errmsg.ERROR_DATABASE
		c.JSON(
			http.StatusOK, serialization.NewResponse(code),
		)
	} else {
		user, _ := model.GetOneUser(model.Db, map[string]interface{}{"id": id})
		code = errmsg.SUCCSE
		c.JSON(
			http.StatusOK, serialization.NewResponseWithData(code, user),
		)
	}

}

//顾问列表接口
func AdviserList(c *gin.Context) {
	var code int
	//1.参数绑定
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	//2.查询数据库
	offset := (pageNum - 1) * pageSize
	where := map[string]interface{}{"_limit": []uint{uint(offset), uint(pageSize)}}
	columns := []string{"adviser_name", "img", "bio"}
	cond, vals, err := builder.BuildSelect("adviser", where, columns)
	if nil != err {
		code = errmsg.ERROR
		c.JSON(
			http.StatusOK, serialization.NewResponse(code),
		)
		return
	}
	row, err := model.Db.Query(cond, vals...)
	if nil != err || nil == row {
		fmt.Println(err)
		code = errmsg.ERROR_DATABASE
		c.JSON(
			http.StatusOK, serialization.NewResponse(code),
		)
		return
	}
	defer row.Close()
	var res []*model.AdviserList
	if err = scanner.Scan(row, &res); err != nil {
		code = errmsg.ERROR
		c.JSON(
			http.StatusOK, serialization.NewResponse(code),
		)
	}
	code = errmsg.SUCCSE
	c.JSON(
		http.StatusOK, serialization.NewResponseWithData(code, res),
	)
}

//顾问主页接口
func AdviserInfoForUser(c *gin.Context) {
	var adviserinfo model.AdviserInfoForUser
	var code int
	//1.参数绑定
	adviser_id, _ := strconv.Atoi(c.Query("adviser_id"))
	//2.查询adviser表
	where := map[string]interface{}{"id": adviser_id}
	columns := []string{"adviser_name", "img", "bio", "about"}
	cond, vals, err := builder.BuildSelect("adviser", where, columns)
	if nil != err {
		code = errmsg.ERROR
		c.JSON(
			http.StatusOK, serialization.NewResponse(code),
		)
		return
	}
	row, err := model.Db.Query(cond, vals...)
	if nil != err || nil == row {
		fmt.Println(err)
		code = errmsg.ERROR_DATABASE
		c.JSON(
			http.StatusOK, serialization.NewResponse(code),
		)
		return
	}
	defer row.Close()
	if err = scanner.Scan(row, &adviserinfo); err != nil {
		code = errmsg.ERROR
		c.JSON(
			http.StatusOK, serialization.NewResponse(code),
		)
		return
	}
	//3.查询service表
	where = map[string]interface{}{"adviser_id": adviser_id}
	var res []*model.Service
	if res, err = model.GetMultiService(model.Db, where); err != nil {
		fmt.Println(err)
		code = errmsg.ERROR_DATABASE
		c.JSON(
			http.StatusOK, serialization.NewResponse(code),
		)
		return
	}
	adviserinfo.Services = res
	code = errmsg.SUCCSE
	c.JSON(
		http.StatusOK, serialization.NewResponseWithData(code, adviserinfo),
	)

}
