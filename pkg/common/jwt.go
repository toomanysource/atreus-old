package common

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"time"
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

// ProduceToken 生成Token
func ProduceToken(tokenKey string, userId uint32, expired time.Duration) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userId
	expiredTime := time.Now().Add(expired)
	claims["exp"] = expiredTime.Unix()
	mySigningKey := []byte(tokenKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

// GenSaltPassword 生成一个含有盐值的密码字符串
func GenSaltPassword(salt, password string) string {
	// 创建一个 sha256 的哈希算法实例
	s1 := sha256.New()
	// 密码转化为字符数组
	s1.Write([]byte(password))
	// 使用 s1 进行哈希运算，并转化为字符串
	str1 := fmt.Sprintf("%x", s1.Sum(nil))

	// 创建另外一个 sha256 哈希算法，并且将 str1 和 salt 连接起来，转换为字符串，并且使用 s2 进行哈希运算
	s2 := sha256.New()
	s2.Write([]byte(str1 + salt))
	return fmt.Sprintf("%x", s2.Sum(nil))
}
