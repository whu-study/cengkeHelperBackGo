package dto

import "time"

// UserPostLike 记录用户对帖子的点赞
// 表名将是 "user_post_likes" (GORM 默认)
type UserPostLike struct {
	UserID    uint32    `gorm:"primaryKey;not null;comment:点赞用户ID"`
	PostID    uint32    `gorm:"primaryKey;not null;comment:被点赞的帖子ID"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:点赞时间"`

	//// 定义外键关联 (可选，但推荐，GORM 自动迁移时会尝试创建)
	//User User `gorm:"foreignKey:UserID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	//Post Post `gorm:"foreignKey:PostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName 自定义表名 (如果 GORM 默认不符合您的期望)
// func (UserPostLike) TableName() string {
//  return "user_post_likes"
// }

// UserPostCollect 记录用户对帖子的收藏
// 表名将是 "user_post_favorites"
type UserPostCollect struct {
	UserID    uint32    `gorm:"primaryKey;not null;comment:收藏用户ID"`
	PostID    uint32    `gorm:"primaryKey;not null;comment:被收藏的帖子ID"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:收藏时间"`

	//// 可选的外键定义
	//User User `gorm:"foreignKey:UserID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	//Post Post `gorm:"foreignKey:PostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName 自定义表名
// func (UserPostFavorite) TableName() string {
//  return "user_post_favorites"
// }

// UserCommentLike 记录用户对评论的点赞
// 表名将是 "user_comment_likes"
type UserCommentLike struct {
	UserID    uint32    `gorm:"primaryKey;not null;comment:点赞用户ID"`
	CommentID uint32    `gorm:"primaryKey;not null;comment:被点赞的评论ID"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:点赞时间"`

	//// 可选的外键定义
	//User    User    `gorm:"foreignKey:UserID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	//Comment Comment `gorm:"foreignKey:CommentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName 自定义表名
// func (UserCommentLike) TableName() string {
//  return "user_comment_likes"
// }
