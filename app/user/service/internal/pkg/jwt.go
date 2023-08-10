package pkg

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// 可提到 config.yaml
const (
	JwtSignKey = "AtReUs"
	JwtExpired = 60 * 60
)

func ProduceToken(userId uint32) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userId
	expiredTime := time.Now().Add(JwtExpired * time.Second)
	claims["exp"] = expiredTime.Unix()

	mySigningKey := []byte(JwtSignKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}
