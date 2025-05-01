package database

import (
	"cengkeHelperBackGo/internal/models"
	"fmt"
	"testing"
)

func TestMysql(t *testing.T) {

	user := models.User{
		Username: "testuser",
		Password: "testpassword",
	}
	_ = Client.Create(&user)

	users := make([]models.User, 0)
	_ = Client.Find(&users)

	fmt.Println(users)
}
