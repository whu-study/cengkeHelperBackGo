package services

import (
	database "cengkeHelperBackGo/internal/db"
	"cengkeHelperBackGo/internal/models/dto"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// 在 services/auth_service.go (或类似文件) 中
func CheckUser(email string, password string) (dto.User, bool) {
	var user dto.User
	// 根据邮箱查找用户
	if err := database.Client.Where("email = ?", email).First(&user).Error; err != nil {
		log.Printf("未找到邮箱为 %s 的用户, 错误: %v", email, err)
		return dto.User{}, false // 用户不存在
	}

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("密码验证失败:", err)
		return dto.User{}, false // 密码不匹配
	}
	return user, true // 密码匹配, 返回用户信息
}
