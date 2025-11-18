package vo

import "time"

// CourseInfoVO 课程信息VO（完全匹配前端接口）
type CourseInfoVO struct {
	ID            uint32       `json:"id"`
	CourseName    string       `json:"courseName"`
	CourseCode    string       `json:"courseCode,omitempty"`
	TeacherName   string       `json:"teacherName"`
	TeacherTitle  string       `json:"teacherTitle"`
	Faculty       string       `json:"faculty"`
	Credits       float32      `json:"credits,omitempty"`
	CourseType    string       `json:"courseType"`
	Room          string       `json:"room"`
	TimeSlots     []TimeSlotVO `json:"timeSlots,omitempty"`
	Capacity      int          `json:"capacity,omitempty"`
	Enrolled      int          `json:"enrolled,omitempty"`
	Description   string       `json:"description,omitempty"`
	CourseTime    string       `json:"courseTime,omitempty"` // 保留用于简单展示
	AverageRating float32      `json:"averageRating,omitempty"`
	ReviewCount   uint32       `json:"reviewCount,omitempty"`
}

// CourseDetailVO 对应前端的 CourseDetail 接口 (课程详情)
type CourseDetailVO struct {
	ID            uint         `json:"id"`
	CourseName    string       `json:"courseName"`
	CourseCode    string       `json:"courseCode,omitempty"`
	TeacherName   string       `json:"teacherName"`
	TeacherTitle  string       `json:"teacherTitle"`
	Faculty       string       `json:"faculty"`
	Credits       float32      `json:"credits,omitempty"`
	CourseType    string       `json:"courseType"`
	Room          string       `json:"room"`
	TimeSlots     []TimeSlotVO `json:"timeSlots,omitempty"`
	Capacity      int          `json:"capacity,omitempty"`
	Enrolled      int          `json:"enrolled,omitempty"`
	Description   string       `json:"description,omitempty"`
	CourseTime    string       `json:"courseTime,omitempty"`
	AverageRating float32      `json:"rating,omitempty"`
	ReviewCount   uint         `json:"reviewCount,omitempty"`
}

// CourseReviewInfoVO 课程评价信息VO
type CourseReviewInfoVO struct {
	ID           uint32    `json:"id"`
	CourseID     uint32    `json:"courseId"`
	Rating       int       `json:"rating"`
	Comment      string    `json:"comment"`
	ReviewerName string    `json:"reviewerName,omitempty"`
	UserID       uint32    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
}
