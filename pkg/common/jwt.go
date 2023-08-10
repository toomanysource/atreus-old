package common

import (
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
)

// ParseToken 接收TokenString进行校验
func ParseToken(tokenKey, tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenKey), nil
	})
	if err != nil {
		log.Errorf("Server failed to convert Token, err :", err.Error())
		return nil, err
	}
	if token.Valid {
		return token, nil
	}
	return nil, errors.New("invalid JWT token")
}

// GetTokenData 获取Token中的用户数据,返回的是map[string]any类型，需要断言
func GetTokenData(token *jwt.Token) (map[string]any, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to extract claims from JWT token")
	}
	_, ok = claims["user_id"]
	if !ok {
		return nil, errors.New("the token does not carry critical data")
	}
	return claims, nil
}
