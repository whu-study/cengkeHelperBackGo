package generator

import (
	"fmt"
	"testing"
)

var weekNums = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 19, 18, 17}
var lessonNums = []int{3, 7, 2, 4, 13, 11, 12}

func TestNearestToDisplay(t *testing.T) {
	//fmt.Println(NearestToDisplay(3, WeekLesson2Bin(weekNums, lessonNums)))
	//fmt.Println(NearestToDisplay(7, WeekLesson2Bin(weekNums, lessonNums)))
	fmt.Println(NearestToDisplay(11, WeekLesson2Bin(weekNums, lessonNums)))
	fmt.Println(NearestToDisplay(13, WeekLesson2Bin(weekNums, lessonNums)))
}
func TestHaveClass(t *testing.T) {
	binNum := WeekLesson2Bin(weekNums, lessonNums)

	for _, num := range weekNums {
		if !IsWeekLessonMatch(num, -1, binNum) {
			t.Error("error")
		}
	}

	for _, num := range lessonNums {
		IsWeekLessonMatch(-1, num, binNum)
		if !IsWeekLessonMatch(-1, num, binNum) {
			t.Error("error")
		}
	}

	for _, num := range weekNums {
		for _, num2 := range lessonNums {
			if !IsWeekLessonMatch(num, num2, binNum) {
				t.Error("error")
			}
		}
	}
}

func TestBase(t *testing.T) {
	fmt.Println(
		Bin2WeekLesson(
			WeekLesson2Bin(weekNums, lessonNums),
		),
	)

	result1, result2 := Bin2WeekLesson(
		WeekLesson2Bin(weekNums, lessonNums),
	)

	if len(result1) != 12 || len(result2) != 7 {
		t.Error("error")
	}

	for i := 1; i <= 19; i++ {
		if res1, res2 := Bin2WeekLesson(
			WeekLesson2Bin([]int{i},
				[]int{13})); res1[0] != i || res2[0] != 13 {
			t.Error("周出错")
		}
	}

	for i := 1; i <= 13; i++ {
		if res1, res2 := Bin2WeekLesson(
			WeekLesson2Bin([]int{19},
				[]int{i})); res1[0] != 19 || res2[0] != i {
			t.Error("节出错")
		}
	}
	//t.Errorf("Add(3, 4) = %d; want 7", result)

}
