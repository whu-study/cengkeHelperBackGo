package generator

import (
	"reflect"
	"testing"
)

func TestWeekLessonRoundtrip(t *testing.T) {
	weekNums := []int{1, 2, 3, 19}
	lessonNums := []int{1, 5, 13}
	bin := WeekLesson2Bin(weekNums, lessonNums)
	w, l := Bin2WeekLesson(bin)
	if !reflect.DeepEqual(w, weekNums) {
		t.Fatalf("weeks mismatch: got %v want %v", w, weekNums)
	}
	if !reflect.DeepEqual(l, lessonNums) {
		t.Fatalf("lessons mismatch: got %v want %v", l, lessonNums)
	}
}

func TestIsWeekLessonMatchVarious(t *testing.T) {
	weekNums := []int{2, 4}
	lessonNums := []int{3, 4}
	bin := WeekLesson2Bin(weekNums, lessonNums)

	// exact matches
	for _, w := range weekNums {
		for _, l := range lessonNums {
			if !IsWeekLessonMatch(w, l, bin) {
				t.Fatalf("expected match for week %d lesson %d", w, l)
			}
		}
	}

	// wildcard week
	for _, l := range lessonNums {
		if !IsWeekLessonMatch(-1, l, bin) {
			t.Fatalf("expected match for any-week lesson %d", l)
		}
	}

	// wildcard lesson
	for _, w := range weekNums {
		if !IsWeekLessonMatch(w, -1, bin) {
			t.Fatalf("expected match for week %d any-lesson", w)
		}
	}

	// both wildcards
	if !IsWeekLessonMatch(-1, -1, bin) {
		t.Fatalf("expected match for both wildcards")
	}
}

func TestNearestToDisplayRanges(t *testing.T) {
	// contiguous lessons 2-4
	bin := WeekLesson2Bin([]int{1}, []int{2, 3, 4})
	s := NearestToDisplay(3, bin)
	if s != "第 2-4 节" {
		t.Fatalf("NearestToDisplay unexpected: %s", s)
	}

	// single lesson
	bin2 := WeekLesson2Bin([]int{1}, []int{7})
	s2 := NearestToDisplay(7, bin2)
	if s2 != "第 7-7 节" {
		t.Fatalf("NearestToDisplay single unexpected: %s", s2)
	}

	// all day
	s3 := NearestToDisplay(-1, 0)
	if s3 != "全天" {
		t.Fatalf("NearestToDisplay all-day unexpected: %s", s3)
	}
}
