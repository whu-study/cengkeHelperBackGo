package vo

// DivisionVO 学部信息VO（匹配前端API设计）
type DivisionVO struct {
	DivisionID     string       `json:"divisionId"`
	DivisionName   string       `json:"divisionName"`
	Description    string       `json:"description"`
	TotalBuildings int          `json:"totalBuildings"`
	TotalFloors    int          `json:"totalFloors"`
	TotalCourses   int          `json:"totalCourses"`
	Buildings      []BuildingVO `json:"buildings"`
}

// BuildingVO 教学楼信息VO（匹配前端API设计）
type BuildingVO struct {
	BuildingID   string    `json:"buildingId"`
	BuildingName string    `json:"buildingName"`
	BuildingCode string    `json:"buildingCode"`
	Address      string    `json:"address,omitempty"`
	Description  string    `json:"description,omitempty"`
	TotalFloors  int       `json:"totalFloors"`
	TotalRooms   int       `json:"totalRooms"`
	TotalCourses int       `json:"totalCourses"`
	Floors       []FloorVO `json:"floors"`
}

// FloorVO 楼层信息VO（匹配前端API设计）
type FloorVO struct {
	FloorID     string         `json:"floorId"`
	FloorName   string         `json:"floorName"`
	FloorNumber int            `json:"floorNumber"`
	Description string         `json:"description,omitempty"`
	Rooms       []RoomVO       `json:"rooms"`
	Courses     []CourseInfoVO `json:"courses"`
}

// RoomVO 教室信息VO（匹配前端API设计）
type RoomVO struct {
	RoomID     string   `json:"roomId"`
	RoomNumber string   `json:"roomNumber"`
	RoomName   string   `json:"roomName,omitempty"`
	Capacity   int      `json:"capacity,omitempty"`
	RoomType   string   `json:"roomType,omitempty"`
	Facilities []string `json:"facilities,omitempty"`
}

// TimeSlotVO 时间段信息VO（匹配前端API设计）
type TimeSlotVO struct {
	DayOfWeek   int    `json:"dayOfWeek"`   // 星期几 (1-7)
	StartPeriod int    `json:"startPeriod"` // 开始节次
	EndPeriod   int    `json:"endPeriod"`   // 结束节次
	Weeks       string `json:"weeks"`       // 周次范围，如 "1-16周"
}

// CurrentCourseTimeVO 当前课程时间信息VO
type CurrentCourseTimeVO struct {
	WeekNum      int    `json:"weekNum"`      // 当前周次
	Weekday      int    `json:"weekday"`      // 星期几 (0-6)
	WeekdayName  string `json:"weekdayName"`  // 星期名称（如"周一"）
	LessonNum    int    `json:"lessonNum"`    // 当前节次（-1表示非上课时间）
	LessonStatus string `json:"lessonStatus"` // 节次状态描述（如"第3节"或"非上课时间"）
	Timestamp    int64  `json:"timestamp"`    // 当前时间戳
}
