package handlers

import (
	"cengkeHelperBackGo/internal/config" // 导入配置以获取错误码
	"cengkeHelperBackGo/internal/models/dto"
	"cengkeHelperBackGo/internal/models/vo" // 导入包含 RespData 和辅助函数的包
	"cengkeHelperBackGo/internal/services"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	// "log" // 如果需要记录详细错误
)

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler() *PostHandler {
	return &PostHandler{

		postService: services.NewPostService(),
	}
}

// GetPosts godoc
// @Summary 获取帖子列表
// @Description 根据查询参数获取帖子列表，支持分页、排序和过滤
// @Tags Posts
// @Accept  json
// @Produce  json
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Param sortBy query string false "排序字段和顺序 (例如: createdAt_desc, likesCount_asc)"
// @Param filterText query string false "标题或内容模糊搜索"
// @Param category query string false "分类过滤"
// @Param tag query string false "标签过滤 (单个标签)"
// @Param authorId query int false "作者ID过滤"
// @Success 200 {object} vo.RespData{data=vo.GetPostsResponseDataVO} "成功"
// @Failure 400 {object} vo.RespData "请求参数错误"
// @Failure 500 {object} vo.RespData "服务器内部错误"
// @Router /posts [get]
func (h *PostHandler) GetPosts(c *gin.Context) {
	var params dto.GetPostsParamsDTO
	if err := c.ShouldBindQuery(&params); err != nil {
		// log.Printf("GetPosts: Bad request params: %v\n", err) // 记录具体错误
		vo.RespondError(c, http.StatusBadRequest, config.CodeBadRequest, "请求参数无效", err)
		return
	}

	// 设定默认值 (如果 DTO 中的 default tag 不生效或需要更复杂逻辑)
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 || params.Limit > 100 { // 限制最大 limit
		params.Limit = 10
	}

	responseVO, err := h.postService.GetPosts(&params)
	if err != nil {
		// log.Printf("GetPosts: Failed to get posts from service: %v\n", err) // 记录具体错误
		// 根据 service 层返回的错误类型判断是客户端错误还是服务端错误
		// 此处简化为服务器内部错误
		vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "获取帖子列表失败", err)
		return
	}

	vo.RespondSuccess(c, "帖子列表获取成功", responseVO)
}

// GetPostByID godoc
// @Summary 获取单个帖子详情
// @Description 根据帖子ID获取帖子的详细信息
// @Tags Posts
// @Accept  json
// @Produce  json
// @Param   id   path      int  true  "帖子ID"
// @Success 200 {object} vo.RespData{data=vo.PostVO} "成功，返回帖子详情"
// @Failure 400 {object} vo.RespData "请求参数错误，例如ID格式无效"
// @Failure 404 {object} vo.RespData "帖子未找到"
// @Failure 500 {object} vo.RespData "服务器内部错误"
// @Router /posts/{id} [get]
func (h *PostHandler) GetPostByID(c *gin.Context) {
	idStr := c.Param("id")
	postIDUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		vo.RespondError(c, http.StatusBadRequest, config.CodeBadRequest, "无效的帖子ID格式", err)
		return
	}
	postID := uint32(postIDUint64)

	// --- 获取当前登录用户的ID (如果存在) ---
	var currentUserID *uint32            // 使用指针，因为用户可能未登录
	userIDVal, exists := c.Get("userID") // 假设您的 JWT 中间件将 userID 存放在 context 中
	if exists {
		if id, ok := userIDVal.(uint32); ok { // 确保类型与 User 模型的 ID 类型一致
			currentUserID = &id
		} else if idFloat, ok := userIDVal.(float64); ok { // 有时 JWT 解析出来可能是 float64
			idUint32 := uint32(idFloat)
			currentUserID = &idUint32
		}
		// 可以添加更多类型断言或转换逻辑
	}
	// 如果取不到 userID，currentUserID 将保持为 nil，Service 层会据此判断用户未登录

	postVO, serviceErr := h.postService.GetPostByID(postID, currentUserID) // ★ 传递 currentUserID
	if serviceErr != nil {
		if errors.Is(serviceErr, gorm.ErrRecordNotFound) || serviceErr.Error() == "帖子未找到" {
			vo.RespondError(c, http.StatusNotFound, config.CodeNotFound, "帖子未找到", nil)
		} else {
			vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "获取帖子详情失败", serviceErr)
		}
		return
	}

	vo.RespondSuccess(c, "帖子详情获取成功", postVO)
}

