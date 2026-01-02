package webhttp

func GetMessageByCode(code int) string {
	switch code {
	case SUCCESS:
		return "成功"
	case ERROR_INVALID_PARAMS:
		return "参数错误"
	case ERROR_UNKNOWN_USER:
		return "用户不存在"
	case ERROR_WRONG_PASSWORD:
		return "密码错误"
	case ERROR_AUTH_CHECK_TOKEN_FAIL:
		return "验证失败, 请重新登录"
	case ERROR_AUTH_CHECK_TOKEN_TIMEOUT:
		return "账户已过期, 请重新登录"
	case ERROR_UPLOAD_ERROR:
		return "上传出现错误"
	case ERROR_DATABASE_CREATE_ERROR:
		return "数据库创建出错"
	default:
		return "未知错误"
	}
}
