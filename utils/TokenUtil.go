package utils

import "github.com/dgrijalva/jwt-go"

func ParseToken(tokenSrt string, SecretKey []byte) (claims jwt.Claims, err error) {
	token, err := jwt.Parse(tokenSrt, func(*jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	claims = token.Claims
	return
}
