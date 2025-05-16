package handlers // Or your actual handler package name

import (
	"cengkeHelperBackGo/internal/config"
	"cengkeHelperBackGo/internal/models/dto"
	"cengkeHelperBackGo/internal/models/vo" // Your common response is here
	"cengkeHelperBackGo/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type CommentHandler struct {
	commentService *services.CommentService
}

func NewCommentHandler() *CommentHandler {
	return &CommentHandler{
		commentService: services.NewCommentService(),
	}
}

// Helper to get user ID from context; adjust key and type as per your auth middleware
func getUserIDFromContext(c *gin.Context) (*uint32, bool) {
	userIDVal, exists := c.Get("userId") // Ensure key matches ("userId", "userID", etc.)
	if !exists {
		return nil, false
	}
	// Handle different possible types from context
	var userID uint32
	switch v := userIDVal.(type) {
	case uint:
		userID = uint32(v)
	case uint32:
		userID = v
	case int:
		userID = uint32(v)
	case int32:
		userID = uint32(v)
	case int64:
		userID = uint32(v)
	case float64: // JWT might store numbers as float64
		userID = uint32(v)
	case string:
		parsedID, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return nil, false
		}
		userID = uint32(parsedID)
	default:
		return nil, false
	}
	return &userID, true
}

// Helper to get user role from context; adjust as needed
func getUserRoleFromContext(c *gin.Context) string {
	roleVal, exists := c.Get("role") // Assuming role is stored
	if !exists {
		return "" // Default role or indication of no specific role
	}
	role, _ := roleVal.(string)
	return role
}

// GetCommentsByPostID godoc
// @Summary 根据帖子ID获取评论列表
// @Description Fetches a list of comments for a given post ID.
// @Tags Comments
// @Accept json
// @Produce json
// @Param postId path string true "帖子ID"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Param sortBy query string false "排序 (eg: createdAt_desc)"
// @Success 200 {object} vo.RespData{data=vo.GetCommentsResponseDataVO} "成功"
// @Failure 400 {object} vo.RespData "请求参数错误"
// @Failure 404 {object} vo.RespData "帖子未找到"
// @Failure 500 {object} vo.RespData "服务器内部错误"
// @Router /posts/{postId}/comments [get]
func (h *CommentHandler) GetCommentsByPostID(c *gin.Context) {
	postIDStr := c.Param("postId")
	postIDUint64, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		vo.RespondError(c, http.StatusBadRequest, config.CodeBadRequest, "无效的帖子ID格式", err)
		return
	}
	postID := uint32(postIDUint64)

	var params dto.GetCommentsParamsDTO
	if err := c.ShouldBindQuery(&params); err != nil {
		vo.RespondError(c, http.StatusBadRequest, config.CodeBadRequest, "请求参数无效", err)
		return
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 || params.Limit > 100 {
		params.Limit = 10
	}

	currentUserID, _ := getUserIDFromContext(c) // User might not be logged in, so error is ignored for reads

	responseVO, serviceErr := h.commentService.GetCommentsByPostID(postID, &params, currentUserID)
	if serviceErr != nil {
		if serviceErr.Error() == "帖子未找到" {
			vo.RespondError(c, http.StatusNotFound, config.CodeNotFound, serviceErr.Error(), nil)
		} else {
			vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "获取评论列表失败", serviceErr)
		}
		return
	}
	vo.RespondSuccess(c, "评论列表获取成功", responseVO)
}

// AddComment godoc
// @Summary 添加新评论
// @Description Adds a new comment (top-level or reply).
// @Tags Comments
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param commentData body dto.AddCommentDTO true "评论数据"
// @Success 201 {object} vo.RespData{data=vo.CommentVO} "评论创建成功"
// @Failure 400 {object} vo.RespData "请求参数错误"
// @Failure 401 {object} vo.RespData "用户未授权"
// @Failure 404 {object} vo.RespData "相关资源未找到 (帖子/父评论/回复用户)"
// @Failure 500 {object} vo.RespData "服务器内部错误"
// @Router /comments [post]
func (h *CommentHandler) AddComment(c *gin.Context) {
	authorIDPtr, authenticated := getUserIDFromContext(c)
	if !authenticated || authorIDPtr == nil {
		vo.RespondError(c, http.StatusUnauthorized, config.CodeUnauthorized, "用户未授权或无法获取用户ID", nil)
		return
	}
	authorID := *authorIDPtr

	var commentData dto.AddCommentDTO
	if err := c.ShouldBindJSON(&commentData); err != nil {
		vo.RespondError(c, http.StatusBadRequest, config.CodeBadRequest, "请求参数无效", err)
		return
	}

	// postId from path could also be used/validated if route was /posts/:postId/comments
	// postIdFromPathStr := c.Param("postId")

	createdCommentVO, serviceErr := h.commentService.AddComment(&commentData, authorID)
	if serviceErr != nil {
		// More granular error checking based on service error messages
		errMsg := serviceErr.Error()
		switch {
		case strings.Contains(errMsg, "未找到"): // "帖子未找到", "父评论未找到", etc.
			vo.RespondError(c, http.StatusNotFound, config.CodeNotFound, errMsg, nil)
		case strings.Contains(errMsg, "不匹配"): // "父评论与当前帖子不匹配"
			vo.RespondError(c, http.StatusBadRequest, config.CodeBadRequest, errMsg, nil)
		default:
			vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "添加评论失败", serviceErr)
		}
		return
	}
	c.JSON(http.StatusCreated, vo.NewSuccessResp("评论添加成功", createdCommentVO))
}

