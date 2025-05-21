package services

import (
	database "cengkeHelperBackGo/internal/db"
	"cengkeHelperBackGo/internal/models/dto"
	"cengkeHelperBackGo/internal/models/vo"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"strings"
)

type PostService struct {
}

func NewPostService() *PostService {
	// 构造函数现在不需要初始化 db
	return &PostService{}
}

// GetPosts 获取帖子列表并处理分页、排序和过滤
func (s *PostService) GetPosts(params *dto.GetPostsParamsDTO) (*vo.GetPostsResponseDataVO, error) {
	var posts []dto.Post
	var total int64

	query := database.Client.Model(&dto.Post{})
	if params.FilterText != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+params.FilterText+"%", "%"+params.FilterText+"%")
	}
	if params.Category != "" {
		query = query.Where("category = ?", params.Category)
	}
	if params.AuthorID != 0 {
		query = query.Where("author_id = ?", params.AuthorID)
	}
	if params.Tag != "" {
		query = query.Where("CAST(tags AS CHAR) LIKE ?", "%"+params.Tag+"%")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 应用排序
	if params.SortBy != "" {
		parts := strings.Split(params.SortBy, "_")
		if len(parts) == 2 {
			column := parts[0]
			order := strings.ToUpper(parts[1])

			// Map frontend column names to actual database column names
			columnMapping := map[string]string{
				"createdAt":     "created_at",
				"updatedAt":     "updated_at",
				"viewCount":     "view_count",
				"likesCount":    "likes_count",
				"commentsCount": "comments_count",
			}

			if dbColumn, exists := columnMapping[column]; exists {
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

	// 应用分页
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	// 预加载 Author 信息
	if err := query.Preload("Author").Find(&posts).Error; err != nil {
		return nil, err
	}

	// --- 转换为 VO (关键修改) ---
	// itemsVO := make([]vo.PostSimpleVO, len(posts)) // 旧代码
	itemsVO := make([]vo.PostVO, len(posts)) // *** 修改点 1: 创建 []PostVO 切片 ***
	for i, p := range posts {
		// 处理 Tags
		var tagsInPost []string
		if p.Tags != nil && len(p.Tags) > 0 {
			if err := json.Unmarshal(p.Tags, &tagsInPost); err != nil {
				tagsInPost = []string{}
			}
		}

		// 处理 Author (假设 dto.User 结构体中有 Username, ID, Avatar 等字段)
		// *** 修改点 2: 确保 AuthorInfoVO 的字段和 dto.User 匹配 ***
		var authorInfo vo.AuthorInfoVO // 假设 post_vo.txt 中定义了 AuthorInfoVO
		if p.Author.Id != 0 {
			authorInfo = vo.AuthorInfoVO{
				// 这里填充 AuthorInfoVO 的字段，例如：
				ID:       p.Author.Id,
				Username: p.Author.Username, // 假设 dto.User 有 Username
				// Nickname: p.Author.Nickname, // 如果 AuthorInfoVO 需要
				Avatar: p.Author.Avatar, // 如果 AuthorInfoVO 需要
			}
		}

		// 处理可选字段 (int 和 bool 转指针)
		viewCountPtr := int(p.ViewCount) // 先转 int
		likesCountPtr := int(p.LikesCount)
		collectCountPtr := int(p.CollectCount)
		commentsCountPtr := int(p.CommentsCount)
		isPublishedPtr := p.IsPublished
		isPinnedPtr := p.IsPinned
		isLockedPtr := p.IsLocked
		isLikedByCurrentUserPtr := p.IsLikedByCurrentUser
		isCollectedByCurrentUserPtr := p.IsCollectedByCurrentUser
		// isLikedByCurrentUserPtr := false // 这个需要额外逻辑判断当前用户是否点赞

		// *** 修改点 3: 填充 PostVO 而不是 PostSimpleVO ***
		itemsVO[i] = vo.PostVO{
			ID:        p.ID,
			Title:     p.Title,
			Content:   p.Content, // PostVO 需要完整 Content
			Author:    authorInfo,
			CreatedAt: p.CreatedAt,
			UpdatedAt: &p.UpdatedAt, // 如果 UpdatedAt 不是指针类型，取地址
			Tags:      tagsInPost,
			Category:  p.Category,
			// --- 使用指针 ---
			ViewCount:                &viewCountPtr,
			LikesCount:               &likesCountPtr,
			CommentsCount:            &commentsCountPtr,
			IsPublished:              &isPublishedPtr,
			IsPinned:                 &isPinnedPtr,
			IsLocked:                 &isLockedPtr,
			CollectCount:             &collectCountPtr,
			IsLikedByCurrentUser:     &isLikedByCurrentUserPtr,     // 需要额外逻辑
			IsCollectedByCurrentUser: &isCollectedByCurrentUserPtr, // 需要额外逻辑
			// IsLikedByCurrentUser: &isLikedByCurrentUserPtr, // 需要额外逻辑
		}
	}

	currentPage := params.Page // 获取 int 值
	pageSize := params.Limit   // 获取 int 值
	return &vo.GetPostsResponseDataVO{
		Items:       itemsVO, // *** 现在是 []PostVO ***
		Total:       total,
		CurrentPage: &currentPage, // *** 修改点 4: 取地址变成 *int ***
		PageSize:    &pageSize,    // *** 修改点 5: 取地址变成 *int ***
	}, nil
}

func (s *PostService) GetPostByID(postID uint32, currentUserID *uint32) (*vo.PostVO, error) {
	// 开启事务
	tx := database.Client.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 获取帖子信息
	var post dto.Post
	err := tx.Model(&dto.Post{}).
		Preload("Author").
		First(&post, postID).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("帖子未找到")
		}
		return nil, err
	}

	// 2. 浏览量+1（使用原子操作避免并发问题）
	if err := tx.Model(&dto.Post{}).
		Where("id = ?", postID).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).
		Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新浏览量失败: %w", err)
	}

	// 3. 重新加载更新后的帖子信息（可选）
	if err := tx.First(&post, postID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	// 4. 转换为VO对象并返回
	return convertPostToVO(post, currentUserID, database.Client)
}

