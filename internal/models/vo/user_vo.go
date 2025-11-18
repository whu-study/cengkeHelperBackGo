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

// ExtendedUserProfileVO 在 UserProfileVO 基础上增加统计字段，返回给前端
type ExtendedUserProfileVO struct {
	UserProfileVO
	PostsCount    int64 `json:"postsCount"`
	CommentsCount int64 `json:"commentsCount"`
	LikesCount    int64 `json:"likesCount"`    // 点赞数（用户自己做出的点赞数）
	LikesReceived int64 `json:"likesReceived"` // 收到的点赞数（别人给该用户的帖子/评论点赞）
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
