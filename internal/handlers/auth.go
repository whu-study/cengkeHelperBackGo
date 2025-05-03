package handlers

import (
	"cengkeHelperBackGo/internal/config"
	"cengkeHelperBackGo/internal/models/dto"
	"cengkeHelperBackGo/internal/models/vo"
	"cengkeHelperBackGo/internal/services"
	"cengkeHelperBackGo/pkg/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func UserLoginHandler(c *gin.Context) {

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, vo.RespData{Code: 400, Msg: err.Error()})
		return
	}

	var user dto.User
	var ok bool
	if user, ok = services.CheckUser(req.Username, req.Password); !ok {
		c.JSON(http.StatusBadRequest, vo.NewBadResp("用户名或密码错误"))
		return
	}

	// 获取到后，发布token
	// 5 天有效期，过期前端自动续签
	expirationTime := time.Now().Add(5 * 24 * time.Hour)
	token, err := utils.GenerateUserToken(user.Username, user.UserRole)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewBadResp("生成令牌失败，请联系管理员"))
		return
	}

	log.Println("用户登录: ", user.Username, user.UserRole)

	c.JSON(http.StatusOK, vo.RespData{
		Code: config.CodeSuccess,
		Data: gin.H{
			"token":        token,
			"expirationAt": expirationTime.Unix(),
		},
		Msg: "登录鉴权成功！",
	})

}