// CreatePost godoc
// @Summary 创建新帖子
// @Description 用户创建一篇新的帖子
// @Tags Posts
// @Accept  json
// @Produce  json
// @Param   Authorization header    string  true  "Bearer <token>"
// @Param   postData      body      dto.CreatePostDTO  true  "帖子内容"
// @Success 201 {object} vo.RespData{data=vo.PostVO} "帖子创建成功"
// @Failure 400 {object} vo.RespData "请求参数错误"
// @Failure 401 {object} vo.RespData "用户未授权"
// @Failure 500 {object} vo.RespData "服务器内部错误"
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var postData dto.CreatePostDTO
	if err := c.ShouldBindJSON(&postData); err != nil {
		vo.RespondError(c, http.StatusBadRequest, config.CodeBadRequest, "请求参数无效: "+err.Error(), nil)
		return
	}

	userIDVal, exists := c.Get("userId") // 从JWT中间件获取 userID
	if !exists {
		vo.RespondError(c, http.StatusUnauthorized, config.CodeUnauthorized, "用户未授权", nil)
		return
	}
	var authorID uint32
	_, err := fmt.Sscanf(userIDVal.(string), "%d", &authorID)
	if err != nil {
		vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "无法获取用户ID", nil)
		return
	}

	createdPostVO, err := h.postService.CreatePost(&postData, authorID)
	if err != nil {
		vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "创建帖子失败: "+err.Error(), nil)
		return
	}

	// HTTP 201 Created
	c.JSON(http.StatusCreated, vo.NewSuccessResp("帖子创建成功", createdPostVO))
}

// --- 新增 UpdatePost Handler ---
// UpdatePost godoc
// @Summary 更新帖子
// @Description 更新指定ID的帖子内容
// @Tags Posts
// @Accept  json
// @Produce  json
// @Param   Authorization header    string  true  "Bearer <token>"
// @Param   id            path      int     true  "帖子ID"
// @Param   postData      body      dto.UpdatePostDTO  true  "要更新的帖子数据"
// @Success 200 {object} vo.RespData{data=vo.PostVO} "帖子更新成功"
// @Failure 400 {object} vo.RespData "请求参数错误"
// @Failure 401 {object} vo.RespData "用户未授权"
// @Failure 403 {object} vo.RespData "无权操作"
// @Failure 404 {object} vo.RespData "帖子未找到"
// @Failure 500 {object} vo.RespData "服务器内部错误"
// @Router /posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	postIDUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		vo.RespondError(c, http.StatusBadRequest, config.CodeBadRequest, "无效的帖子ID格式", err)
		return
	}
	postID := uint32(postIDUint64)

	var postData dto.UpdatePostDTO
	if err := c.ShouldBindJSON(&postData); err != nil {
		vo.RespondError(c, http.StatusBadRequest, config.CodeBadRequest, "请求参数无效: "+err.Error(), nil)
		return
	}

	userIDVal, exists := c.Get("userId")

	if !exists {
		vo.RespondError(c, http.StatusUnauthorized, config.CodeUnauthorized, "用户未授权", nil)
		return
	}
	var authorID uint32
	_, err = fmt.Sscanf(userIDVal.(string), "%d", &authorID)
	if err != nil {
		vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "无法获取用户ID", nil)
		return
	}

	updatedPostVO, serviceErr := h.postService.UpdatePost(postID, authorID, &postData)
	if serviceErr != nil {
		if serviceErr.Error() == "帖子未找到" {
			vo.RespondError(c, http.StatusNotFound, config.CodeNotFound, serviceErr.Error(), nil)
		} else if serviceErr.Error() == "无权修改此帖子" {
			vo.RespondError(c, http.StatusForbidden, config.CodeForbidden, serviceErr.Error(), nil)
		} else {
			vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "更新帖子失败: "+serviceErr.Error(), nil)
		}
		return
	}
	vo.RespondSuccess(c, "帖子更新成功", updatedPostVO)
}

