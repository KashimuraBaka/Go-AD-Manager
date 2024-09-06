package middleware

import (
	"net/http"
	"time"

	"gitee.com/Kashimura/go-baka-control/utils/webhttp"
	"gitee.com/Kashimura/go-baka-control/utils/webhttp/jwt"
	"github.com/gin-gonic/gin"
)

func VerifyAccountToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code := webhttp.SUCCESS
		token := ctx.Request.Header.Get("Authorization")

		if token == "" {
			code = webhttp.ERROR_INVALID_PARAMS
		} else {
			// 解析token
			claims, err := jwt.ParseUserToken(token)
			if err != nil {
				// 非法token
				code = webhttp.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				// 过期token
				code = webhttp.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		// 返回结果
		if code != webhttp.SUCCESS {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, webhttp.Response{
				Code: code,
				Msg:  webhttp.GetMessageByCode(code),
			})
			return
		}

		ctx.Next()
	}
}
