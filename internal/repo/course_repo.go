package repo

import (
	database "cengkeHelperBackGo/internal/db"
	"log"
)

// CourseRow 是 SQL 查询的返回行
type CourseRow struct {
	ID           uint32 `gorm:"column:id"`
	CourseNum    string `gorm:"column:course_num"`
	Classroom    string `gorm:"column:classroom"`
	Building     string `gorm:"column:building"`
	CourseType   string `gorm:"column:course_type"`
	Faculty      string `gorm:"column:faculty"`
	CourseName   string `gorm:"column:course_name"`
	Teacher      string `gorm:"column:teacher"`
	TeacherTitle string `gorm:"column:teacher_title"`
	WeekAndTime  uint32 `gorm:"column:week_and_time"`
	DayOfWeek    uint8  `gorm:"column:day_of_week"`
	Area         uint8  `gorm:"column:area"`
}

// QueryStr 按教学楼+course_num 去重，返回每门课的一条代表行
var QueryStr = `
SELECT 
    MAX(ci.id) AS id,
    ci.course_num,
    ANY_VALUE(ti.classroom) AS classroom,
    ti.building,
    ANY_VALUE(ci.course_type) AS course_type,
    ANY_VALUE(ci.faculty) AS faculty,
    ANY_VALUE(ci.course_name) AS course_name,
    ANY_VALUE(ci.teacher) AS teacher,
    ANY_VALUE(ci.teacher_title) AS teacher_title,
    MAX(ti.week_and_time) AS week_and_time,
    MAX(ti.day_of_week) AS day_of_week,
    ti.area AS area
FROM time_infos ti
JOIN course_infos ci ON ci.id = ti.course_info_id
WHERE ti.day_of_week = ?
  AND (? = -1 OR ti.area = ?)
  AND (? = -1 OR (ti.week_and_time & (1 << (32 - ?))) != 0)
  AND (? = -1 OR (ti.week_and_time & (1 << (? - 1))) != 0)
GROUP BY
    ti.building,
    ci.course_num
`

// SearchByAreaAndWeekday 查询函数：dayOfWeek, area (-1 表示所有), weekNum (-1 表示不限), lessonNum (-1 表示不限)
func SearchByAreaAndWeekday(dayOfWeek, area, weekNum, lessonNum int) ([]CourseRow, error) {
	rows := make([]CourseRow, 0)
	err := database.Client.Raw(QueryStr, dayOfWeek, area, area, weekNum, weekNum, lessonNum, lessonNum).Scan(&rows).Error
	if err != nil {
		log.Printf("repo.SearchByAreaAndWeekday error: %v", err)
		return nil, err
	}
	return rows, nil
}