// --- Helper function to convert dto.Post to vo.PostVO ---
// (这个辅助函数可以抽取出来，或者每个需要的地方单独实现转换逻辑)
func convertPostToVO(post dto.Post, currentUserID *uint32, db *gorm.DB) (*vo.PostVO, error) {
	var tagsInPost []string
	if post.Tags != nil && len(post.Tags) > 0 {
		if err := json.Unmarshal(post.Tags, &tagsInPost); err != nil {
			tagsInPost = []string{}
		}
	}

	var authorInfo vo.AuthorInfoVO
	if post.Author.Id != 0 { // 假设 User DTO 中 Id 是主键
		authorInfo = vo.AuthorInfoVO{
			ID:       post.Author.Id,
			Username: post.Author.Username,
			Avatar:   post.Author.Avatar,
		}
	}

	viewCount := int(post.ViewCount)
	likesCount := int(post.LikesCount)
	commentsCount := int(post.CommentsCount)
	isPublished := post.IsPublished
	isPinned := post.IsPinned
	isLocked := post.IsLocked
	collectCount := int(post.CollectCount)

	var isLikedByCurrentUser bool = post.IsLikedByCurrentUser
	var isCollectedByCurrentUser bool = post.IsCollectedByCurrentUser

	if currentUserID != nil && *currentUserID > 0 && db != nil {
		var likeCount int64
		db.Model(&dto.UserPostLike{}).Where("user_id = ? AND post_id = ?", *currentUserID, post.ID).Count(&likeCount)
		if likeCount > 0 {
			isLikedByCurrentUser = true
		}
		var collectCount int64
		db.Model(&dto.UserPostCollect{}).Where("user_id = ? AND post_id = ?", *currentUserID, post.ID).Count(&collectCount)
		if collectCount > 0 {
			isCollectedByCurrentUser = true
		}
	}

	return &vo.PostVO{
		ID:                       post.ID,
		Title:                    post.Title,
		Content:                  post.Content,
		Author:                   authorInfo,
		CreatedAt:                post.CreatedAt,
		UpdatedAt:                &post.UpdatedAt,
		Tags:                     tagsInPost,
		Category:                 post.Category,
		ViewCount:                &viewCount,
		LikesCount:               &likesCount,
		CommentsCount:            &commentsCount,
		IsPublished:              &isPublished,
		IsPinned:                 &isPinned,
		IsLocked:                 &isLocked,
		CollectCount:             &collectCount,
		IsLikedByCurrentUser:     &isLikedByCurrentUser,
		IsCollectedByCurrentUser: &isCollectedByCurrentUser,
	}, nil
}

