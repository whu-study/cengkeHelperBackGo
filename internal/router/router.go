package router

import (
	"cengkeHelperBackGo/internal/filter"
	"cengkeHelperBackGo/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

var app *gin.Engine

func Routers() *gin.Engine {
	v1 := app.Group("/api/v1")
	{
		v1.GET("/ping", handlers.PingHandler)
		v1.POST("/auth/user-login", handlers.UserLoginHandler)

		v1.Use(filter.UserAuthChecker())
		v1.GET("/users/echo", handlers.UserEchoHandler)
		v1.GET("/users/profile", handlers.UserProfileHandler)

		v1.Use(filter.AdminAuthChecker())
		v1.GET("/admins/echo", handlers.AdminEchoHandler)
	}
	return app
}

func init() {
	//gin.SetMode(gin.ReleaseMode)
	//gin.DefaultWriter = io.Discard
	app = gin.Default()

	// 中间件，解决开发时的跨域问题
	app.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))
}
