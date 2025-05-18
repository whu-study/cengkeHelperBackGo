package vo

import (
	"time"
)

// PostSimpleVO 用于帖子列表中的单个帖子项 (简化版，不含完整 Author 对象，按需调整)
type PostSimpleVO struct {
	ID                       uint32       `json:"id"`
	Title                    string       `json:"title"`
	Author                   UserSimpleVO `json:"author"` // 嵌套简化的作者信息
	CreatedAt                time.Time    `json:"createdAt"`
	UpdatedAt                time.Time    `json:"updatedAt"`
	Category                 *string      `json:"category,omitempty"`
	Tags                     []string     `json:"tags,omitempty"` // 反序列化后的 tags
	ViewCount                uint         `json:"viewCount"`
	LikesCount               uint         `json:"likesCount"`
	CommentsCount            uint         `json:"commentsCount"`
	IsPublished              bool         `json:"isPublished"`
	IsPinned                 bool         `json:"isPinned"`
	IsLocked                 bool         `json:"isLocked"`
	ContentExcerpt           string       `json:"contentExcerpt,omitempty"` // 可选：内容摘要
	IsCollectedByCurrentUser bool         `json:"isCollectedByCurrentUser"`
	IsLikedByCurrentUser     bool         `json:"isLikedByCurrentUser"`
	// 根据前端 Post 类型，可能还需要 'content' 字段。
	// 如果列表项不需要完整content，则可以只在详情接口返回。
	// content: string;
}

// PostDetailVO 用于帖子详情 (可以更完整)
type PostDetailVO struct {
	ID                       uint32       `json:"id"`
	Title                    string       `json:"title"`
	Content                  string       `json:"content"`
	Author                   UserSimpleVO `json:"author"`
	CreatedAt                time.Time    `json:"createdAt"`
	UpdatedAt                time.Time    `json:"updatedAt"`
	Category                 *string      `json:"category,omitempty"`
	Tags                     []string     `json:"tags,omitempty"`
	ViewCount                uint         `json:"viewCount"`
	LikesCount               uint         `json:"likesCount"`
	CollectCount             uint         `json:"collectCount"`
	CommentsCount            uint         `json:"commentsCount"`
	IsPublished              bool         `json:"isPublished"`
	IsPinned                 bool         `json:"isPinned"`
	IsLocked                 bool         `json:"isLocked"`
	Comments                 []CommentVO  `json:"comments,omitempty"` // 帖子下的评论
	IsCollectedByCurrentUser bool         `json:"isCollectedByCurrentUser"`
	IsLikedByCurrentUser     bool         `json:"isLikedByCurrentUser"`
}

// UserSimpleVO 用于内嵌在 PostVO 中的作者信息
type UserSimpleVO struct {
	ID       uint32 `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar,omitempty"`
}

// PostVO 对应前端的 Post 类型，用于API响应
type PostVO struct {
	ID                       uint32       `json:"id"`
	Title                    string       `json:"title"`
	Content                  string       `json:"content"` // 后端通常会做XSS清理
	Author                   AuthorInfoVO `json:"author"`
	CreatedAt                time.Time    `json:"createdAt"`
	UpdatedAt                *time.Time   `json:"updatedAt,omitempty"` // 指针表示可选
	Tags                     []string     `json:"tags,omitempty"`
	Category                 *string      `json:"category,omitempty"`
	ViewCount                *int         `json:"viewCount,omitempty"`
	LikesCount               *int         `json:"likesCount,omitempty"`
	CollectCount             *int         `json:"collectCount,omitempty"`
	CommentsCount            *int         `json:"commentsCount,omitempty"`
	IsPublished              *bool        `json:"isPublished,omitempty"`
	IsPinned                 *bool        `json:"isPinned,omitempty"`
	IsLocked                 *bool        `json:"isLocked,omitempty"`
	IsLikedByCurrentUser     *bool        `json:"isLikedByCurrentUser,omitempty"`     // 当前用户是否点赞
	IsCollectedByCurrentUser *bool        `json:"isCollectedByCurrentUser,omitempty"` // 当前用户是否收藏
}

// GetPostsResponseDataVO 对应前端 GetPostsResponseData
type GetPostsResponseDataVO struct {
	Items       []PostVO `json:"items"`
	Total       int64    `json:"total"` // 通常总数用 int64
	CurrentPage *int     `json:"currentPage,omitempty"`
	PageSize    *int     `json:"pageSize,omitempty"`
}

// ToggleLikeResponseDataVO 对应前端 ToggleLikeResponseData
type ToggleLikeResponseDataVO struct {
	IsLiked    bool `json:"isLiked"`
	LikesCount int  `json:"likesCount"`
}

// ToggleCollectResponseDataVO 对应前端 ToggleCollectResponseData
type ToggleCollectResponseDataVO struct {
	IsCollected  bool `json:"isCollected"`
	CollectCount int  `json:"collectCount"` // 注意类型
}