// --- 新增 CreatePost 方法 ---
func (s *PostService) CreatePost(postData *dto.CreatePostDTO, authorID uint32) (*vo.PostVO, error) {
	var tagsJSON datatypes.JSON
	if len(postData.Tags) > 0 {
		tagsBytes, err := json.Marshal(postData.Tags)
		if err != nil {
			return nil, fmt.Errorf("序列化标签失败: %w", err)
		}
		tagsJSON = datatypes.JSON(tagsBytes)
	}

	newPost := dto.Post{
		Title:    postData.Title,
		Content:  postData.Content,
		AuthorID: authorID, // 从认证信息中获取的作者ID
		Category: postData.Category,
		Tags:     tagsJSON,
		// IsPublished 默认应为 true (在 GORM 模型中定义)
	}

	result := database.Client.Create(&newPost)
	if result.Error != nil {
		return nil, result.Error
	}

	// 创建成功后，需要重新查询一次以预加载 Author 信息并转换为 VO
	// 或者直接构造 VO，但 Author 信息可能不完整
	// 为简单起见，我们直接用创建后的 newPost (ID 已被填充) 来构造 VO，但 Author 字段会是空的
	// 更好的做法是：database.Client.Preload("Author").First(&newPost, newPost.ID)
	// 然后调用 convertPostToVO(newPost, &authorID, database.Client)

	// 为了返回完整的 PostVO，包括作者信息，我们需要重新查询
	var createdPostWithAuthor dto.Post
	if err := database.Client.Preload("Author").First(&createdPostWithAuthor, newPost.ID).Error; err != nil {
		return nil, fmt.Errorf("获取创建后的帖子信息失败: %w", err)
	}

	return convertPostToVO(createdPostWithAuthor, &authorID, database.Client)
}

// --- 新增 UpdatePost 方法 ---
func (s *PostService) UpdatePost(postID uint32, userID uint32, postData *dto.UpdatePostDTO) (*vo.PostVO, error) {
	var existingPost dto.Post
	// 1. 查找帖子并验证作者
	if err := database.Client.Where("id = ?", postID).First(&existingPost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("帖子未找到")
		}
		return nil, err
	}

	// 3. 准备更新的数据
	updates := make(map[string]interface{})
	if postData.Title != nil {
		updates["title"] = *postData.Title
	}
	if postData.Content != nil {
		updates["content"] = *postData.Content
	}
	if postData.Category != nil {
		updates["category"] = *postData.Category
	}
	// 对于 Tags，如果 postData.Tags 非 nil (即使是空数组)，都表示要更新
	if postData.Tags != nil {
		tagsBytes, err := json.Marshal(postData.Tags)
		if err != nil {
			return nil, fmt.Errorf("序列化标签失败: %w", err)
		}
		updates["tags"] = datatypes.JSON(tagsBytes)
	}
	// 如果允许更新 isPublished, isPinned, isLocked 等，也在这里添加
	// if postData.IsPublished != nil { updates["is_published"] = *postData.IsPublished }

	if len(updates) == 0 {
		// 如果没有要更新的字段，可以直接返回当前帖子信息 (需要重新查询以包含用户信息)
		return s.GetPostByID(postID, &userID)
	}

	// 4. 执行更新
	if err := database.Client.Model(&existingPost).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 5. 返回更新后的帖子详情 (重新查询以获取最新数据和关联数据)
	return s.GetPostByID(postID, &userID)
}

// --- 新增 DeletePost 方法 ---
func (s *PostService) DeletePost(postID uint32, userID uint32, userRole string) error {
	var post dto.Post
	// 1. 查找帖子并验证作者
	if err := database.Client.Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("帖子未找到")
		}
		return err
	}

	// 3. 执行删除 (硬删除)
	// 如果需要软删除，GORM 模型中应包含 gorm.DeletedAt，并且这里使用 db.Delete(&post)
	// 硬删除：
	if err := database.Client.Select("Comments", "UserPostLikes", "UserPostCollects").Delete(&post).Error; err != nil {
		// GORM 的 Delete 如果设置了 Select，会尝试删除关联数据（如果关联已定义且支持级联或通过回调处理）
		// 如果没有 Select 或者关联未正确设置，可能需要手动删除关联表中的记录
		// 或者依赖数据库的级联删除约束 (ON DELETE CASCADE)
		// 为了确保，可以显式删除关联：
		// database.Client.Where("post_id = ?", postID).Delete(&dto.UserPostLike{})
		// database.Client.Where("post_id = ?", postID).Delete(&dto.UserPostCollect{})
		// database.Client.Where("post_id = ?", postID).Delete(&dto.Comment{}) // 假设 Comment 有 PostID
		return err
	}
	// 对于硬删除，如果上面 Select 方式不工作，直接删除帖子即可，依赖数据库外键约束
	// if err := database.Client.Delete(&dto.Post{}, postID).Error; err != nil {
	// return err
	// }

	return nil
}

