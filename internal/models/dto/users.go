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
	Id        uint32    `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"not null;unique;type:varchar(255)" json:"email"`
	Username  string    `gorm:"not null;unique;type:varchar(100)" json:"username"`
	Password  string    `gorm:"not null;type:varchar(255)" json:"password"`
	UserRole  uint8     `gorm:"not null;type:tinyint(1)" json:"userRole"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`

	Avatar string `gorm:"not null;type:varchar(255)" json:"avatar"` // 用户头像路径或 URL
	Bio    string `gorm:"not null;type:varchar(255)" json:"bio"`    // 用户简介
}
