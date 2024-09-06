package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

// MyClaims 自定义声明结构体并内嵌 jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段，若需要额外记录其他字段，就可以自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中

type FileClaims struct {
	jwt.StandardClaims
	FileID int64 `json:"file_id"`
}

func GenFileToken(fileID int64, timeout time.Duration) (string, error) {
	// 创建一个我们自己的声明的数据
	c := FileClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(timeout).Unix(), // 过期时间
			Issuer:    "Kashimrua",                    // 签发人
		},
		FileID: fileID,
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseFileToken 解析JWT
func ParseFileToken(tokenString string) (*FileClaims, error) {
	// 解析token
	var fc = new(FileClaims)
	token, err := jwt.ParseWithClaims(tokenString, fc, keyFunc)
	if err != nil {
		return nil, err
	}
	// 校验token
	if token.Valid {
		return fc, nil
	}
	return nil, errors.New("invalid token")
}
