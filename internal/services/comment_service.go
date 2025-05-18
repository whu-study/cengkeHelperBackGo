package services

import (
	database "cengkeHelperBackGo/internal/db" // 确保与 post_service.go 一致
	"cengkeHelperBackGo/internal/models/dto"
	"cengkeHelperBackGo/internal/models/vo"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"strings"
)

type CommentService struct {
	// 不需要数据库字段，因为我们直接使用 database.Client
}

func NewCommentService() *CommentService {
	return &CommentService{}
}

// convertCommentToVO 辅助函数：将 dto.Comment (模型) 转换为 vo.CommentVO
func convertCommentToVO(comment dto.Comment, currentUserID *uint32, fetchChildren bool, maxDepth int, currentDepth int) (*vo.CommentVO, error) {
	authorVO := vo.AuthorInfoVO{}
	if comment.Author.Id != 0 { // 假设 Author 结构体已预加载，且 Id 是其有效性标识
		authorVO.ID = comment.Author.Id
		authorVO.Username = comment.Author.Username
		authorVO.Avatar = comment.Author.Avatar
	}

	var replyToUserVO *vo.AuthorInfoVO
	if comment.ReplyToUser != nil && comment.ReplyToUser.Id != 0 { // 假设 ReplyToUser 已预加载
		replyToUserVO = &vo.AuthorInfoVO{
			ID:       comment.ReplyToUser.Id,
			Username: comment.ReplyToUser.Username,
			Avatar:   comment.ReplyToUser.Avatar,
		}
	}

	likesCount := int(comment.LikesCount)
	isLiked := false
	if currentUserID != nil && *currentUserID > 0 {
		var likeRecordCount int64
		// 使用 database.Client 构建查询
		queryLike := database.Client.Model(&dto.UserCommentLike{}).Where("user_id = ? AND comment_id = ?", *currentUserID, comment.ID)
		if err := queryLike.Count(&likeRecordCount).Error; err != nil {
			log.Printf("检查评论 %d 的点赞状态失败 (用户 %d): %v", comment.ID, *currentUserID, err)
			// 此处不中断，isLiked 默认为 false
		}
		if likeRecordCount > 0 {
			isLiked = true
		}
	}

	var childrenVO []vo.CommentVO
	if fetchChildren && currentDepth < maxDepth {
		var childModels []dto.Comment
		// 使用 database.Client 构建查询
		queryChildren := database.Client.Model(&dto.Comment{}).
			Preload("Author").
			Preload("ReplyToUser").
			Where("parent_id = ?", comment.ID).
			Order("created_at ASC") // 链式操作

		if err := queryChildren.Find(&childModels).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("获取评论 %d 的子评论失败: %v", comment.ID, err)
			}
			// 即使获取子评论失败，也继续，只是子评论列表为空
		}

		for _, childModel := range childModels {
			childVO, convErr := convertCommentToVO(childModel, currentUserID, fetchChildren, maxDepth, currentDepth+1)
			if convErr == nil && childVO != nil {
				childrenVO = append(childrenVO, *childVO)
			} else if convErr != nil {
				log.Printf("转换子评论 %d (父评论 %d) 到 VO 失败: %v", childModel.ID, comment.ID, convErr)
			}
		}
	}

	commentVO := &vo.CommentVO{
		ID:                   comment.ID,
		PostID:               comment.PostID,
		Author:               authorVO,
		Content:              comment.Content,
		CreatedAt:            comment.CreatedAt,
		UpdatedAt:            &comment.UpdatedAt,
		LikesCount:           &likesCount,
		ParentID:             comment.ParentID,
		ReplyToUser:          replyToUserVO,
		IsLikedByCurrentUser: &isLiked,
		Children:             childrenVO,
	}
	if comment.UpdatedAt.IsZero() {
		commentVO.UpdatedAt = nil
	}

	return commentVO, nil
}

