package dto

import "time"

// GetCommentsParamsDTO 对应前端 GetCommentsParams，用于获取评论列表的查询参数
type GetCommentsParamsDTO struct {
	Page   int    `form:"page,default=1"`
	Limit  int    `form:"limit,default=10"` // 默认每页10条评论
	SortBy string `form:"sortBy,omitempty"` // 例如 "createdAt_asc" 或 "likesCount_desc"
}

// AddCommentDTO 对应前端 AddCommentPayload，用于添加新评论的请求体
type AddCommentDTO struct {
	PostID        uint32  `json:"postId" binding:"required"`
	Content       string  `json:"content" binding:"required,min=1,max=1000"`
	ParentID      *uint32 `json:"parentId,omitempty"`      // 指针表示可选，用于回复
	ReplyToUserID *uint32 `json:"replyToUserId,omitempty"` // 指针表示可选，回复的目标用户ID
}

// Comment 对应数据库中的 'comments' 表
type Comment struct {
	ID                   uint32    `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID               uint32    `gorm:"not null;index;comment:评论所属的帖子ID" json:"postId"` // 直接对应 DTO 的 postId
	Post                 Post      `gorm:"foreignKey:PostID" json:"-"`                     // 关联帖子，通常不在JSON中完整返回避免循环
	AuthorID             uint32    `gorm:"not null;index;comment:评论作者的用户ID" json:"authorId"`
	Author               User      `gorm:"foreignKey:AuthorID" json:"author"` // 关联作者信息
	Content              string    `gorm:"type:text;not null;comment:评论内容" json:"content"`
	CreatedAt            time.Time `gorm:"autoCreateTime;index" json:"createdAt"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime" json:"updatedAt"`           // 虽然前端没直接显示，但通常会有
	IsLikedByCurrentUser bool      `gorm:"default:false" json:"isLikedByCurrentUser"` // 标记当前用户是否已点赞
	LikesCount           uint      `gorm:"default:0;comment:评论点赞数量" json:"likesCount"`

	// --- 用于支持回复功能 ---
	ParentID      *uint32  `gorm:"index;comment:父评论ID (用于回复)" json:"parentId,omitempty"`   // 指针表示可选
	ParentComment *Comment `gorm:"foreignKey:ParentID" json:"-"`                           // 关联父评论，避免循环JSON
	ReplyToUserID *uint32  `gorm:"index;comment:回复的目标用户ID" json:"replyToUserId,omitempty"` // 指针表示可选
	ReplyToUser   *User    `gorm:"foreignKey:ReplyToUserID" json:"replyToUser,omitempty"`  // 关联被回复的用户信息

	Children []*Comment `gorm:"foreignKey:ParentID" json:"children,omitempty"` // 嵌套的子回复列表

	// IsDeleted bool `gorm:"default:false;comment:是否已删除 (用于软删除)" json:"-"` // 软删除标记
}

// TableName 自定义 Comment 模型对应的表名
func (Comment) TableName() string {
	return "comments"
}
