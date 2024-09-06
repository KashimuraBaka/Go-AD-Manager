package webhttp

const (
	// 成功
	SUCCESS = 200

	// 错误参数
	ERROR_INVALID_PARAMS = 410
	// 用户不存在
	ERROR_UNKNOWN_USER = 411
	// 密码错误
	ERROR_WRONG_PASSWORD = 412
	// Token错误
	ERROR_AUTH_CHECK_TOKEN_FAIL = 415
	// Token超时
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 416
	// 上传错误
	ERROR_UPLOAD_ERROR = 420
	// 数据创建失败
	ERROR_DATABASE_CREATE_ERROR = 430
	// 其他错误
	ERROR_OTHER = 500
)

type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}