func (s *CommentService) GetCommentsByPostID(postID uint32, params *dto.GetCommentsParamsDTO, currentUserID *uint32) (*vo.GetCommentsResponseDataVO, error) {
	var comments []dto.Comment
	var total int64

	// 1. 检查帖子是否存在
	var post dto.Post
	if err := database.Client.First(&post, postID).Error; err != nil { // 使用 database.Client
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("帖子未找到")
		}
		return nil, fmt.Errorf("验证帖子有效性失败: %w", err)
	}

	// 2. 构建基础查询，查询顶级评论
	query := database.Client.Model(&dto.Comment{}).Where("post_id = ? AND parent_id IS NULL", postID) // 使用 database.Client

	// 3. 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("统计评论总数失败: %w", err)
	}

	// 4. 应用排序 (确保 query 被重新赋值)
	if params.SortBy != "" {
		parts := strings.Split(params.SortBy, "_")
		if len(parts) == 2 {
			column := parts[0]
			order := strings.ToUpper(parts[1])
			allowedSorts := map[string]string{"createdAt": "created_at", "likesCount": "likes_count"}
			if dbColumn, ok := allowedSorts[column]; ok {
				query = query.Order(dbColumn + " " + order)
			} else {
				query = query.Order("created_at DESC")
			}
		} else {
			query = query.Order("created_at DESC")
		}
	} else {
		query = query.Order("created_at DESC")
	}

	// 5. 应用分页 (确保 query 被重新赋值)
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	// 6. 执行查询并预加载关联数据
	if err := query.Preload("Author").Preload("ReplyToUser").Find(&comments).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			currentPage := params.Page
			pageSize := params.Limit
			return &vo.GetCommentsResponseDataVO{Items: []vo.CommentVO{}, Total: 0, CurrentPage: &currentPage, PageSize: &pageSize}, nil
		}
		return nil, fmt.Errorf("获取评论列表失败: %w", err)
	}

	// 7. 将查询结果转换为VO
	itemsVO := make([]vo.CommentVO, 0, len(comments))
	for _, commentModel := range comments {
		commentVO, err := convertCommentToVO(commentModel, currentUserID, true, 2, 0) // 获取最多2层子评论
		if err == nil && commentVO != nil {
			itemsVO = append(itemsVO, *commentVO)
		} else if err != nil {
			log.Printf("转换评论 %d 到 VO 失败: %v", commentModel.ID, err)
		}
	}

	currentPage := params.Page
	pageSize := params.Limit
	return &vo.GetCommentsResponseDataVO{
		Items:       itemsVO,
		Total:       total,
		CurrentPage: &currentPage,
		PageSize:    &pageSize,
	}, nil
}

func (s *CommentService) AddComment(data *dto.AddCommentDTO, authorID uint32) (*vo.CommentVO, error) {
	// 1. 验证 PostID 是否有效
	var post dto.Post
	if err := database.Client.First(&post, data.PostID).Error; err != nil { // 使用 database.Client
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("关联的帖子未找到")
		}
		return nil, fmt.Errorf("验证帖子ID失败: %w", err)
	}

	// 2. 如果是回复，验证 ParentID 和 ReplyToUserID
	if data.ParentID != nil && *data.ParentID > 0 {
		var parentComment dto.Comment
		if err := database.Client.First(&parentComment, *data.ParentID).Error; err != nil { // 使用 database.Client
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("回复的父评论未找到")
			}
			return nil, fmt.Errorf("验证父评论ID失败: %w", err)
		}
		if parentComment.PostID != data.PostID {
			return nil, errors.New("父评论与当前帖子不匹配")
		}
	}
	if data.ReplyToUserID != nil && *data.ReplyToUserID > 0 {
		var replyToUser dto.User
		if err := database.Client.First(&replyToUser, *data.ReplyToUserID).Error; err != nil { // 使用 database.Client
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("回复的目标用户未找到")
			}
			return nil, fmt.Errorf("验证回复目标用户ID失败: %w", err)
		}
	}

	newComment := dto.Comment{
		PostID:        data.PostID,
		AuthorID:      authorID,
		Content:       data.Content,
		ParentID:      data.ParentID,
		ReplyToUserID: data.ReplyToUserID,
	}

	// 3. 使用事务创建评论
	tx := database.Client.Begin() // 使用 database.Client
	if tx.Error != nil {
		return nil, fmt.Errorf("开启事务失败: %w", tx.Error)
	}

	if err := tx.Create(&newComment).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("创建评论记录失败: %w", err)
	}

	// 可选：在同一事务中更新帖子的评论数 (示例)
	// if err := tx.Model(&dto.Post{}).Where("id = ?", data.PostID).UpdateColumn("comments_count", gorm.Expr("comments_count + 1")).Error; err != nil {
	//    tx.Rollback()
	//    return nil, fmt.Errorf("更新帖子评论数失败: %w", err)
	// }

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交评论事务失败: %w", err)
	}

	// 4. 创建成功后，重新加载评论以包含关联的Author和ReplyToUser信息
	var reloadedComment dto.Comment
	if err := database.Client.Preload("Author").Preload("ReplyToUser").First(&reloadedComment, newComment.ID).Error; err != nil { // 使用 database.Client
		log.Printf("警告: 评论 %d 创建成功，但重新加载关联信息失败: %v。将返回部分信息。", newComment.ID, err)
		// 即使重新加载失败，也尝试用已有的 newComment (无预加载) 进行转换
		return convertCommentToVO(newComment, &authorID, false, 0, 0)
	}

	return convertCommentToVO(reloadedComment, &authorID, false, 0, 0)
}

