package dto

import (
	//"cengkeHelperBackGo/internal/models/vo"
	"gorm.io/gorm"
	"time"
)

// 课程信息表
// 课程信息表
type CourseInfo struct {
	ID            uint                `gorm:"primaryKey" json:"id"`
	Room          string              `gorm:"type:varchar(100)" json:"room"`
	Faculty       string              `gorm:"type:varchar(100)" json:"faculty"`
	CourseName    string              `gorm:"type:varchar(255);not null" json:"courseName"`
	TeacherName   string              `gorm:"type:varchar(100)" json:"teacherName"`
	TeacherTitle  string              `gorm:"type:varchar(50)" json:"teacherTitle"`
	CourseTime    string              `gorm:"type:varchar(100)" json:"courseTime"`
	CourseType    string              `gorm:"type:varchar(50)" json:"courseType"`
	Description   string              `gorm:"type:text" json:"description,omitempty"`
	Credits       float32             `gorm:"type:decimal(3,1)" json:"credits,omitempty"`
	Building      string              `gorm:"type:varchar(100);index" json:"building"` // 教学楼名称，可以保留用于显示或冗余
	BuildingID    uint                `gorm:"index;comment:所属教学楼ID" json:"buildingId"` // <--- 新增的外键字段
	AverageRating float32             `gorm:"default:0" json:"rating,omitempty"`
	ReviewCount   uint                `gorm:"default:0" json:"reviewCount,omitempty"`
	CreatedAt     time.Time           `gorm:"autoCreateTime" json:"-"`
	UpdatedAt     time.Time           `gorm:"autoUpdateTime" json:"-"`
	DeletedAt     gorm.DeletedAt      `gorm:"index" json:"-"`
	Reviews       []CourseReviewModel `gorm:"foreignKey:CourseID" json:"-"`
}

type CourseReviewModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CourseID  uint      `gorm:"not null;index" json:"courseId"`      // 关联的课程ID
	UserID    uint32    `gorm:"not null;index" json:"userId"`        // 评价用户ID (关联 User 模型)
	Rating    int       `gorm:"type:tinyint;not null" json:"rating"` // 评分 (例如 1-5)
	Comment   string    `gorm:"type:text" json:"comment"`            // 评论内容
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`

	Course CourseInfo `gorm:"foreignKey:CourseID" json:"-"` // 反向关联，方便查询，但通常不在 JSON 中序列化
	User   User       `gorm:"foreignKey:UserID" json:"-"`   // 假设有一个 UserModel
	// 如果需要匿名评价或展示用户名，可以冗余存储或 JOIN 查询
	// ReviewerName string `gorm:"type:varchar(100)" json:"reviewerName,omitempty"` // 冗余的评价者名称，可选
}

// TableName 自定义表名
func (CourseReviewModel) TableName() string {
	return "course_reviews"
}

type CourseReviewCreateDTO struct {
	CourseID uint   `json:"courseId" binding:"required"`           // 课程ID
	Rating   int    `json:"rating" binding:"required,min=1,max=5"` // 课程评分 (例如 1-5)
	Comment  string `json:"comment" binding:"required,max=1000"`   // 课程评论内容
	// UserID 会从当前登录用户获取，不需要前端传递
}
