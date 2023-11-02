package global

const (
	RES_SUCCESS = 200
	RES_ERROR   = 201
)

type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

func ResponseError(msg string) *Response {
	return &Response{
		Code: RES_ERROR,
		Data: nil,
		Msg:  msg,
	}
}

func ResponseSuccess(data any, msg string) *Response {
	return &Response{
		Code: RES_SUCCESS,
		Data: data,
		Msg:  msg,
	}
}
