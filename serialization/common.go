package serialization

import "soulight/utils/errmsg"

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
