package auth

import (
	database "cengkeHelperBackGo/internal/db"
	"cengkeHelperBackGo/internal/models/dto"
	"cengkeHelperBackGo/internal/services"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

var emailService = services.NewEmailService()

func checkUser(email string, password string) (dto.User, bool) {
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

// SendEmailCode 发送邮箱验证码
func SendEmailCode(email string) (string, error) {
	return emailService.SendVerificationCode(email)
}

// checkEmailCode 验证邮箱验证码
func checkEmailCode(email, emailCode string) bool {
	return emailService.VerifyEmailCode(email, emailCode)
}
