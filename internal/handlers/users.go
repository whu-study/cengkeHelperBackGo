package handlers

import (
	"cengkeHelperBackGo/internal/config"
	database "cengkeHelperBackGo/internal/db"
	"cengkeHelperBackGo/internal/models/dto"
	"cengkeHelperBackGo/internal/models/vo"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserProfileHandler(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, vo.NewBadResp("userId is required"))
		return
	}

	user := dto.User{}
	if err := database.Client.
		Model(&dto.User{}).
		Where("id = ?", userId).
		First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, vo.RespData{
			Code: config.CodeUserNotFound,
			Msg:  "failed to find user",
		})
		return
	}

	fmt.Println(user)
	c.JSON(http.StatusOK, vo.NewSuccessResp(user))
	// select * from users where id = userId

}
