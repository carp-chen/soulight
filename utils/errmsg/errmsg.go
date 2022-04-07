package errmsg

const (
	SUCCSE         = 200
	ERROR          = 500
	INVALID_PARAMS = 400

	// code= 1000... 用户模块的错误
	ERROR_USERNAME_USED    = 1001
	ERROR_PASSWORD_WRONG   = 1002
	ERROR_USER_NOT_EXIST   = 1003
	ERROR_GENARATE_TOKEN   = 1004
	ERROR_TOKEN_NOT_EXIST  = 1005
	ERROR_TOKEN_TIMEOUT    = 1006
	ERROR_TOKEN_WRONG      = 1007
	ERROR_TOKEN_TYPE_WRONG = 1008
	ERROR_USER_NO_RIGHT    = 1009

	// code= 2000... 顾问模块的错误

	//code=3000 数据库错误
	ERROR_DATABASE = 3000
)

var codeMsg = map[int]string{
	SUCCSE:                 "OK",
	ERROR:                  "FAIL",
	INVALID_PARAMS:         "请求参数错误",
	ERROR_USERNAME_USED:    "用户名已存在！",
	ERROR_PASSWORD_WRONG:   "密码错误",
	ERROR_USER_NOT_EXIST:   "用户不存在",
	ERROR_GENARATE_TOKEN:   "Token生成失败",
	ERROR_TOKEN_NOT_EXIST:  "TOKEN不存在,请重新登陆",
	ERROR_TOKEN_TIMEOUT:    "TOKEN已过期,请重新登陆",
	ERROR_TOKEN_WRONG:      "TOKEN不正确,请重新登陆",
	ERROR_TOKEN_TYPE_WRONG: "TOKEN格式错误,请重新登陆",
	ERROR_USER_NO_RIGHT:    "该用户无权限",
	ERROR_DATABASE:         "数据库操作出错,请重试",
}

func GetErrMsg(code int) string {
	msg, ok := codeMsg[code]
	if ok {
		return msg
	}
	return codeMsg[ERROR]
}
