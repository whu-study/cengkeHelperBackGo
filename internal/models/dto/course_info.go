package dto

import (
	"cengkeHelperBackGo/internal/models/vo"
	"time"
)

// 课程信息表
type CourseInfo struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BuildingID   uint      `gorm:"index;not null;comment:所属教学楼ID" json:"-"`
	Room         string    `gorm:"type:varchar(50);not null;comment:教室编号" json:"room"`
	Faculty      string    `gorm:"type:varchar(100);not null;comment:所属院系" json:"faculty"`
	CourseName   string    `gorm:"type:varchar(255);not null;comment:课程名称" json:"courseName"`
	TeacherName  string    `gorm:"type:varchar(50);not null;comment:教师姓名" json:"teacherName"`
	TeacherTitle string    `gorm:"type:varchar(20);comment:教师职称" json:"teacherTitle"`
	CourseTime   string    `gorm:"type:varchar(50);not null;comment:上课时间" json:"courseTime"`
	CourseType   string    `gorm:"type:varchar(50);not null;comment:课程类型" json:"courseType"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func convertCoursesToVO(courses []CourseInfo) []vo.CourseInfoVO {
	var voCourses []vo.CourseInfoVO
	for _, c := range courses {
		voCourses = append(voCourses, vo.CourseInfoVO{
			Room:         c.Room,
			Faculty:      c.Faculty,
			CourseName:   c.CourseName,
			TeacherName:  c.TeacherName,
			TeacherTitle: c.TeacherTitle,
			CourseTime:   c.CourseTime,
			CourseType:   c.CourseType,
		})
	}
	return voCourses
}
