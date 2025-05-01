package handlers

import (
	database "cengkeHelperBackGo/internal/db"
	"cengkeHelperBackGo/internal/models/dto"
	"cengkeHelperBackGo/internal/models/vo"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserProfileHandler(c *gin.Context) {
	userId := c.Query("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, vo.NewBadResp("userId is required"))
		return
	}

	user := dto.User{}
	if err := database.Client.
		Model(&dto.User{}).
		Where("id = ?", userId).
		Find(&user).Error; err != nil {

	}

	fmt.Println(user)
	c.JSON(http.StatusOK, vo.RespData{
		Code: 200,
		Data: user,
		Msg:  "success",
	})
	// select * from users where id = userId

}
