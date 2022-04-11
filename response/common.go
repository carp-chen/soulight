package response

import (
	"soulight/utils/errmsg"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ResponseWithData struct {
	*Response
	Data interface{} `json:"data"`
}

type ResponseWithToken struct {
	*ResponseWithData
	Token string `json:"token"`
}

func NewResponse(code int) *Response {
	return &Response{code, errmsg.GetErrMsg(code)}
}

func NewResponseWithData(code int, data interface{}) *ResponseWithData {
	return &ResponseWithData{
		NewResponse(code),
		data,
	}
}

func NewResponseWithToken(code int, data interface{}, token string) *ResponseWithToken {
	return &ResponseWithToken{
		NewResponseWithData(code, data),
		token,
	}
}

func SendResponse(c *gin.Context, arg ...interface{}) {
	switch len(arg) {
	case 1:
		code := arg[0].(int)
		c.JSON(200, NewResponse(code))
	case 2:
		code := arg[0].(int)
		c.JSON(200, NewResponseWithData(code, arg[1]))
	case 3:
		code := arg[0].(int)
		token := arg[2].(string)
		c.JSON(200, NewResponseWithToken(code, arg[1], token))
	}
}
