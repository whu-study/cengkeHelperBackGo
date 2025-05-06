package dto

import (
	"cengkeHelperBackGo/internal/models/vo"
	"time"
)

// 教学楼信息表
type BuildingInfo struct {
	ID         uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	DivisionID uint         `gorm:"index;not null;comment:所属学部ID" json:"divisionId"`
	Name       string       `gorm:"type:varchar(255);not null;comment:教学楼名称" json:"building"` // 对应前端的building字段
	Label      string       `gorm:"type:varchar(255);comment:显示名称" json:"label,omitempty"`
	Value      int          `gorm:"type:int;default:0;comment:排序值" json:"value,omitempty"`
	Courses    []CourseInfo `gorm:"foreignKey:BuildingID" json:"infos"`
	CreatedAt  time.Time    `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time    `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (b *BuildingInfo) ToVO() vo.BuildingInfoVO {
	return vo.BuildingInfoVO{
		Building: b.Name, // 模型中的 Name 映射到前端的 building
		Label:    b.Label,
		Value:    b.Value,
		Infos:    convertCoursesToVO(b.Courses),
	}
}
