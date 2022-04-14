package api

import (
	"net/http"
	"soulight/model"
	"soulight/response"
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
	//1.绑定参数
	if err := c.ShouldBind(&user); err != nil {
		response.SendResponse(c, errmsg.INVALID_PARAMS)
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
			token, err := utils.GenerateToken(u.ID, u.Username, "user")
			if err != nil {
				response.SendResponse(c, errmsg.ERROR_GENARATE_TOKEN)
			} else {
				response.SendResponse(c, errmsg.SUCCSE, u, token)
			}
		} else {
			response.SendResponse(c, errmsg.ERROR_PASSWORD_WRONG)
		}
		return
	}
	//4.加密密码并写入数据库
	passwordDigest, _ := utils.SetPassword(user.Password)
	//	if _, err := model.Db.Exec("insert into user(username,password) values(?,?)", user.Username, passwordDigest)
	if _, err := model.InsertUser(model.Db, []map[string]interface{}{{"username": user.Username, "password": passwordDigest}}); err != nil {
		response.SendResponse(c, errmsg.ERROR_DATABASE)
	} else {
		user, _ := model.GetOneUser(model.Db, map[string]interface{}{"username": user.Username})
		response.SendResponse(c, errmsg.SUCCSE, user)
	}
}

//用户修改信息接口
func UserEdit(c *gin.Context) {
	var edit_user model.EditUser
	//1.参数绑定
	id := c.GetInt("id")
	if err := c.ShouldBind(&edit_user); err != nil {
		response.SendResponse(c, errmsg.INVALID_PARAMS)
		return
	}
	//2.判断要修改的用户名是否存在
	u, _ := model.GetOneUser(model.Db, map[string]interface{}{"username": edit_user.Username})
	if u != nil && u.ID != id {
		response.SendResponse(c, errmsg.ERROR_USERNAME_USED)
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
		response.SendResponse(c, errmsg.ERROR_DATABASE)
	} else {
		user, _ := model.GetOneUser(model.Db, map[string]interface{}{"id": id})
		response.SendResponse(c, errmsg.SUCCSE, user)
	}
}

//顾问列表接口
func AdviserList(c *gin.Context) {
	//1.参数绑定
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	//2.查询数据库
	offset := (pageNum - 1) * pageSize
	where := map[string]interface{}{"_limit": []uint{uint(offset), uint(pageSize)}}
	// columns := []string{"adviser_name", "img", "bio"}
	cond, vals, err := builder.BuildSelect("adviser", where, nil)
	if nil != err {
		response.SendResponse(c, errmsg.ERROR)
		return
	}
	rows, err := model.Db.Query(cond, vals...)
	if nil != err || nil == rows {
		response.SendResponse(c, errmsg.ERROR_DATABASE)
		return
	}
	defer rows.Close()
	var res []*model.Adviser
	if err = scanner.Scan(rows, &res); err != nil {
		response.SendResponse(c, errmsg.ERROR)
		return
	}
	response.SendResponse(c, errmsg.SUCCSE, res)
}

//顾问主页接口
func AdviserInfoForUser(c *gin.Context) {
	var adviserinfo model.AdviserInfoForUser
	//1.参数绑定
	adviser_id, _ := strconv.Atoi(c.Query("adviser_id"))
	//2.查询adviser表
	where := map[string]interface{}{"id": adviser_id}
	columns := []string{"adviser_name", "img", "bio", "about"}
	cond, vals, err := builder.BuildSelect("adviser", where, columns)
	if nil != err {
		response.SendResponse(c, errmsg.ERROR)
		return
	}
	row, err := model.Db.Query(cond, vals...)
	if nil != err || nil == row {
		response.SendResponse(c, errmsg.ERROR_DATABASE)
		return
	}
	defer row.Close()
	if err = scanner.Scan(row, &adviserinfo); err != nil {
		response.SendResponse(c, errmsg.ERROR)
		return
	}
	//3.查询service表
	where = map[string]interface{}{"adviser_id": adviser_id}
	var res []*model.Service
	if res, err = model.GetMultiService(model.Db, where); err != nil {
		response.SendResponse(c, errmsg.ERROR_DATABASE)
		return
	}
	adviserinfo.Services = res
	response.SendResponse(c, errmsg.SUCCSE, adviserinfo)
}

//收藏顾问
func AddFavorite(c *gin.Context) {
	//1.参数绑定
	user_id := c.GetInt("id")
	adviser_id, _ := strconv.Atoi(c.Query("adviser_id"))
	//2.添加到数据库
	if _, err := model.Db.Exec("insert into favorite(user_id,adviser_id) values(?,?)", user_id, adviser_id); err != nil {
		response.SendResponse(c, errmsg.ERROR_DUPLICATE_FAVORITE)
		return
	}
	response.SendResponse(c, errmsg.SUCCSE)
}

//获取收藏顾问列表
func GetFavorites(c *gin.Context) {
	//1.参数绑定
	user_id := c.GetInt("id")
	//2.查询数据库，获取用户收藏的所有顾问
	rows, err := model.Db.Query("select * from adviser where id in(select adviser_id from favorite where user_id=?)", user_id)
	if nil != err || nil == rows {
		response.SendResponse(c, errmsg.ERROR_DATABASE)
		return
	}
	defer rows.Close()
	var res []*model.Adviser
	if err = scanner.Scan(rows, &res); err != nil {
		response.SendResponse(c, errmsg.ERROR)
		return
	}
	response.SendResponse(c, errmsg.SUCCSE, res)
}
