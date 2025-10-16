package dto

import "time"

// Division 学部信息表
type Division struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	DivisionID  string    `gorm:"type:varchar(50);uniqueIndex;not null;comment:学部唯一标识" json:"divisionId"`
	Name        string    `gorm:"type:varchar(100);not null;comment:学部名称" json:"name"`
	Description string    `gorm:"type:text;comment:学部描述" json:"description"`
	Icon        string    `gorm:"type:varchar(255);comment:学部图标URL" json:"icon"`
	SortOrder   int       `gorm:"type:int;default:0;comment:排序值" json:"sortOrder"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName 自定义表名
func (Division) TableName() string {
	return "divisions"
}
