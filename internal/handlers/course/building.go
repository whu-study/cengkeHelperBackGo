package course

import (
	database "cengkeHelperBackGo/internal/db"
	"cengkeHelperBackGo/pkg/generator"
	"log"
	"slices"
)

type RespTeachInfo struct {
	ID           uint32 `json:"id"`
	Room         string `json:"room"`
	Faculty      string `json:"faculty"`
	CourseName   string `json:"courseName"`
	TeacherName  string `json:"teacherName"`
	TeacherTitle string `json:"teacherTitle"`
	CourseTime   string `json:"courseTime"`
	CourseType   string `json:"courseType"`
	Description  string `json:"description,omitempty"`

	AverageRating float32 `json:"rating,omitempty"`
	ReviewCount   uint32  `json:"reviewCount,omitempty"`
}
type MapTeachInfo struct {
	CourseNum    string
	ID           uint32
	Classroom    string
	Faculty      string
	CourseName   string
	Teacher      string
	TeacherTitle string
	WeekAndTime  uint32
	Building     string

	DayOfWeek string

	CourseType string
}

// BuildingTeachInfos 每个学部各个教学楼的课程信息
type BuildingTeachInfos struct {
	Building string          `json:"building"`
	Infos    []RespTeachInfo `json:"infos"`
}

var RespTeachInfos = make([][]BuildingTeachInfos, 5)

func searchByAreaAndWeekday(areaNum int, weekday int) []MapTeachInfo {
	tempInfo := make([]MapTeachInfo, 0)
	if err := database.Client.
		Raw(queryStr,
			weekday, areaNum).
		Find(&tempInfo).Error; err != nil {
		log.Fatal(err)
	}

	return tempInfo
}

func getInfos(weekNum, weekday, lessonNum int) [][]BuildingTeachInfos {
	for i := 0; i < 5; i++ {
		RespTeachInfos[i] = make([]BuildingTeachInfos, 0)
	}

	for i := 1; i <= 4; i++ {
		buildingMap := make(map[string][]RespTeachInfo)

		for _, info := range searchByAreaAndWeekday(i, weekday) {
			if !generator.IsWeekLessonMatch(weekNum, lessonNum, info.WeekAndTime) {
				continue
			}

			res := RespTeachInfo{
				ID:           info.ID,
				Room:         info.Classroom,
				Faculty:      info.Faculty,
				CourseName:   info.CourseName,
				TeacherName:  info.Teacher,
				TeacherTitle: info.TeacherTitle,
				CourseTime:   generator.NearestToDisplay(lessonNum, info.WeekAndTime),
				CourseType:   info.CourseType,
			}
			//_, lesson := generator.Bin2WeekLesson(info.WeekAndTime)
			//logger.Info(res, lesson)
			buildingMap[info.Building] = append(buildingMap[info.Building], res)
		}
		for key, infos := range buildingMap {
			RespTeachInfos[i-1] = append(RespTeachInfos[i-1], BuildingTeachInfos{
				Building: key,
				Infos:    infos,
			})
		}

	}

	// 每个学部的教学楼按照课程数量排序
	for i := range RespTeachInfos {
		slices.SortFunc(RespTeachInfos[i], func(a, b BuildingTeachInfos) int {
			return len(b.Infos) - len(a.Infos)
		})
	}

	return RespTeachInfos
}