// DeleteComment godoc
// @Summary 删除指定ID的评论
// @Description Deletes a comment with the given ID. Requires authentication and authorization.
// @Tags Comments
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param commentId path string true "评论ID"
// @Success 200 {object} vo.RespData{data=nil} "评论删除成功"
// @Failure 400 {object} vo.RespData "无效的评论ID格式"
// @Failure 401 {object} vo.RespData "用户未授权"
// @Failure 403 {object} vo.RespData "无权限删除"
// @Failure 404 {object} vo.RespData "评论未找到"
// @Failure 500 {object} vo.RespData "服务器内部错误"
// @Router /comments/{commentId} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentIDStr := c.Param("commentId")
	commentIDUint64, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		vo.RespondError(c, http.StatusBadRequest, config.CodeBadRequest, "无效的评论ID格式", err)
		return
	}
	commentID := uint32(commentIDUint64)

	userIDPtr, authenticated := getUserIDFromContext(c)
	if !authenticated || userIDPtr == nil {
		vo.RespondError(c, http.StatusUnauthorized, config.CodeUnauthorized, "用户未授权", nil)
		return
	}
	userID := *userIDPtr
	userRole := getUserRoleFromContext(c)

	serviceErr := h.commentService.DeleteComment(commentID, userID, userRole)
	if serviceErr != nil {
		errMsg := serviceErr.Error()
		switch {
		case errMsg == "评论未找到":
			vo.RespondError(c, http.StatusNotFound, config.CodeNotFound, errMsg, nil)
		case errMsg == "无权限删除此评论":
			vo.RespondError(c, http.StatusForbidden, config.CodeForbidden, errMsg, nil)
		default:
			vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "删除评论失败", serviceErr)
		}
		return
	}
	vo.RespondSuccess(c, "评论删除成功", nil)
}

// ToggleLikeComment godoc
// @Summary 切换评论点赞状态
// @Description Toggles the like status of a comment. Requires authentication.
// @Tags Comments
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param commentId path string true "评论ID"
// @Success 200 {object} vo.RespData{data=vo.ToggleLikeCommentResponseDataVO} "操作成功"
// @Failure 400 {object} vo.RespData "无效的评论ID格式"
// @Failure 401 {object} vo.RespData "用户未授权"
// @Failure 404 {object} vo.RespData "评论未找到"
// @Failure 500 {object} vo.RespData "服务器内部错误"
// @Router /comments/{commentId}/toggle-like [post]
func (h *CommentHandler) ToggleLikeComment(c *gin.Context) {
	commentIDStr := c.Param("commentId")
	commentIDUint64, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		vo.RespondError(c, http.StatusBadRequest, config.CodeBadRequest, "无效的评论ID格式", err)
		return
	}
	commentID := uint32(commentIDUint64)

	userIDPtr, authenticated := getUserIDFromContext(c)
	if !authenticated || userIDPtr == nil {
		vo.RespondError(c, http.StatusUnauthorized, config.CodeUnauthorized, "用户未授权", nil)
		return
	}
	userID := *userIDPtr

	responseVO, serviceErr := h.commentService.ToggleLikeComment(commentID, userID)
	if serviceErr != nil {
		if serviceErr.Error() == "评论未找到" {
			vo.RespondError(c, http.StatusNotFound, config.CodeNotFound, serviceErr.Error(), nil)
		} else {
			vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "操作点赞失败", serviceErr)
		}
		return
	}
	vo.RespondSuccess(c, "操作成功", responseVO)
}
