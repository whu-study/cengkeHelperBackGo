package dto

import "time"

// 学部信息表
type Division struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null;comment:学部名称" json:"name"`
	Buildings []BuildingInfo `gorm:"foreignKey:DivisionID" json:"buildings"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
}
