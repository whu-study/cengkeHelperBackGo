package handlers

import "github.com/golang-jwt/jwt/v5"

// UserClaims 用于储存用户的JWT的声明信息, 这里负载只携带openid
type UserClaims struct {
	Openid string `json:"openid"`
	jwt.RegisteredClaims
}
