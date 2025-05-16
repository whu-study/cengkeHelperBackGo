package vo

import "time"

// --- DTOs (Data Transfer Objects) ---

// --- VOs (View Objects) ---
// UserProfileVO 对应前端的 UserProfile 类型，用于API响应
type UserProfileVO struct {
	ID        uint32    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	UserRole  uint8     `json:"userRole"` // 前端类型中为 role: number | string
	CreatedAt time.Time `json:"createdAt"`
	Avatar    string    `json:"avatar,omitempty"`
	Bio       string    `json:"bio,omitempty"`
}

// LoginResponseDataVO 对应前端 LoginResponseData
type LoginResponseDataVO struct {
	Token string        `json:"token"`
	User  UserProfileVO `json:"user,omitempty"` // 登录后可选返回用户信息
}

// RegisterResponseDataVO 对应前端 RegisterResponseData
type RegisterResponseDataVO struct {
	Token string        `json:"token"`
	User  UserProfileVO `json:"user"` // 注册成功后返回创建的用户信息
}

// AuthorInfoVO 用于帖子和评论中嵌入的作者信息 (对应前端的 AuthorInfo pick)
type AuthorInfoVO struct {
	ID       uint32 `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar,omitempty"`
}
