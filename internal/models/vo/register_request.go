package vo

type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=100"`
	Password  string `json:"password" binding:"required,min=6"`
	Email     string `json:"email" binding:"required,email"`
	EmailCode string `json:"emailCode" binding:"required,len=6"`
	// 可选字段
	Avatar string `json:"avatar"`
	Bio    string `json:"bio"`
}
