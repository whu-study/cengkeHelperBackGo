package vo

// CourseInfoVO 课程信息VO（完全匹配前端接口）
type CourseInfoVO struct {
	Room         string `json:"room"`
	Faculty      string `json:"faculty"`
	CourseName   string `json:"courseName"`
	TeacherName  string `json:"teacherName"`
	TeacherTitle string `json:"teacherTitle"`
	CourseTime   string `json:"courseTime"`
	CourseType   string `json:"courseType"`
}
