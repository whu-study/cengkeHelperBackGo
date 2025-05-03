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
	c.JSON(http.StatusOK, vo.NewSuccessResp(user, "用户信息查询成功"))
	// select * from users where id = userId

}
func UpdateUserProfileHandler(c *gin.Context) {
	// 获取用户ID
	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, vo.NewBadResp("userId is required"))
		return
	}

	// 解析请求体
	var updateData struct {
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
		Bio      string `json:"bio"`
		// 注意：这里不包含敏感字段如密码、邮箱等
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewBadResp("invalid request data"))
		return
	}

	// 检查用户名是否已存在（如果请求中提供了用户名）
	if updateData.Username != "" {
		var count int64
		if err := database.Client.
			Model(&dto.User{}).
			Where("username = ? AND id != ?", updateData.Username, userId).
			Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, vo.NewBadResp("failed to check username"))
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, vo.RespData{
				Code: config.CodeUsernameExists,
				Msg:  "username already exists",
			})
			return
		}
	}

	// 更新用户信息
	updates := map[string]interface{}{}
	if updateData.Username != "" {
		updates["username"] = updateData.Username
	}
	if updateData.Avatar != "" {
		updates["avatar"] = updateData.Avatar
	}
	if updateData.Bio != "" {
		updates["bio"] = updateData.Bio
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, vo.NewBadResp("no fields to update"))
		return
	}

	// 执行更新
	result := database.Client.
		Model(&dto.User{}).
		Where("id = ?", userId).
		Updates(updates)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, vo.NewBadResp("failed to update user profile"))
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, vo.RespData{
			Code: config.CodeUserNotFound,
			Msg:  "user not found",
		})
		return
	}

	// 返回更新后的用户信息
	updatedUser := dto.User{}
	if err := database.Client.
		Model(&dto.User{}).
		Where("id = ?", userId).
		First(&updatedUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewBadResp("failed to fetch updated user data"))
		return
	}

	c.JSON(http.StatusOK, vo.NewSuccessResp(updatedUser, "用户信息修改成功"))
	return
}