func (s *CommentService) DeleteComment(commentID uint32, userID uint32, userRole string) error {
	var comment dto.Comment
	if err := database.Client.First(&comment, commentID).Error; err != nil { // 使用 database.Client
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评论未找到")
		}
		return fmt.Errorf("查找待删除评论失败: %w", err)
	}

	tx := database.Client.Begin() // 使用 database.Client
	if tx.Error != nil {
		return fmt.Errorf("开启删除事务失败: %w", tx.Error)
	}

	// 1. 删除关联的点赞记录 (如果存在且需要手动处理)
	if err := tx.Where("comment_id = ?", commentID).Delete(&dto.UserCommentLike{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除评论关联的点赞记录失败: %w", err)
	}

	// 2. 处理子评论 (根据业务逻辑，可能需要级联删除或置空 parent_id)
	// 此处简化，不显式处理。

	// 3. 删除评论本身
	if err := tx.Delete(&dto.Comment{}, commentID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除评论记录失败: %w", err)
	}

	// 可选：更新帖子评论数
	// if err := tx.Model(&dto.Post{}).Where("id = ?", comment.PostID).UpdateColumn("comments_count", gorm.Expr("comments_count - 1")).Error; err != nil {
	//    tx.Rollback()
	//    return fmt.Errorf("更新帖子评论数（删除后）失败: %w", err)
	// }

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交删除评论事务失败: %w", err)
	}
	return nil
}

func (s *CommentService) ToggleLikeComment(commentID uint32, userID uint32) (*vo.ToggleLikeCommentResponseDataVO, error) {
	var comment dto.Comment
	if err := database.Client.First(&comment, commentID).Error; err != nil { // 使用 database.Client
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("评论未找到")
		}
		return nil, fmt.Errorf("查找评论失败: %w", err)
	}

	isLiked := false
	currentLikesCount := comment.LikesCount

	tx := database.Client.Begin() // 使用 database.Client
	if tx.Error != nil {
		return nil, fmt.Errorf("开启点赞事务失败: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil { // 处理事务中的 panic
			tx.Rollback()
		}
	}()

	var existingLike dto.UserCommentLike
	err := tx.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&existingLike).Error

	if errors.Is(err, gorm.ErrRecordNotFound) { // 未点赞 -> 点赞
		newLike := dto.UserCommentLike{UserID: userID, CommentID: commentID}
		if errCreate := tx.Create(&newLike).Error; errCreate != nil {
			tx.Rollback()
			return nil, fmt.Errorf("创建点赞记录失败: %w", errCreate)
		}
		if errUpdate := tx.Model(&dto.Comment{}).Where("id = ?", commentID).UpdateColumn("likes_count", gorm.Expr("likes_count + 1")).Error; errUpdate != nil {
			tx.Rollback()
			return nil, fmt.Errorf("增加评论点赞数失败: %w", errUpdate)
		}
		if err := tx.Model(&dto.Post{}).Where("id = ?", commentID).UpdateColumn("is_liked_by_current_user", gorm.Expr("true")).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		isLiked = true
		currentLikesCount++
	} else if err == nil { // 已点赞 -> 取消点赞
		if errDelete := tx.Delete(&existingLike).Error; errDelete != nil {
			tx.Rollback()
			return nil, fmt.Errorf("删除点赞记录失败: %w", errDelete)
		}
		if currentLikesCount > 0 {
			if errUpdate := tx.Model(&dto.Comment{}).Where("id = ?", commentID).UpdateColumn("likes_count", gorm.Expr("likes_count - 1")).Error; errUpdate != nil {
				tx.Rollback()
				return nil, fmt.Errorf("减少评论点赞数失败: %w", errUpdate)
			}
			if err := tx.Model(&dto.Post{}).Where("id = ?", commentID).UpdateColumn("is_liked_by_current_user", gorm.Expr("false")).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			currentLikesCount--
		}
		isLiked = false
	} else { // 查询点赞记录时发生其他错误
		tx.Rollback()
		return nil, fmt.Errorf("检查用户点赞状态失败: %w", err)
	}

	if errCommit := tx.Commit().Error; errCommit != nil {
		return nil, fmt.Errorf("提交点赞事务失败: %w", errCommit)
	}

	return &vo.ToggleLikeCommentResponseDataVO{
		IsLiked:    isLiked,
		LikesCount: int(currentLikesCount),
	}, nil
}
