package util

import (
	"errors"
	"tiktok/cmd/api/global"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	UserID uint64
	jwt.RegisteredClaims
}

// GenerateToken 生成 token
func GenerateToken(userID uint64) (string, error) {
	// 获取签名密钥
	signingKey := []byte(global.Config.JWTConfig.SigningKey)
	// 生成token
	claims := UserClaims{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

// ParseToken 解析 token
func ParseToken(tokenString string) (*UserClaims, error) {
	// 获取签名密钥
	signingKey := []byte(global.Config.JWTConfig.SigningKey)
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	} else if token == nil {
		return nil, errors.New("invalid token")
	}
	if claims, ok := token.Claims.(*UserClaims); ok {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
