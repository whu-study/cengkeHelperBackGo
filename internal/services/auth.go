package services

import (
	database "cengkeHelperBackGo/internal/db"
	"cengkeHelperBackGo/internal/models/dto"
	"log"
)

func CheckUser(username, password string) (dto.User, bool) {
	user := dto.User{}
	var count int64
	if err := database.Client.Model(&dto.User{}).
		Where("username = ?", username).
		First(&user).Count(&count).Error; err != nil {
		log.Println(err)
		return user, false
	}

	return user, count > 0 && user.Password == password
}
