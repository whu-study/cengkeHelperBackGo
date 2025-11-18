package dto

import "time"

type CourseInfo struct {
	ID         uint32 `gorm:"not null;uint;primaryKey;autoIncrement" json:"id"`
	Years      string `gorm:"not null;type:varchar(255)" json:"years"`     // 2021-2022
	Semester   string `gorm:"not null;type:varchar(255)" json:"semester"`  // 秋季学期
	CourseNum  string `gorm:"not null;type:varchar(255)" json:"courseNum"` // 2021-2022-1-1001
	CourseName string `gorm:"not null;type:varchar(255)" json:"courseName"`
	Faculty    string `gorm:"not null;type:varchar(255)" json:"faculty"`

	Credit string `gorm:"not null;type:varchar(255)" json:"credit"` // 2.0

	CourseComplexion string `gorm:"not null;type:varchar(255)" json:"courseComplexion"`
	CourseType       string `gorm:"not null;type:varchar(255)" json:"courseType"`
	Grade            string `gorm:"not null;type:varchar(255)" json:"grade"`
	Major            string `gorm:"not null;type:varchar(255)" json:"major"`
	Teacher          string `gorm:"not null;type:varchar(255)" json:"teacher"`
	TeacherTitle     string `gorm:"not null;type:varchar(255)" json:"teacherTitle"`

	Description string `gorm:"type:text" json:"description,omitempty"`

	AverageRating float32 `gorm:"default:0" json:"rating,omitempty"`
	ReviewCount   uint32  `gorm:"default:0" json:"reviewCount,omitempty"`

	Reviews []CourseReviewModel `gorm:"foreignKey:CourseID" json:"-"`
}

type TimeInfo struct {
	ID           uint32 `gorm:"not null;primaryKey;autoIncrement" json:"id"`
	CourseInfoId uint32 `gorm:"not null" json:"courseInfo"`

	WeekAndTime uint32 `gorm:"not null" json:"weekAndTime"`

	DayOfWeek uint8 `gorm:"not null" json:"dayOfWeek"` // 0-6

	Area      uint8  `gorm:"not null" json:"area"` // 1-4
	Building  string `gorm:"not null;type:varchar(255)" json:"building"`
	Classroom string `gorm:"not null;type:varchar(255)" json:"classroom"`
}

type CourseReviewModel struct {
	ID        uint32    `gorm:"primaryKey" json:"id"`
	CourseID  uint32    `gorm:"not null;index" json:"courseId"`      // 关联的课程ID
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
	CourseID uint32 `json:"courseId" binding:"required"` // 课程ID
	// 使用 float32 接收前端可能传来的分数（例如 4.5），后端会在保存前转换为整数以兼容现有存储结构
	Rating  float32 `json:"rating" binding:"required,gte=1,lte=5"` // 课程评分 (支持小数，如 4.5)
	Comment string  `json:"comment" binding:"required,max=1000"`   // 课程评论内容
	// UserID 会从当前登录用户获取，不需要前端传递
}
