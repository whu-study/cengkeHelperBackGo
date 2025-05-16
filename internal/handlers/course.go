package handlers

import (
	"cengkeHelperBackGo/internal/config"
	database "cengkeHelperBackGo/internal/db"
	"cengkeHelperBackGo/internal/models/dto"
	"cengkeHelperBackGo/internal/models/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetCoursesHandler(c *gin.Context) {
	var divisions []dto.Division

	if err := database.Client.
		Preload("Buildings.Courses").
		Find(&divisions).
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.RespData{
			Code: config.CodeDatabaseError,
			Msg:  "数据库查询失败",
		})
		return
	}

	// 转换为前端需要的格式
	var result [][]vo.BuildingInfoVO
	for _, div := range divisions {
		var buildings []vo.BuildingInfoVO
		for _, b := range div.Buildings {
			buildings = append(buildings, ToVO(b))
		}
		result = append(result, buildings)
	}

	c.JSON(http.StatusOK, vo.NewSuccessResp("课程数据获取成功", result))
}
func convertCoursesToVO(courses []dto.CourseInfo) []vo.CourseInfoVO {
	var voCourses []vo.CourseInfoVO
	for _, c := range courses {
		voCourses = append(voCourses, vo.CourseInfoVO{
			Room:         c.Room,
			Faculty:      c.Faculty,
			CourseName:   c.CourseName,
			TeacherName:  c.TeacherName,
			TeacherTitle: c.TeacherTitle,
			CourseTime:   c.CourseTime,
			CourseType:   c.CourseType,
		})
	}
	return voCourses
}
func ToVO(b dto.BuildingInfo) vo.BuildingInfoVO {
	return vo.BuildingInfoVO{
		Building: b.Name,
		Label:    b.Label,
		Value:    b.Value,
		Infos:    convertCoursesToVO(b.Courses),
	}
}
