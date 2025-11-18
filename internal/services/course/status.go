package course

import (
	database "cengkeHelperBackGo/internal/db"
	"cengkeHelperBackGo/pkg/generator"
	"fmt"
)

func getNumOfCourses(dayOfWeek, weekNum int, lessonNum []int) int {
	// 计算 weekAndTime 掩码
	weekAndTime := generator.WeekLesson2Bin([]int{weekNum}, lessonNum)

	var count int
	if err := database.Client.Raw(
		`SELECT COUNT(DISTINCT ci.course_num)
		 FROM time_infos ti 
		 JOIN course_infos ci ON ci.id = ti.course_info_id
		 WHERE ti.day_of_week = ? 
           AND (ti.week_and_time & ? & (-1<<13)) != 0
           AND (ti.week_and_time & ? & 0x1FFF) != 0`,
		dayOfWeek, weekAndTime, weekAndTime,
	).Scan(&count).Error; err != nil {
		fmt.Println(err)
	}

	return count
}

func GetSingleNumOfCourses(dayOfWeek, weekNum int, lessonNum int) int {
	return getNumOfCourses(dayOfWeek, weekNum, []int{lessonNum})
}

func GetOneDayNumOfCourses(dayOfWeek, weekNum int) int {
	lessonNums := make([]int, 13)
	for i := 1; i <= 13; i++ {
		lessonNums[i-1] = i
	}
	return getNumOfCourses(dayOfWeek, weekNum, lessonNums)
}
