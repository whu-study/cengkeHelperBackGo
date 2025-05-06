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
			buildings = append(buildings, b.ToVO())
		}
		result = append(result, buildings)
	}

	c.JSON(http.StatusOK, vo.NewSuccessResp(result, "课程数据获取成功"))
}
