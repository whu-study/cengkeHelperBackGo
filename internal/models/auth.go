package models

import "github.com/golang-jwt/jwt/v5"

// UserClaims 用于储存用户的JWT的声明信息
type UserClaims struct {
	Username string `json:"username"`
	Role     uint8  `json:"role"`
	jwt.RegisteredClaims
}
