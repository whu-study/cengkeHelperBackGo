package dto

import (
	"gorm.io/datatypes" // 导入 GORM 的 datatypes 包
	"time"
)

// --- DTOs ---

// GetPostsParamsDTO 对应前端 GetPostsParams，用于获取帖子列表的查询参数
type GetPostsParamsDTO struct {
	Page       int    `form:"page,default=1"` // gin 中用 form tag 接收 query 参数
	Limit      int    `form:"limit,default=10"`
	SortBy     string `form:"sortBy,omitempty"`     // 例如 "createdAt_desc"
	FilterText string `form:"filterText,omitempty"` // 搜索关键词
	Category   string `form:"category,omitempty"`
	Tag        string `form:"tag,omitempty"`
	AuthorID   uint32 `form:"authorId,omitempty"`
}

// CreatePostDTO 对应前端 CreatePostBody，用于创建新帖子的请求体
// 前端 CreatePostBody Omit 了很多字段，这里我们定义实际需要创建的字段
type CreatePostDTO struct {
	Title    string   `json:"title" binding:"required,min=5,max=100"`
	Content  string   `json:"content" binding:"required,min=20"`
	Tags     []string `json:"tags,omitempty"`
	Category *string  `json:"category,omitempty"` // 使用指针表示可选
	// AuthorID uint32 `json:"authorId"` // 通常由后端从JWT获取，不由前端传递
}

// UpdatePostDTO 对应前端 UpdatePostBody，用于更新帖子的请求体
type UpdatePostDTO struct {
	Title    *string  `json:"title,omitempty" binding:"omitempty,min=5,max=100"` // 指针表示可选更新
	Content  *string  `json:"content,omitempty" binding:"omitempty,min=20"`
	Tags     []string `json:"tags,omitempty"` // 如果传空数组表示清空，如果 omitempty 且为nil则不更新
	Category *string  `json:"category,omitempty"`
	// isPublished, isPinned, isLocked 等状态的更新也可以放在这里
}

// Post 对应数据库中的 'posts' 表
type Post struct {
	ID                       uint32         `gorm:"primaryKey;autoIncrement" json:"id"`
	Title                    string         `gorm:"type:varchar(255);not null;comment:帖子标题" json:"title"`
	Content                  string         `gorm:"type:text;not null;comment:帖子内容 (HTML 或 Markdown)" json:"content"`
	AuthorID                 uint32         `gorm:"not null;type:uint;index;comment:帖子作者的用户ID" json:"authorId"`
	Author                   User           `gorm:"foreignKey:AuthorID" json:"author"` // 关联作者信息
	CreatedAt                time.Time      `gorm:"autoCreateTime;index" json:"createdAt"`
	UpdatedAt                time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	Category                 *string        `gorm:"type:varchar(100);index;comment:帖子分类" json:"category,omitempty"`
	Tags                     datatypes.JSON `gorm:"type:json;comment:帖子标签数组" json:"tags,omitempty"` // 存储为 JSON 字符串或 JSONB
	ViewCount                uint           `gorm:"default:0;comment:帖子浏览次数" json:"viewCount"`
	LikesCount               uint           `gorm:"default:0;comment:帖子点赞数量" json:"likesCount"`
	CollectCount             uint           `gorm:"default:0;comment:帖子收藏数量" json:"collectCount"`
	CommentsCount            uint           `gorm:"default:0;comment:帖子评论数量" json:"commentsCount"`
	IsPublished              bool           `gorm:"default:true;comment:是否已发布" json:"isPublished"`
	IsPinned                 bool           `gorm:"default:false;comment:是否置顶" json:"isPinned"`
	IsLocked                 bool           `gorm:"default:false;comment:是否锁定评论" json:"isLocked"`
	IsLikedByCurrentUser     bool           `gorm:"default:false;comment:当前用户是否点赞" json:"isLikedByCurrentUser"`
	IsCollectedByCurrentUser bool           `gorm:"default:false;comment:当前用户是否收藏" json:"IsCollectedByCurrentUser"`
	// 如果需要追踪最后评论信息，可以添加以下字段，但通常这些可以通过查询动态获取或在评论创建时更新
	// LastCommentAt      *time.Time `gorm:"comment:最后评论时间" json:"lastCommentAt,omitempty"`
	// LastCommentUserID *uint      `gorm:"comment:最后评论的用户ID" json:"lastCommentUserId,omitempty"`
	// LastCommentUser   User       `gorm:"foreignKey:LastCommentUserID" json:"lastCommentUser,omitempty"`

	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"` // 帖子的评论列表
}

// TableName 自定义 Post 模型对应的表名
func (Post) TableName() string {
	return "posts"
}
