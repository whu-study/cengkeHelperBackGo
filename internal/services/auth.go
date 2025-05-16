package services

import (
	database "cengkeHelperBackGo/internal/db"
	"cengkeHelperBackGo/internal/models/dto"
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

	return user, user.Password == password // 验证成功
}
