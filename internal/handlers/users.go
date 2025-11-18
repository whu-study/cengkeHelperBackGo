package handlers

import (
	"cengkeHelperBackGo/internal/config"
	database "cengkeHelperBackGo/internal/db"
	"cengkeHelperBackGo/internal/models/dto"
	"cengkeHelperBackGo/internal/models/vo"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserProfileHandler(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, vo.NewBadResp("userId is required"))
		return
	}

	// we store userId as string in context; assert that here
	uid, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, vo.NewBadResp("invalid userId type"))
		return
	}

	// fetch basic user
	var user dto.User
	if err := database.Client.Model(&dto.User{}).Where("id = ?", uid).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, vo.RespData{
			Code: config.CodeUserNotFound,
			Msg:  "failed to find user",
		})
		return
	}

	// compute counts
	var postsCount int64
	_ = database.Client.Model(&dto.Post{}).Where("author_id = ?", uid).Count(&postsCount).Error

	var commentsCount int64
	_ = database.Client.Model(&dto.Comment{}).Where("author_id = ?", uid).Count(&commentsCount).Error

	// likes count: count of post likes given + comment likes given
	var postLikesGiven int64
	_ = database.Client.Model(&dto.UserPostLike{}).Where("user_id = ?", uid).Count(&postLikesGiven).Error
	var commentLikesGiven int64
	_ = database.Client.Model(&dto.UserCommentLike{}).Where("user_id = ?", uid).Count(&commentLikesGiven).Error

	likesCount := postLikesGiven + commentLikesGiven

	// likes received: sum of likes_count on user's posts + likes_count on user's comments
	var postLikesReceived int64
	_ = database.Client.Model(&dto.Post{}).Select("COALESCE(SUM(likes_count),0)").Where("author_id = ?", uid).Scan(&postLikesReceived).Error
	var commentLikesReceived int64
	_ = database.Client.Model(&dto.Comment{}).Select("COALESCE(SUM(likes_count),0)").Where("author_id = ?", uid).Scan(&commentLikesReceived).Error
	likesReceived := postLikesReceived + commentLikesReceived

	profile := vo.ExtendedUserProfileVO{
		UserProfileVO: vo.UserProfileVO{
			ID:        user.Id,
			Email:     user.Email,
			Username:  user.Username,
			UserRole:  user.UserRole,
			CreatedAt: user.CreatedAt,
			Avatar:    user.Avatar,
			Bio:       user.Bio,
		},
		PostsCount:    postsCount,
		CommentsCount: commentsCount,
		LikesCount:    likesCount,
		LikesReceived: likesReceived,
	}

	c.JSON(http.StatusOK, vo.NewSuccessResp("用户信息查询成功", profile))
	// select * from users where id = userId

}
func UpdateUserProfileHandler(c *gin.Context) {
	// 获取用户ID
	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, vo.NewBadResp("userId is required"))
		return
	}

	// expect userId in context to be string
	uid, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, vo.NewBadResp("invalid userId type"))
		return
	}

	// 解析请求体
	var updateData struct {
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
		Bio      string `json:"bio"`
		Email    string `json:"email"`
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
			Where("username = ? AND id != ?", updateData.Username, uid).
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
	if updateData.Email != "" {
		updates["email"] = updateData.Email
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, vo.NewBadResp("no fields to update"))
		return
	}

	// 执行更新
	result := database.Client.
		Model(&dto.User{}).
		Where("id = ?", uid).
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
		Where("id = ?", uid).
		First(&updatedUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewBadResp("failed to fetch updated user data"))
		return
	}

	c.JSON(http.StatusOK, vo.NewSuccessResp("用户信息修改成功", updatedUser))
	return
}
