package config

import (
	"cengkeHelperBackGo/internal/models/dto"
	"fmt"
	"testing"
)

func TestMysql(t *testing.T) {
	LoadConfig("../config/config.yaml")
	user := dto.User{
		Username: "testuser",
		Password: "testpassword",
	}
	_ = Client.Create(&user)

	users := make([]dto.User, 0)
	_ = Client.Find(&users)

	fmt.Println(users)
}
