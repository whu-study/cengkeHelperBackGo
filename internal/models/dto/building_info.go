package dto

import (
	//"cengkeHelperBackGo/internal/models/vo"
	"time"
)

// BuildingInfo 教学楼信息表
type BuildingInfo struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BuildingID  string    `gorm:"type:varchar(50);uniqueIndex;not null;comment:教学楼唯一标识" json:"buildingId"`
	DivisionID  string    `gorm:"type:varchar(50);index;not null;comment:所属学部ID" json:"divisionId"`
	Name        string    `gorm:"type:varchar(100);not null;comment:教学楼名称" json:"name"`
	Code        string    `gorm:"type:varchar(10);not null;comment:教学楼代码" json:"code"`
	Address     string    `gorm:"type:text;comment:地址" json:"address"`
	Description string    `gorm:"type:text;comment:描述" json:"description"`
	TotalFloors int       `gorm:"type:int;default:0;comment:楼层数" json:"totalFloors"`
	SortOrder   int       `gorm:"type:int;default:0;comment:排序值" json:"sortOrder"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName 自定义表名
func (BuildingInfo) TableName() string {
	return "buildings"
}

// Floor 楼层信息表
type Floor struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FloorID     string    `gorm:"type:varchar(50);uniqueIndex;not null;comment:楼层唯一标识" json:"floorId"`
	BuildingID  string    `gorm:"type:varchar(50);index;not null;comment:所属教学楼ID" json:"buildingId"`
	Name        string    `gorm:"type:varchar(100);not null;comment:楼层名称" json:"name"`
	FloorNumber int       `gorm:"type:int;not null;comment:楼层号" json:"floorNumber"`
	Description string    `gorm:"type:text;comment:描述" json:"description"`
	SortOrder   int       `gorm:"type:int;default:0;comment:排序值" json:"sortOrder"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName 自定义表名
func (Floor) TableName() string {
	return "floors"
}

// Room 教室信息表
type Room struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RoomID     string    `gorm:"type:varchar(50);uniqueIndex;not null;comment:教室唯一标识" json:"roomId"`
	FloorID    string    `gorm:"type:varchar(50);index;not null;comment:所属楼层ID" json:"floorId"`
	RoomNumber string    `gorm:"type:varchar(20);not null;comment:教室编号" json:"roomNumber"`
	RoomName   string    `gorm:"type:varchar(100);comment:教室名称" json:"roomName"`
	Capacity   int       `gorm:"type:int;comment:容纳人数" json:"capacity"`
	RoomType   string    `gorm:"type:varchar(50);comment:教室类型" json:"roomType"`
	Facilities string    `gorm:"type:json;comment:设施设备JSON" json:"facilities"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName 自定义表名
func (Room) TableName() string {
	return "rooms"
}
