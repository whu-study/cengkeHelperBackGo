package services

import (
	database "cengkeHelperBackGo/internal/db"
	"cengkeHelperBackGo/internal/models/dto"
	"cengkeHelperBackGo/internal/models/vo"
	"cengkeHelperBackGo/pkg/generator"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// CourseStructureService 课程结构化数据服务
type CourseStructureService struct{}

// NewCourseStructureService 创建课程结构化服务实例
func NewCourseStructureService() *CourseStructureService {
	return &CourseStructureService{}
}

// GetStructuredCourses 获取结构化的课程数据（学部 → 教学楼 → 楼层 → 课程）
func (s *CourseStructureService) GetStructuredCourses() ([]vo.DivisionVO, error) {
	// 获取所有课程和时间信息
	var timeInfos []dto.TimeInfo
	if err := database.Client.Find(&timeInfos).Error; err != nil {
		log.Printf("查询时间信息失败: %v", err)
		return nil, fmt.Errorf("查询课程时间信息失败: %w", err)
	}

	var courses []dto.CourseInfo
	if err := database.Client.Find(&courses).Error; err != nil {
		log.Printf("查询课程信息失败: %v", err)
		return nil, fmt.Errorf("查询课程信息失败: %w", err)
	}

	// 创建课程ID到课程的映射
	courseMap := make(map[uint32]*dto.CourseInfo)
	for i := range courses {
		courseMap[courses[i].ID] = &courses[i]
	}

	// 按学部组织数据
	divisionMap := make(map[string]*vo.DivisionVO)

	// 学部名称映射（根据Area字段）
	divisionNames := map[uint8]string{
		1: "文理学部",
		2: "信息学部",
		3: "工学部",
		4: "医学部",
	}

	for _, timeInfo := range timeInfos {
		course, exists := courseMap[timeInfo.CourseInfoId]
		if !exists {
			continue
		}

		// 获取或创建学部
		divisionID := fmt.Sprintf("division_%d", timeInfo.Area)
		divisionName := divisionNames[timeInfo.Area]
		if divisionName == "" {
			divisionName = fmt.Sprintf("学部%d", timeInfo.Area)
		}

		division, exists := divisionMap[divisionID]
		if !exists {
			division = &vo.DivisionVO{
				DivisionID:   divisionID,
				DivisionName: divisionName,
				Description:  fmt.Sprintf("%s教学区域", divisionName),
				Icon:         fmt.Sprintf("/assets/icons/%s.svg", divisionID),
				Buildings:    []vo.BuildingVO{},
			}
			divisionMap[divisionID] = division
		}

		// 获取或创建教学楼
		buildingID := fmt.Sprintf("%s_%s", divisionID, s.normalizeBuilding(timeInfo.Building))
		var building *vo.BuildingVO
		for i := range division.Buildings {
			if division.Buildings[i].BuildingID == buildingID {
				building = &division.Buildings[i]
				break
			}
		}
		if building == nil {
			building = &vo.BuildingVO{
				BuildingID:   buildingID,
				BuildingName: timeInfo.Building,
				BuildingCode: s.extractBuildingCode(timeInfo.Building),
				Address:      fmt.Sprintf("武汉大学%s", divisionName),
				Floors:       []vo.FloorVO{},
			}
			division.Buildings = append(division.Buildings, *building)
			building = &division.Buildings[len(division.Buildings)-1]
		}

		// 从教室编号提取楼层信息
		floorNumber := s.extractFloorNumber(timeInfo.Classroom)
		floorID := fmt.Sprintf("%s_F%d", buildingID, floorNumber)

		// 获取或创建楼层
		var floor *vo.FloorVO
		for i := range building.Floors {
			if building.Floors[i].FloorID == floorID {
				floor = &building.Floors[i]
				break
			}
		}
		if floor == nil {
			floor = &vo.FloorVO{
				FloorID:     floorID,
				FloorName:   fmt.Sprintf("%s %d层", building.BuildingCode, floorNumber),
				FloorNumber: floorNumber,
				Rooms:       []vo.RoomVO{},
				Courses:     []vo.CourseInfoVO{},
			}
			building.Floors = append(building.Floors, *floor)
			floor = &building.Floors[len(building.Floors)-1]
		}

		// 添加教室（如果不存在）
		roomID := fmt.Sprintf("%s_%s", floorID, timeInfo.Classroom)
		roomExists := false
		for _, room := range floor.Rooms {
			if room.RoomID == roomID {
				roomExists = true
				break
			}
		}
		if !roomExists {
			floor.Rooms = append(floor.Rooms, vo.RoomVO{
				RoomID:     roomID,
				RoomNumber: timeInfo.Classroom,
				RoomName:   fmt.Sprintf("教室 %s", timeInfo.Classroom),
				Facilities: []string{"投影仪", "空调", "网络"},
			})
		}

		// 解析学分
		credits := s.parseCredits(course.Credit)

		// 解析时间段
		timeSlots := s.parseTimeSlots(timeInfo)

		// 添加课程信息
		courseVO := vo.CourseInfoVO{
			ID:            course.ID,
			CourseName:    course.CourseName,
			CourseCode:    course.CourseNum,
			TeacherName:   course.Teacher,
			TeacherTitle:  course.TeacherTitle,
			Faculty:       course.Faculty,
			Credits:       credits,
			CourseType:    course.CourseType,
			Room:          timeInfo.Classroom,
			TimeSlots:     timeSlots,
			Description:   course.Description,
			AverageRating: course.AverageRating,
			ReviewCount:   course.ReviewCount,
		}
		floor.Courses = append(floor.Courses, courseVO)
	}

	// 计算统计信息并转换为数组
	result := make([]vo.DivisionVO, 0, len(divisionMap))
	for _, division := range divisionMap {
		totalFloors := 0
		totalCourses := 0
		totalRooms := 0

		for i := range division.Buildings {
			building := &division.Buildings[i]
			buildingCourses := 0
			buildingRooms := 0

			for j := range building.Floors {
				floor := &building.Floors[j]
				buildingCourses += len(floor.Courses)
				buildingRooms += len(floor.Rooms)
			}

			building.TotalFloors = len(building.Floors)
			building.TotalCourses = buildingCourses
			building.TotalRooms = buildingRooms

			totalFloors += len(building.Floors)
			totalCourses += buildingCourses
			totalRooms += buildingRooms
		}

		division.TotalBuildings = len(division.Buildings)
		division.TotalFloors = totalFloors
		division.TotalCourses = totalCourses

		result = append(result, *division)
	}

	return result, nil
}

// normalizeBuilding 规范化教学楼名称作为ID的一部分
func (s *CourseStructureService) normalizeBuilding(building string) string {
	// 移除空格和特殊字符，保留字母、数字和中文
	// Go 的正则表达式使用 \p{Han} 匹配中文字符
	reg := regexp.MustCompile(`[^\p{Han}a-zA-Z0-9]+`)
	normalized := reg.ReplaceAllString(building, "_")
	return strings.ToLower(normalized)
}

// extractBuildingCode 从教学楼名称提取代码
func (s *CourseStructureService) extractBuildingCode(building string) string {
	// 尝试提取英文字母或数字编号
	reg := regexp.MustCompile(`[A-Z]+\d*|\d+号楼`)
	match := reg.FindString(building)
	if match != "" {
		return match
	}
	// 如果没有匹配，返回前3个字符
	runes := []rune(building)
	if len(runes) > 3 {
		return string(runes[:3])
	}
	return building
}

// extractFloorNumber 从教室编号提取楼层号
func (s *CourseStructureService) extractFloorNumber(classroom string) int {
	// 常见格式: A101, 201, 3-101 等
	// 提取第一个或第二个数字作为楼层号
	reg := regexp.MustCompile(`\d+`)
	matches := reg.FindAllString(classroom, -1)

	if len(matches) > 0 {
		// 取第一个数字序列
		numStr := matches[0]
		if len(numStr) >= 3 {
			// 如果是3位或以上数字，第一位是楼层
			floor, _ := strconv.Atoi(string(numStr[0]))
			if floor > 0 {
				return floor
			}
		} else if len(numStr) >= 1 {
			// 否则整个数字可能是楼层
			floor, _ := strconv.Atoi(numStr)
			if floor > 0 && floor < 50 {
				return floor
			}
		}
	}

	return 1 // 默认第一层
}

// parseCredits 解析学分字符串为浮点数
func (s *CourseStructureService) parseCredits(creditStr string) float32 {
	// 移除非数字字符
	reg := regexp.MustCompile(`[\d.]+`)
	match := reg.FindString(creditStr)
	if match != "" {
		if credit, err := strconv.ParseFloat(match, 32); err == nil {
			return float32(credit)
		}
	}
	return 0
}

// parseTimeSlots 解析时间段信息
func (s *CourseStructureService) parseTimeSlots(timeInfo dto.TimeInfo) []vo.TimeSlotVO {
	// 从二进制数据中提取周次和节次
	weeks, lessons := generator.Bin2WeekLesson(timeInfo.WeekAndTime)

	if len(lessons) == 0 {
		return []vo.TimeSlotVO{}
	}

	// 格式化周次范围
	weekRange := s.formatWeekRange(weeks)

	// 格式化节次范围
	startPeriod := lessons[0]
	endPeriod := lessons[len(lessons)-1]

	return []vo.TimeSlotVO{
		{
			DayOfWeek:   int(timeInfo.DayOfWeek),
			StartPeriod: startPeriod,
			EndPeriod:   endPeriod,
			Weeks:       weekRange,
		},
	}
}

// formatWeekRange 格式化周次范围
func (s *CourseStructureService) formatWeekRange(weeks []int) string {
	if len(weeks) == 0 {
		return ""
	}
	if len(weeks) == 1 {
		return fmt.Sprintf("第%d周", weeks[0])
	}
	// 简化处理：如果是连续的，显示范围
	if weeks[len(weeks)-1]-weeks[0] == len(weeks)-1 {
		return fmt.Sprintf("%d-%d周", weeks[0], weeks[len(weeks)-1])
	}
	// 否则列出所有周次
	return fmt.Sprintf("%d-%d周", weeks[0], weeks[len(weeks)-1])
}

// GetRoomFacilities 解析设施JSON字符串
func (s *CourseStructureService) getRoomFacilities(facilitiesJSON string) []string {
	if facilitiesJSON == "" {
		return []string{}
	}
	var facilities []string
	if err := json.Unmarshal([]byte(facilitiesJSON), &facilities); err != nil {
		return []string{}
	}
	return facilities
}
