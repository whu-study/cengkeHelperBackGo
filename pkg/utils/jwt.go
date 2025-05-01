package utils

import (
	"cengkeHelperBackGo/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// UserClaims 用于储存用户的JWT的声明信息
type UserClaims struct {
	Username string `json:"username"`
	UserId   uint32 `json:"userId"`
	Role     uint8  `json:"role"`
	jwt.RegisteredClaims
}

func ParseUserJwt(token string) (*UserClaims, error) {
	t, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Conf.JwtSecurityKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Go 语言中的类型断言（Type Assertion）语法
	// 如果类型转换成功，则ok为true
	// 否则ok为false
	if claims, ok := t.Claims.(*UserClaims); ok && t.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func GenerateUserToken(username string, role uint8) (string, error) {

	// 5 天有效期
	expirationTime := time.Now().Add(5 * 24 * time.Hour)
	// 验证成功，准备签名生成token
	claims := &UserClaims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // 过期时间
		},
	}

	// 因为服务器性能问题，这里用一个简单的算法(HS256)，配上私钥，生成令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Conf.JwtSecurityKey))
	return tokenString, err

}
