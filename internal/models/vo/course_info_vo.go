package vo

import "time"

//import "cengkeHelperBackGo/internal/models/dto"

// CourseInfoVO 课程信息VO（完全匹配前端接口）

// CourseInfoVO 对应前端的 CourseInfo 接口 (列表项)
type CourseInfoVO struct {
	ID           uint32 `json:"id"`
	Room         string `json:"room"`
	Faculty      string `json:"faculty"`
	CourseName   string `json:"courseName"`
	TeacherName  string `json:"teacherName"`
	TeacherTitle string `json:"teacherTitle"`
	CourseTime   string `json:"courseTime"`
	CourseType   string `json:"courseType"`
}

// CourseDetailVO 对应前端的 CourseDetail 接口 (课程详情)
type CourseDetailVO struct {
	ID            uint    `json:"id"`
	Room          string  `json:"room"`
	Faculty       string  `json:"faculty"`
	CourseName    string  `json:"courseName"`
	TeacherName   string  `json:"teacherName"`
	TeacherTitle  string  `json:"teacherTitle"`
	CourseTime    string  `json:"courseTime"`
	CourseType    string  `json:"courseType"`
	Description   string  `json:"description,omitempty"`
	Credits       float32 `json:"credits,omitempty"`
	AverageRating float32 `json:"rating,omitempty"`      // 平均评分
	ReviewCount   uint    `json:"reviewCount,omitempty"` // 评价数量
	// 可以根据需要添加更多详细信息
}

type CourseReviewInfoVO struct {
	ID           uint32    `json:"id"`                     // 评价本身的ID
	CourseID     uint32    `json:"courseId"`               // 课程ID
	Rating       int       `json:"rating"`                 // 课程评分
	Comment      string    `json:"comment"`                // 课程评论内容
	ReviewerName string    `json:"reviewerName,omitempty"` // 评价人名称 (需要从关联的 User 表获取)
	UserID       uint32    `json:"-"`                      // 评价者ID，前端可能不需要，但后端可能用
	CreatedAt    time.Time `json:"createdAt"`              // 评价时间
}
