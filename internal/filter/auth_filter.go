package filter

import (
	"cengkeHelperBackGo/internal/config"
	"cengkeHelperBackGo/internal/models"
	"cengkeHelperBackGo/internal/models/dto"
	"cengkeHelperBackGo/internal/models/vo"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func UserAuthChecker() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 获取请求头中的 token
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, vo.RespData{
				Code: 401,
				Msg:  "缺少token",
			})

			// 拦截无token的请求，只是终止了http请求的继续， 但不终止代码流程
			c.Abort()
			return
		}

		// 解析token
		claims, err := ParseUserJwt(tokenString)
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrSignatureInvalid):
				c.JSON(http.StatusUnauthorized, vo.RespData{
					Code: 401,
					Msg:  "无效的token签名",
				})
				c.Abort()
			default:
				c.JSON(http.StatusUnauthorized, vo.RespData{
					Code: 401,
					Msg:  "token无效: " + err.Error(),
				})
				c.Abort()
			}
			return

		}

		// 将数据保存到请求上下文，传递给下一级请求链
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		// 进入下一步请求链
		c.Next()
	}
}

func AdminAuthChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		if roleStr, ok := c.Get("role"); ok {
			role := roleStr.(uint8)
			switch role {
			case dto.UserRoleAdmin:
				c.Next()
			default:
				c.JSON(http.StatusForbidden, vo.RespData{
					Code: 403,
					Data: nil,
					Msg:  "没有权限访问该资源",
				})
				c.Abort()
			}
		} else {
			c.JSON(http.StatusUnauthorized, vo.RespData{
				Code: 401,
				Data: nil,
				Msg:  "没有权限访问该资源",
			})
			c.Abort()
		}
	}
}

func ParseUserJwt(token string) (*models.UserClaims, error) {
	t, err := jwt.ParseWithClaims(token, &models.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Conf.JwtSecurityKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Go 语言中的类型断言（Type Assertion）语法
	// 如果类型转换成功，则ok为true
	// 否则ok为false
	if claims, ok := t.Claims.(*models.UserClaims); ok && t.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