// --- 新增 DeletePost Handler ---
// DeletePost godoc
// @Summary 删除帖子
// @Description 删除指定ID的帖子
// @Tags Posts
// @Accept  json
// @Produce  json
// @Param   Authorization header    string  true  "Bearer <token>"
// @Param   id            path      int     true  "帖子ID"
// @Success 200 {object} vo.RespData{data=nil} "帖子删除成功"
// @Failure 400 {object} vo.RespData "请求参数错误"
// @Failure 401 {object} vo.RespData "用户未授权"
// @Failure 403 {object} vo.RespData "无权操作"
// @Failure 404 {object} vo.RespData "帖子未找到"
// @Failure 500 {object} vo.RespData "服务器内部错误"
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	postIDUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		vo.RespondError(c, http.StatusBadRequest, config.CodeBadRequest, "无效的帖子ID格式", err)
		return
	}
	postID := uint32(postIDUint64)

	userIDVal, exists := c.Get("userId")
	if !exists {
		vo.RespondError(c, http.StatusUnauthorized, config.CodeUnauthorized, "用户未授权", nil)
		return
	}
	var authorID uint32
	_, err = fmt.Sscanf(userIDVal.(string), "%d", &authorID)
	if err != nil {
		vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "无法获取用户ID", nil)
		return
	}
	role := c.Param("role")
	serviceErr := h.postService.DeletePost(postID, authorID, role)
	if serviceErr != nil {
		if serviceErr.Error() == "帖子未找到" {
			vo.RespondError(c, http.StatusNotFound, config.CodeNotFound, serviceErr.Error(), nil)
		} else if serviceErr.Error() == "无权删除此帖子" {
			vo.RespondError(c, http.StatusForbidden, config.CodeForbidden, serviceErr.Error(), nil)
		} else {
			vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "删除帖子失败: "+serviceErr.Error(), nil)
		}
		return
	}
	vo.RespondSuccess(c, "帖子删除成功", nil) // 删除成功通常不返回 data
}

// --- 新增 ToggleLikePost Handler ---
// ToggleLikePost godoc
// @Summary 切换帖子点赞状态
// @Description 用户点赞或取消点赞一个帖子
// @Tags Posts
// @Accept  json
// @Produce  json
// @Param   Authorization header    string  true  "Bearer <token>"
// @Param   id            path      int     true  "帖子ID"
// @Success 200 {object} vo.RespData{data=vo.ToggleLikeResponseDataVO} "操作成功"
// @Failure 400 {object} vo.RespData "请求参数错误"
// @Failure 401 {object} vo.RespData "用户未授权"
// @Failure 404 {object} vo.RespData "帖子未找到"
// @Failure 500 {object} vo.RespData "服务器内部错误"
// @Router /posts/{id}/toggle-like [post]
func (h *PostHandler) ToggleLikePost(c *gin.Context) {
	idStr := c.Param("id")
	postIDUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		vo.RespondError(c, http.StatusBadRequest, config.CodeBadRequest, "无效的帖子ID格式", err)
		return
	}
	postID := uint32(postIDUint64)

	userIDVal, exists := c.Get("userId")
	if !exists {
		vo.RespondError(c, http.StatusUnauthorized, config.CodeUnauthorized, "用户未授权", nil)
		return
	}
	userID64, err := strconv.ParseUint(userIDVal.(string), 10, 32)
	if err != nil {
		vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "无法获取用户ID", nil)
		return
	}
	userID := uint32(userID64)

	responseVO, serviceErr := h.postService.ToggleLikePost(postID, userID)
	if serviceErr != nil {
		if serviceErr.Error() == "帖子未找到" {
			vo.RespondError(c, http.StatusNotFound, config.CodeNotFound, serviceErr.Error(), nil)
		} else {
			vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "操作失败: "+serviceErr.Error(), nil)
		}
		return
	}
	vo.RespondSuccess(c, "操作成功", responseVO)
}

