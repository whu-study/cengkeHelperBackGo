package services

import (
	database "cengkeHelperBackGo/internal/db"
	"cengkeHelperBackGo/internal/models/dto"
	"cengkeHelperBackGo/internal/models/vo"
	"cengkeHelperBackGo/pkg/generator"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// CourseStructureService 课程结构化数据服务
type CourseStructureService struct{}

// NewCourseStructureService 创建课程结构化服务实例
func NewCourseStructureService() *CourseStructureService {
	return &CourseStructureService{}
}

// CourseQueryParams 课程查询参数
type CourseQueryParams struct {
	WeekNum    int  // 周次，-1 表示不限，0 表示使用当前时间
	Weekday    int  // 星期几，-1 表示不限，0 表示使用当前时间
	LessonNum  int  // 节次，-1 表示不限，0 表示使用当前时间
	DivisionID *int // 学部ID (1-4)，nil 表示不限
	UseCache   bool // 是否使用缓存
}

// GetCurrentCourseTime 获取当前的课程时间（周次、星期、节次）
func (s *CourseStructureService) GetCurrentCourseTime() (weekNum int, weekday int, lessonNum int) {
	now := time.Now()
	// 计算第几周（学期开始日期：2025年9月8日）
	beginDate := time.Date(2025, time.September, 8, 0, 0, 0, 0, time.Local)
	sub := now.Sub(beginDate)
	durationDay := int(sub.Hours()) / 24
	weekNum = durationDay/7 + 1

	// 计算周几（0=周日, 1=周一, ..., 6=周六）
	weekday = int(now.Weekday())

	// 计算第几节课
	hour := now.Hour()
	minute := now.Minute()

	switch {
	case hour < 7 || (hour == 7 && minute < 50):
		lessonNum = -1 // 早上还没开始上课
	case hour < 8 || (hour == 8 && minute < 45):
		lessonNum = 1
	case hour < 9 || (hour == 9 && minute < 35):
		lessonNum = 2
	case hour < 10 || (hour == 10 && minute < 35):
		lessonNum = 3
	case hour < 11 || (hour == 11 && minute < 25):
		lessonNum = 4
	case hour < 12 || (hour == 12 && minute < 15):
		lessonNum = 5
	case hour < 13 || (hour == 13 && minute < 55):
		lessonNum = -1 // 中午休息
	case hour < 14 || (hour == 14 && minute < 50):
		lessonNum = 6
	case hour < 15 || (hour == 15 && minute < 40):
		lessonNum = 7
	case hour < 16 || (hour == 16 && minute < 35):
		lessonNum = 8
	case hour < 17 || (hour == 17 && minute < 25):
		lessonNum = 9
	case hour < 18 || (hour == 18 && minute < 15):
		lessonNum = 10
	case hour < 18 || (hour == 18 && minute < 20):
		lessonNum = -1 // 晚饭时间
	case hour < 19 || (hour == 19 && minute < 15):
		lessonNum = 11
	case hour < 20 || (hour == 20 && minute < 5):
		lessonNum = 12
	case hour < 20 || (hour == 20 && minute < 55):
		lessonNum = 13
	default:
		lessonNum = -1 // 晚上没课了
	}

	return weekNum, weekday, lessonNum
}

// GetStructuredCourses 获取结构化的课程数据（学部 → 教学楼 → 楼层 → 课程）
// 默认返回当前时间的课程
func (s *CourseStructureService) GetStructuredCourses(params *CourseQueryParams) ([]vo.DivisionVO, error) {
	// 如果没有传参数，使用当前时间
	if params == nil {
		weekNum, weekday, lessonNum := s.GetCurrentCourseTime()
		params = &CourseQueryParams{
			WeekNum:   weekNum,
			Weekday:   weekday,
			LessonNum: lessonNum,
			UseCache:  true,
		}
	} else {
		// 如果参数中某些值为 0，表示使用当前时间
		currentWeekNum, currentWeekday, currentLessonNum := s.GetCurrentCourseTime()

		if params.WeekNum == 0 {
			params.WeekNum = currentWeekNum
		}
		if params.Weekday == 0 {
			params.Weekday = currentWeekday
		}
		if params.LessonNum == 0 {
			params.LessonNum = currentLessonNum
		}
	}

	// 尝试从缓存获取
	if params.UseCache {
		cacheKey := s.getCacheKey(params)
		if cachedData, err := s.getFromCache(cacheKey); err == nil && cachedData != nil {
			log.Printf("从缓存获取课程数据: %s", cacheKey)
			return cachedData, nil
		}
	}

	// 构建查询条件
	var timeInfos []dto.TimeInfo
	query := database.Client.Model(&dto.TimeInfo{})

	// 添加学部过滤
	if params.DivisionID != nil {
		query = query.Where("area = ?", *params.DivisionID)
	}

	// 添加星期过滤（只有当weekday >= 0时才过滤）
	if params.Weekday >= 0 {
		query = query.Where("day_of_week = ?", params.Weekday)
	}

	if err := query.Find(&timeInfos).Error; err != nil {
		log.Printf("查询时间信息失败: %v", err)
		return nil, fmt.Errorf("查询课程时间信息失败: %w", err)
	}

	// 获取课程ID列表
	courseIDs := make([]uint32, 0, len(timeInfos))
	for _, ti := range timeInfos {
		courseIDs = append(courseIDs, ti.CourseInfoId)
	}

	// 只查询需要的课程信息
	var courses []dto.CourseInfo
	if len(courseIDs) > 0 {
		if err := database.Client.Where("id IN ?", courseIDs).Find(&courses).Error; err != nil {
			log.Printf("查询课程信息失败: %v", err)
			return nil, fmt.Errorf("查询课程信息失败: %w", err)
		}
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

	// 课程去重：使用 courseID + classroom 作为唯一键
	addedCourses := make(map[string]bool)

	for _, timeInfo := range timeInfos {
		// 如果指定了周次和节次，进行二进制匹配过滤
		if params.WeekNum != -1 || params.LessonNum != -1 {
			if !generator.IsWeekLessonMatch(params.WeekNum, params.LessonNum, timeInfo.WeekAndTime) {
				continue
			}
		}

		course, exists := courseMap[timeInfo.CourseInfoId]
		if !exists {
			continue
		}

		// 去重：同一门课在同一教室只添加一次
		courseKey := fmt.Sprintf("%d_%s", course.ID, timeInfo.Classroom)
		if addedCourses[courseKey] {
			continue
		}
		addedCourses[courseKey] = true

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

		// 生成课程时间文本（如 "1-2节"）
		courseTime := s.formatCourseTime(timeInfo, params.LessonNum)

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
			CourseTime:    courseTime, // 填充课程时间文本
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

	// 存入缓存（5分钟过期）
	if params.UseCache {
		cacheKey := s.getCacheKey(params)
		_ = s.saveToCache(cacheKey, result, 5*time.Minute)
	}

	return result, nil
}

// getCacheKey 生成缓存键
func (s *CourseStructureService) getCacheKey(params *CourseQueryParams) string {
	divisionStr := "all"
	if params.DivisionID != nil {
		divisionStr = fmt.Sprintf("%d", *params.DivisionID)
	}
	return fmt.Sprintf("course_structure:%s:w%d:d%d:l%d",
		divisionStr, params.WeekNum, params.Weekday, params.LessonNum)
}

// getFromCache 从Redis缓存获取数据
func (s *CourseStructureService) getFromCache(key string) ([]vo.DivisionVO, error) {
	ctx := context.Background()
	val, err := database.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var result []vo.DivisionVO
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// saveToCache 保存数据到Redis缓存
func (s *CourseStructureService) saveToCache(key string, data []vo.DivisionVO, expiration time.Duration) error {
	ctx := context.Background()
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return database.RedisClient.Set(ctx, key, jsonData, expiration).Err()
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

// formatCourseTime 生成课程时间文本
func (s *CourseStructureService) formatCourseTime(timeInfo dto.TimeInfo, lessonNum int) string {
	// 从二进制数据中提取周次和节次
	weeks, lessons := generator.Bin2WeekLesson(timeInfo.WeekAndTime)

	// 处理当前时间的课程（周次为0）
	if timeInfo.WeekAndTime == 0 {
		return "当前时间"
	}

	// 课程节次范围
	startPeriod := lessons[0]
	endPeriod := lessons[len(lessons)-1]

	// 生成时间文本
	var weekText string
	if len(weeks) == 1 {
		weekText = fmt.Sprintf("第%d周", weeks[0])
	} else {
		weekText = fmt.Sprintf("%d-%d周", weeks[0], weeks[len(weeks)-1])
	}

	var lessonText string
	if startPeriod == endPeriod {
		lessonText = fmt.Sprintf("%d节", startPeriod)
	} else {
		lessonText = fmt.Sprintf("%d-%d节", startPeriod, endPeriod)
	}

	return fmt.Sprintf("%s %s", weekText, lessonText)
}
