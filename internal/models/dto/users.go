package dto

import (
	"time"
)

const (
	UserRoleCommon uint8 = 0
	UserRoleAdmin  uint8 = 1
)

type SimpleUser struct {
	Id    uint32
	Email string
}
type User struct {
	Id        uint32    `gorm:"not null;type:uint;primaryKey" json:"id"`
	Email     string    `gorm:"not null;unique;type:varchar(255)" json:"email"`
	Username  string    `gorm:"not null;unique;type:varchar(100)" json:"username"`
	Password  string    `gorm:"not null;type:varchar(255)" json:"-"` // 密码不返回给前端
	UserRole  uint8     `gorm:"not null;type:tinyint(1)" json:"userRole"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`

	Avatar string `gorm:"not null;type:varchar(255)" json:"avatar"` // 用户头像路径或 URL
	Bio    string `gorm:"not null;type:varchar(255)" json:"bio"`    // 用户简介
}

// RegisterRequestDTO 对应前端注册请求的数据
type RegisterRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=100"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	// UserRole uint8  `json:"userRole"` // 通常角色由后端分配或有默认值
	// Bio      string `json:"bio,omitempty"`
}

// LoginRequestDTO 对应前端登录请求的数据
type LoginRequestDTO struct {
	Account  string `json:"account" binding:"required"` // 可以是用户名或邮箱
	Password string `json:"password" binding:"required"`
	// Captcha  string `json:"captcha,omitempty"` // 如果有验证码
}

// UpdateUserProfileDTO 用于更新用户个人资料的请求
type UpdateUserProfileDTO struct {
	Username *string `json:"username,omitempty"` // 使用指针表示可选更新
	Email    *string `json:"email,omitempty"`
	Bio      *string `json:"bio,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
	// Password *string `json:"password,omitempty"` // 修改密码通常有单独接口
}
