package generator

import (
	"fmt"
	"log"
)

func NearestToDisplay(lessonNum int, binNum uint32) string {
	if lessonNum == -1 {
		return "全天"
	}
	begin := lessonNum
	end := lessonNum
	for (1<<(begin-1))&binNum != 0 {
		if begin <= 1 {
			begin = 0
			break
		}
		begin--
	}

	for (1<<(end-1))&binNum != 0 {
		end++
	}
	if end > 14 {
		end = 14
	}

	return fmt.Sprintf("第 %d-%d 节", begin+1, end-1)
}

// IsWeekLessonMatch 判断周次和节次是否在所给的二进制数中
func IsWeekLessonMatch(weekNum, lessonNum int, binNum uint32) bool {
	if lessonNum < -1 {
		return false
	}

	if weekNum == -1 && lessonNum == -1 {
		return true
	}

	if weekNum == -1 {
		if (1<<(lessonNum-1))&binNum != 0 {
			return true
		}
		return false
	}
	if lessonNum == -1 {
		if (1<<(32-weekNum))&binNum != 0 {
			return true
		}
		return false
	}
	if (1<<(32-weekNum))&binNum != 0 && (1<<(lessonNum-1))&binNum != 0 {
		return true
	}
	return false
}
func Bin2WeekLesson(binNum uint32) ([]int, []int) {
	weekNums := make([]int, 0)
	lessonNums := make([]int, 0)

	for i := 1; i <= 19; i++ {
		if (1<<(32-i))&binNum == 0 {
			continue
		}
		weekNums = append(weekNums, i)
	}

	for i := 1; i <= 13; i++ {
		if (1<<(i-1))&binNum == 0 {
			continue
		}
		lessonNums = append(lessonNums, i)
	}

	return weekNums, lessonNums

}

func WeekLesson2Bin(weekNums, lessonNums []int) uint32 {
	var res uint32 = 0
	for _, num := range weekNums {
		if num > 19 {
			log.Fatal("weekNum不能超过19", num)
		}
		if res&((1<<31)>>(num-1)) != 0 {
			log.Println("weekNum重复覆盖！", weekNums)
		}
		res = res | ((1 << 31) >> (num - 1))
	}

	for _, num := range lessonNums {
		if num > 13 {
			log.Fatal("lessonNums不能超过13", num)
		}
		if res&(1<<(num-1)) != 0 {
			log.Println("lessonNum重复覆盖！", lessonNums)
		}
		res = res | (1 << (num - 1))
	}

	return res
}
