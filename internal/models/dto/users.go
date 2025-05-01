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
	Id        uint32    `gorm:"primaryKey"`
	Email     string    `gorm:"not null;unique;type:varchar(255)"`
	Username  string    `gorm:"not null;unique;type:varchar(100)"`
	Password  string    `gorm:"not null;type:varchar(255)"`
	UserRole  uint8     `gorm:"not null;type:tinyint(1)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	Avatar string `gorm:"not null;type:varchar(255)"` // 用户头像路径或 URL
	Bio    string `gorm:"not null;type:varchar(255)"` // 用户简介
}
