package course

import (
	"time"
)

func GetTeachInfos() [][]BuildingTeachInfos {
	return getInfos()
}

func CurCourseTime() (weekNum int, weekday int, lessonNum int) {
	now := time.Now()
	// 计算第几周
	beginDate := time.Date(2024, time.September, 9,
		0, 0, 0, 0, time.Local)

	sub := now.Sub(beginDate)
	durationDay := int(sub.Hours()) / 24
	weekNum = durationDay/7 + 1

	// 计算周几
	weekday = int(now.Weekday())

	// 计算第几节课
	if isTimeBeforeHourAndMin(now, 7, 50) { // 8点前，早上
		lessonNum = -2
	} else if isTimeBeforeHourAndMin(now, 8, 45) { // 8点45前，第1节
		lessonNum = 1
	} else if isTimeBeforeHourAndMin(now, 9, 35) { // 9点35前，第2节
		lessonNum = 2
	} else if isTimeBeforeHourAndMin(now, 10, 35) { // 10点35前，第3节
		lessonNum = 3
	} else if isTimeBeforeHourAndMin(now, 11, 25) { // 11点25前，第4节
		lessonNum = 4
	} else if isTimeBeforeHourAndMin(now, 12, 15) { // 12点15前，第5节
		lessonNum = 5
	} else if isTimeBeforeHourAndMin(now, 13, 55) { // 中午，13点55之前, 没课
		lessonNum = -3
		// 中文没课捏
	} else if isTimeBeforeHourAndMin(now, 14, 50) { // 14点50前，第6节
		lessonNum = 6
	} else if isTimeBeforeHourAndMin(now, 15, 40) { // 15点40前，第7节课
		lessonNum = 7
	} else if isTimeBeforeHourAndMin(now, 16, 35) { // 16点35前，第8节课
		lessonNum = 8
	} else if isTimeBeforeHourAndMin(now, 17, 25) { // 17点25前，第9节课
		lessonNum = 9
	} else if isTimeBeforeHourAndMin(now, 18, 15) { // 18点15前，第10节课（一般不排课吧）
		lessonNum = 10
	} else if isTimeBeforeHourAndMin(now, 18, 20) { // 18点20前，晚饭时间，没课
		lessonNum = -4
		// 晚饭没课捏
	} else if isTimeBeforeHourAndMin(now, 19, 15) { // 19点15前，第11节课
		lessonNum = 11
	} else if isTimeBeforeHourAndMin(now, 20, 05) { // 20点05前，第12节课
		lessonNum = 12
	} else if isTimeBeforeHourAndMin(now, 20, 55) { // 20点55前，第13节课
		lessonNum = 13
	} else { // 晚上了
		// 今天的课上完了
		lessonNum = -5
	}
	return weekNum, weekday, lessonNum
}

func isTimeBeforeHourAndMin(inputTime time.Time, hour int, min int) bool {
	return inputTime.Before(
		time.Date(inputTime.Year(), inputTime.Month(), inputTime.Day(),
			hour, min, 0, 0,
			inputTime.Location()))
}