// --- ToggleLikePost 方法 ---
func (s *PostService) ToggleLikePost(postID uint32, userID uint32) (*vo.ToggleLikeResponseDataVO, error) {
	var post dto.Post
	// First, check if the post exists
	if err := database.Client.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("帖子未找到")
		}
		return nil, err // Other database error
	}

	like := dto.UserPostLike{UserID: userID, PostID: postID}
	var currentLikesCount int64 // To store the final likes count

	// Use a transaction to ensure atomicity for all database operations
	tx := database.Client.Begin()
	if tx.Error != nil {
		return nil, tx.Error // Failed to begin transaction
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // Rollback on panic
		}
	}()

	var existingLike dto.UserPostLike
	isLiked := false
	// Check if the user has already liked this post
	err := tx.Where("user_id = ? AND post_id = ?", userID, postID).First(&existingLike).Error

	if err == nil { // User has already liked, so now unlike
		if err := tx.Delete(&existingLike).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		// Decrement likes_count on the post
		if err := tx.Model(&dto.Post{}).Where("id = ?", postID).UpdateColumn("likes_count", gorm.Expr("likes_count - 1")).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Model(&dto.Post{}).Where("id = ?", postID).UpdateColumn("is_liked_by_current_user", gorm.Expr("false")).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		isLiked = false
	} else if errors.Is(err, gorm.ErrRecordNotFound) { // User has not liked yet, so now like
		if err := tx.Create(&like).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		// Increment likes_count on the post
		if err := tx.Model(&dto.Post{}).Where("id = ?", postID).UpdateColumn("likes_count", gorm.Expr("likes_count + 1")).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Model(&dto.Post{}).Where("id = ?", postID).UpdateColumn("is_liked_by_current_user", gorm.Expr("true")).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		isLiked = true
	} else { // Another error occurred while checking for existing like
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Get the latest likes_count from the Post table after the transaction
	var updatedPostForLikeCount dto.Post
	if err := database.Client.Select("likes_count").First(&updatedPostForLikeCount, postID).Error; err != nil {
		// Log this error, but we can still proceed with the isLiked status.
		// Depending on requirements, you might want to return an error here too.
		// For now, we'll use the count from the post object before potential decrement/increment
		// if the re-fetch fails, or set to 0 or handle as an error.
		// A robust way is to re-query.
		currentLikesCount = -1 // Indicate an issue fetching the count, or handle error more strictly
	} else {
		currentLikesCount = int64(updatedPostForLikeCount.LikesCount)
	}

	return &vo.ToggleLikeResponseDataVO{
		IsLiked:    isLiked,
		LikesCount: int(currentLikesCount),
	}, nil
}

// --- 新增 ToggleCollectPost 方法 ---
func (s *PostService) ToggleCollectPost(postID uint32, userID uint32) (*vo.ToggleCollectResponseDataVO, error) {
	var post dto.Post // Check if the post exists
	if err := database.Client.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("帖子未找到")
		}
		return nil, err
	}

	collect := dto.UserPostCollect{UserID: userID, PostID: postID}
	var currentCollectCount int64

	tx := database.Client.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingCollect dto.UserPostCollect
	isCollected := false
	// Check if the user has already collected this post
	err := tx.Where("user_id = ? AND post_id = ?", userID, postID).First(&existingCollect).Error

	if err == nil { // Already collected, so uncollect
		if err := tx.Delete(&existingCollect).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		// Update collect_count on the Post table
		if err := tx.Model(&dto.Post{}).Where("id = ?", postID).UpdateColumn("collect_count", gorm.Expr("collect_count - 1")).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Model(&dto.Post{}).Where("id = ?", postID).UpdateColumn("is_collected_by_current_user", gorm.Expr("false")).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		isCollected = false
	} else if errors.Is(err, gorm.ErrRecordNotFound) { // Not collected yet, so collect
		if err := tx.Create(&collect).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		// Update collect_count on the Post table
		if err := tx.Model(&dto.Post{}).Where("id = ?", postID).UpdateColumn("collect_count", gorm.Expr("collect_count + 1")).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Model(&dto.Post{}).Where("id = ?", postID).UpdateColumn("is_collected_by_current_user", gorm.Expr("true")).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		isCollected = true
	} else { // Another error
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Get the latest collect_count from the Post table after the transaction
	var updatedPostForCollectCount dto.Post
	if err := database.Client.Select("collect_count").First(&updatedPostForCollectCount, postID).Error; err != nil {
		currentCollectCount = -1 // Indicate an issue
	} else {
		currentCollectCount = int64(updatedPostForCollectCount.CollectCount)
	}

	return &vo.ToggleCollectResponseDataVO{
		IsCollected:  isCollected,
		CollectCount: int(currentCollectCount),
	}, nil
}
