package router

import (
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
