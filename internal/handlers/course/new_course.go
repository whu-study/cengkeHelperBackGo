package course

import (
	"cengkeHelperBackGo/internal/models/vo"
	"cengkeHelperBackGo/internal/services/course"
	"fmt"
)

var divisionNames = map[int]string{
	1: "文理学部",
	2: "信息学部",
	3: "工学部",
	4: "医学部",
}

func GetStructuredCourses(dayOfWeek int, weekNum int, lessonNum int) []vo.DivisionVO {

	infos := GetInfos(weekNum, dayOfWeek, lessonNum)

	result := make([]vo.DivisionVO, 0, 5)

	for i, buildingInfos := range infos {
		division := vo.DivisionVO{
			DivisionID:     "division_" + fmt.Sprintf("%d", i+1),
			DivisionName:   divisionNames[i],
			Description:    fmt.Sprintf("%s教学区域", divisionNames[i]),
			TotalBuildings: len(buildingInfos),
			TotalFloors:    0,
			TotalCourses:   0,
			Buildings:      make([]*vo.BuildingVO, 0),
		}

		for _, building := range buildingInfos {
			buildingVO := &vo.BuildingVO{
				BuildingID:   fmt.Sprintf("division_%d_%s", i+1, building.Building),
				BuildingName: building.Building,
				BuildingCode: course.ExtractBuildingCode(building.Building),
				Address:      fmt.Sprintf("武汉大学%s", divisionNames[i]),
				Description:  "",
				TotalFloors:  0,
				TotalRooms:   len(building.Infos),
				TotalCourses: len(building.Infos),
				Floors:       make([]*vo.FloorVO, 0),
			}
			floors := make(map[int]*vo.FloorVO)
			for _, info := range building.Infos {
				floorNumber := course.ExtractFloorNumber(info.Room)
				if _, exists := floors[floorNumber]; !exists {
					floors[floorNumber] = &vo.FloorVO{
						FloorID:     fmt.Sprintf("division_%d_%s_F%d", i+1, building.Building, floorNumber),
						FloorName:   fmt.Sprintf("%s %d层", course.ExtractBuildingCode(building.Building), floorNumber),
						FloorNumber: floorNumber,
						Description: "",
						Rooms:       make([]*vo.RoomVO, 0),
						Courses:     make([]*vo.CourseInfoVO, 0),
					}
				}

				floors[floorNumber].Rooms = append(floors[floorNumber].Rooms, &vo.RoomVO{
					RoomID:     fmt.Sprintf("division_%d_%s_F%d_%s", i+1, building.Building, floorNumber, info.Room),
					RoomNumber: info.Room,
					RoomName:   fmt.Sprintf("教室 %s", info.Room),
					Capacity:   0,
					RoomType:   "",
					Facilities: []string{"投影仪", "空调", "网络"},
				})

				floors[floorNumber].Courses = append(floors[floorNumber].Courses, &vo.CourseInfoVO{
					ID:            info.ID,
					CourseName:    info.CourseName,
					CourseCode:    info.CourseNum,
					TeacherName:   info.TeacherName,
					TeacherTitle:  info.TeacherTitle,
					Faculty:       info.Faculty,
					Credits:       course.ParseCredits(info.Credits),
					CourseType:    info.CourseType,
					Room:          info.Room,
					TimeSlots:     course.ParseTimeSlots(info.WeekAndTime, info.DayOfWeek),
					Capacity:      0,
					Enrolled:      0,
					Description:   "",
					CourseTime:    info.CourseTime,
					AverageRating: info.AverageRating,
					ReviewCount:   info.ReviewCount,
				})
			}
			buildingVO.TotalFloors = len(floors)
			for _, floor := range floors {
				buildingVO.Floors = append(buildingVO.Floors, floor)
			}
			division.Buildings = append(division.Buildings, buildingVO)
		}

		result = append(result, division)

	}

	return result
}
