package jwt

import "github.com/golang-jwt/jwt"

var mySecret = []byte("KashimuraSrcds")

func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return mySecret, nil
}
