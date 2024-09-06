package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

// UserClaims 自定义声明结构体并内嵌 jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段，若需要额外记录其他字段，就可以自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中

type UserClaims struct {
	jwt.StandardClaims        // 必须
	UserID             int    `json:"user_id"`
	UserName           string `json:"user_name"`
}

func GenUserToken(userID int, username string, timeout time.Duration) (string, error) {
	// 创建一个我们自己的声明的数据
	c := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(timeout).Unix(), // 过期时间
			Issuer:    "Kashimrua",                    // 签发人
		},
		UserID:   userID,
		UserName: username,
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseUserToken 解析JWT
func ParseUserToken(tokenString string) (*UserClaims, error) {
	// 解析token
	var mc = new(UserClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, keyFunc)
	if err != nil {
		return nil, err
	}
	// 校验token
	if token.Valid {
		return mc, nil
	}
	return nil, errors.New("invalid token")
}

// RefreshUserToken 刷新AccessToken
func RefreshUserToken(token string) (string, error) {
	// 从旧access token中解析出claims数据
	var claims UserClaims
	_, err := jwt.ParseWithClaims(token, &claims, keyFunc)
	if err != nil {
		v, _ := err.(*jwt.ValidationError)
		// 当access token是过期错误 并且 refresh token没有过期时就创建一个新的access token
		if v.Errors == jwt.ValidationErrorExpired {
			return GenUserToken(claims.UserID, claims.UserName, 1*time.Hour)
		}
		return "", err
	}
	return GenUserToken(claims.UserID, claims.UserName, 1*time.Hour)
}