// --- 新增 ToggleCollectPost Handler ---
// ToggleCollectPost godoc
// @Summary 切换帖子收藏状态
// @Description 用户收藏或取消收藏一个帖子
// @Tags Posts
// @Accept  json
// @Produce  json
// @Param   Authorization header    string  true  "Bearer <token>"
// @Param   id            path      int     true  "帖子ID"
// @Success 200 {object} vo.RespData{data=vo.ToggleCollectResponseDataVO} "操作成功"
// @Failure 400 {object} vo.RespData "请求参数错误"
// @Failure 401 {object} vo.RespData "用户未授权"
// @Failure 404 {object} vo.RespData "帖子未找到"
// @Failure 500 {object} vo.RespData "服务器内部错误"
// @Router /posts/{id}/toggle-collect [post]
func (h *PostHandler) ToggleCollectPost(c *gin.Context) {
	idStr := c.Param("id")
	postIDUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		vo.RespondError(c, http.StatusBadRequest, config.CodeBadRequest, "无效的帖子ID格式", err)
		return
	}
	postID := uint32(postIDUint64)

	userIDVal, exists := c.Get("userId")
	if !exists {
		vo.RespondError(c, http.StatusUnauthorized, config.CodeUnauthorized, "用户未授权", nil)
		return
	}
	userID64, err := strconv.ParseUint(userIDVal.(string), 10, 32)
	if err != nil {
		vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "无法获取用户ID", nil)
		return
	}
	userID := uint32(userID64)

	responseVO, serviceErr := h.postService.ToggleCollectPost(postID, userID)
	if serviceErr != nil {
		if serviceErr.Error() == "帖子未找到" {
			vo.RespondError(c, http.StatusNotFound, config.CodeNotFound, serviceErr.Error(), nil)
		} else {
			vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "操作失败: "+serviceErr.Error(), nil)
		}
		return
	}
	vo.RespondSuccess(c, "操作成功", responseVO)
}

// GetActiveUsersHandler godoc
// @Summary 获取最近 N 天发帖最活跃的用户
// @Description 返回最近 N 天内发帖数量最多的用户列表，默认 days=3, limit=10
// @Tags Posts
// @Accept json
// @Produce json
// @Param days query int false "天数，默认3"
// @Param limit query int false "返回数量，默认10"
// @Success 200 {object} vo.RespData{data=[]vo.ActiveUserVO} "成功"
// @Failure 500 {object} vo.RespData "服务器内部错误"
// @Router /posts/active-users [get]
func (h *PostHandler) GetActiveUsersHandler(c *gin.Context) {
	days := 3
	limit := 10
	if dstr := c.Query("days"); dstr != "" {
		if v, err := strconv.Atoi(dstr); err == nil && v > 0 {
			days = v
		}
	}
	if lstr := c.Query("limit"); lstr != "" {
		if v, err := strconv.Atoi(lstr); err == nil && v > 0 {
			limit = v
		}
	}

	users, err := h.postService.GetActiveUsers(days, limit)
	if err != nil {
		vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "获取活跃用户失败", err)
		return
	}
	vo.RespondSuccess(c, "成功", users)
}

// GetCommunityStatsHandler godoc
// @Summary 获取社区统计信息
// @Description 返回总帖子数、注册用户数和今日新增帖子数
// @Tags Posts
// @Accept json
// @Produce json
// @Success 200 {object} vo.RespData{data=vo.CommunityStatsVO} "成功"
// @Failure 500 {object} vo.RespData "服务器内部错误"
// @Router /community/stats [get]
func (h *PostHandler) GetCommunityStatsHandler(c *gin.Context) {
	stats, err := h.postService.GetCommunityStats()
	if err != nil {
		vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "��取社区统计失败", err)
		return
	}
	vo.RespondSuccess(c, "成功", stats)
}
