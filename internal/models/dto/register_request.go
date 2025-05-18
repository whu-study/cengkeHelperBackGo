package dto

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Email    string `json:"email" binding:"required,email"`
	// 可选字段
	Avatar string `json:"avatar"`
	Bio    string `json:"bio"`
}
