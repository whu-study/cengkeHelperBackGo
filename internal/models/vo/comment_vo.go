package vo

import "time"

// CommentVO 对应前端的 Comment 类型，用于API响应
type CommentVO struct {
	ID                   uint32        `json:"id"`
	PostID               uint32        `json:"postId"` // 明确评论所属的帖子
	Author               AuthorInfoVO  `json:"author"`
	Content              string        `json:"content"`
	CreatedAt            time.Time     `json:"createdAt"`
	UpdatedAt            *time.Time    `json:"updatedAt,omitempty"`
	LikesCount           *int          `json:"likesCount,omitempty"`
	ParentID             *uint32       `json:"parentId,omitempty"`             // 如果是回复，父评论ID
	ReplyToUser          *AuthorInfoVO `json:"replyToUser,omitempty"`          // 如果是回复，被回复的用户信息
	Children             []CommentVO   `json:"children,omitempty"`             // 嵌套的子回复
	IsLikedByCurrentUser *bool         `json:"isLikedByCurrentUser,omitempty"` // 当前用户是否点赞此评论
}

// GetCommentsResponseDataVO 对应前端 GetCommentsResponseData
type GetCommentsResponseDataVO struct {
	Items       []CommentVO `json:"items"`
	Total       int64       `json:"total"` // 顶级评论的总数
	CurrentPage *int        `json:"currentPage,omitempty"`
	PageSize    *int        `json:"pageSize,omitempty"`
}

// ToggleLikeCommentResponseDataVO 对应前端 ToggleLikeCommentResponseData
type ToggleLikeCommentResponseDataVO struct {
	IsLiked    bool `json:"isLiked"`
	LikesCount int  `json:"likesCount"`
}
