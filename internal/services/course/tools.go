package course

import (
	"cengkeHelperBackGo/internal/models/vo"
	"cengkeHelperBackGo/pkg/generator"
	"fmt"
	"regexp"
	"strconv"
	"unicode"
)

var reg = regexp.MustCompile(`\d+`)

// ExtractFloorNumber 从教室编号提取楼层号
func ExtractFloorNumber(classroom string) int {
	// 常见格式: A101, 201, 3-101 等
	// 提取第一个或第二个数字作为楼层号
	matches := reg.FindAllString(classroom, -1)
	// 提取第一个不为0的数字序列
	if len(matches) == 0 {
		return 1 // 默认第一层
	}

	// 取第一个数字序列
	numStr := matches[0]
	for i := 0; i < len(numStr); i++ {
		if numStr[i] != '0' && unicode.IsDigit(rune(numStr[i])) {
			return int(numStr[i] - '0')
		}
	}
	return 1
}

// ParseCredits 解析学分字符串为浮点数
func ParseCredits(creditStr string) float32 {
	credit, err := strconv.ParseFloat(creditStr, 32)
	if err != nil {
		fmt.Println(err)
		return 0.0
	}

	return float32(credit)
}

func formatWeekRange(weeks []int) string {
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

// ParseTimeSlots 解析时间段信息
func ParseTimeSlots(weekAndTime uint32, dayOfWeek int) []vo.TimeSlotVO {
	// 从二进制数据中提取周次和节次
	weeks, lessons := generator.Bin2WeekLesson(weekAndTime)

	if len(lessons) == 0 {
		return []vo.TimeSlotVO{}
	}

	// 格式化周次范围
	weekRange := formatWeekRange(weeks)

	// 格式化节次范围
	startPeriod := lessons[0]
	endPeriod := lessons[len(lessons)-1]

	return []vo.TimeSlotVO{
		{
			DayOfWeek:   dayOfWeek,
			StartPeriod: startPeriod,
			EndPeriod:   endPeriod,
			Weeks:       weekRange,
		},
	}
}

var bcReg = regexp.MustCompile(`[A-Z]+\d*|\d+号楼`)

func ExtractBuildingCode(building string) string {
	// 尝试提取英文字母或数字编号
	match := bcReg.FindString(building)
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
