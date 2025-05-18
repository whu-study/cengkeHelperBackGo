package vo

// 教学楼信息VO（完全匹配前端接口）
type BuildingInfoVO struct {
	Building string         `json:"building"` // 对应前端 building 字段
	Label    string         `json:"label"`
	Value    int            `json:"value"`
	Infos    []CourseInfoVO `json:"infos"` // 对应前端 infos 字段
}
