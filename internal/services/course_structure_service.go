package services

import (
	database "cengkeHelperBackGo/internal/db"
	"cengkeHelperBackGo/internal/models/dto"
	"cengkeHelperBackGo/internal/models/vo"
	"cengkeHelperBackGo/pkg/generator"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
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
	beginDate := time.Date(2025, time.September, 8, 0, 0, 0, 0, time.Local)
	return s.TimeToNums(now, beginDate)
}

func (s *CourseStructureService) TimeToNums(t, beginDate time.Time) (weekNum int, weekday int, lessonNum int) {
	// 计算第几周（学期开始日期：2025年9月8日）
	sub := t.Sub(beginDate)
	durationDay := int(sub.Hours()) / 24
	weekNum = durationDay/7 + 1

	// 计算周几（0=周日, 1=周一, ..., 6=周六）
	weekday = int(t.Weekday())

	// 计算第几节课
	hour := t.Hour()
	minute := t.Minute()

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

func (s *CourseStructureService) ValidParams(params *CourseQueryParams) *CourseQueryParams {
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
	return params
}

// GetEmptyDivisionStructure 获取空的学部结构（保留学部层级，但buildings为空）
// 用于非上课时间或没有课程数据时返回基本的学部框架
func (s *CourseStructureService) GetEmptyDivisionStructure(divisionID *int) []vo.DivisionVO {
	// 学部名称映射（根据Area字段）
	divisionNames := map[uint8]string{
		1: "文理学部",
		2: "信息学部",
		3: "工学部",
		4: "医学部",
	}

	var result []vo.DivisionVO

	if divisionID != nil {
		// 只返回指定学部
		if *divisionID >= 1 && *divisionID <= 4 {
			divisionName := divisionNames[uint8(*divisionID)]
			result = append(result, vo.DivisionVO{
				DivisionID:     fmt.Sprintf("division_%d", *divisionID),
				DivisionName:   divisionName,
				Description:    fmt.Sprintf("%s教学区域", divisionName),
				TotalBuildings: 0,
				TotalFloors:    0,
				TotalCourses:   0,
				Buildings:      []*vo.BuildingVO{},
			})
		}
	} else {
		// 返回所有学部
		for i := 1; i <= 4; i++ {
			divisionName := divisionNames[uint8(i)]
			result = append(result, vo.DivisionVO{
				DivisionID:     fmt.Sprintf("division_%d", i),
				DivisionName:   divisionName,
				Description:    fmt.Sprintf("%s教学区域", divisionName),
				TotalBuildings: 0,
				TotalFloors:    0,
				TotalCourses:   0,
				Buildings:      []*vo.BuildingVO{},
			})
		}
	}

	return result
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
