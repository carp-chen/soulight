package errmsg

const (
	SUCCSE         = 200
	ERROR          = 500
	INVALID_PARAMS = 400

	// code= 10000... 用户与顾问模块的错误
	ERROR_USERNAME_USED    = 10001
	ERROR_PASSWORD_WRONG   = 10002
	ERROR_USER_NOT_EXIST   = 10003
	ERROR_GENARATE_TOKEN   = 10004
	ERROR_TOKEN_NOT_EXIST  = 10005
	ERROR_TOKEN_TIMEOUT    = 10006
	ERROR_TOKEN_WRONG      = 10007
	ERROR_TOKEN_TYPE_WRONG = 10008
	ERROR_USER_NO_RIGHT    = 10009

	// code= 20000... 订单模块的错误
	ERROR_COINS_NOT_ENOUGH = 20001
	ERROR_ORDER_NOT_EXIST  = 20002
	ERROR_ORDER_TIMEOUT    = 20003

	//code=30000 数据库错误
	ERROR_DATABASE = 30000
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
	ERROR_COINS_NOT_ENOUGH: "金币不足",
	ERROR_ORDER_NOT_EXIST:  "订单不存在",
	ERROR_ORDER_TIMEOUT:    "订单已过期",
	ERROR_DATABASE:         "数据库操作出错,请重试",
}

func GetErrMsg(code int) string {
	msg, ok := codeMsg[code]
	if ok {
		return msg
	}
	return codeMsg[ERROR]
}
