package handlers

import (
	"cengkeHelperBackGo/internal/models/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminEchoHandler(c *gin.Context) {
	value, _ := c.Get("username")

	c.JSON(http.StatusOK, vo.RespData{
		Code: 200,
		Data: "admin: " + value.(string),
		Msg:  "success",
	})
}
