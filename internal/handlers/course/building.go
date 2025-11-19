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
	CourseNum    string `json:"courseNum,omitempty"`
	CourseName   string `json:"courseName"`
	TeacherName  string `json:"teacherName"`
	TeacherTitle string `json:"teacherTitle"`
	CourseTime   string `json:"courseTime"`
	CourseType   string `json:"courseType"`
	Description  string `json:"description,omitempty"`
	Credit       string `json:"credit,omitempty"`

	WeekAndTime uint32 `json:"weekAndTime,omitempty"`
	DayOfWeek   int    `json:"dayOfWeek,omitempty"`

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

	ReviewCount   uint32
	AverageRating float32

	Credit string

	DayOfWeek int

	CourseType string
}

// BuildingTeachInfos 每个学部各个教学楼的课程信息
type BuildingTeachInfos struct {
	Building string          `json:"building"`
	Infos    []RespTeachInfo `json:"infos"`
}

var RespTeachInfos = make([][]BuildingTeachInfos, 5)

func searchByAreaAndWeekday(areaNum int, weekday int, weekNum int, lessonNum int) []MapTeachInfo {
	tempInfo := make([]MapTeachInfo, 0)
	weekLessonBin := generator.WeekLesson2Bin([]int{weekNum}, []int{lessonNum})
	if err := database.Client.
		Raw(queryStr,
			weekday, areaNum, weekLessonBin, weekLessonBin).
		Find(&tempInfo).Error; err != nil {
		log.Fatal(err)
	}

	return tempInfo
}

func GetInfos(weekNum, weekday, lessonNum int) [][]BuildingTeachInfos {
	for i := 0; i < 5; i++ {
		RespTeachInfos[i] = make([]BuildingTeachInfos, 0)
	}

	for i := 1; i <= 4; i++ {
		buildingMap := make(map[string][]RespTeachInfo)

		for _, info := range searchByAreaAndWeekday(i, weekday, weekNum, lessonNum) {
			// 数据库已经完成过滤，不再需要内存中的二次过滤
			res := RespTeachInfo{
				ID:            info.ID,
				CourseNum:     info.CourseNum,
				Room:          info.Classroom,
				Faculty:       info.Faculty,
				CourseName:    info.CourseName,
				TeacherName:   info.Teacher,
				TeacherTitle:  info.TeacherTitle,
				CourseTime:    generator.NearestToDisplay(lessonNum, info.WeekAndTime),
				CourseType:    info.CourseType,
				Credit:        info.Credit,
				AverageRating: info.AverageRating,
				ReviewCount:   info.ReviewCount,

				WeekAndTime: info.WeekAndTime,
				DayOfWeek:   info.DayOfWeek,
			}
			// 去重：按 courseNum 在同一教学楼内去重，避免同一课程因为不同教室/时段重复出现
			existing := false
			for _, ex := range buildingMap[info.Building] {
				if ex.CourseNum != "" && ex.CourseNum == info.CourseNum {
					existing = true
					break
				}
				// 作为兜底，如果 CourseNum 缺失，则用 ID 检查去重
				if ex.CourseNum == "" && ex.ID == info.ID {
					existing = true
					break
				}
			}
			if !existing {
				buildingMap[info.Building] = append(buildingMap[info.Building], res)
			}
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
