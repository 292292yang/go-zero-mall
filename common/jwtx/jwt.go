package jwtx

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(secretKey string, now int64, accessExpire int64, userId int64) (string, int64, error) {
	expireAt := now + accessExpire

	claims := make(jwt.MapClaims)
	claims["exp"] = expireAt
	claims["iat"] = now
	claims["userId"] = userId

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expireAt, nil
}

func Now() int64 {
	return time.Now().Unix()
}
