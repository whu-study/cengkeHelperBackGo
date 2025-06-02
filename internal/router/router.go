package router

import (
	"cengkeHelperBackGo/internal/filter"
	"cengkeHelperBackGo/internal/handlers"
	"cengkeHelperBackGo/internal/handlers/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

var app *gin.Engine

func Routers() *gin.Engine {
	postHandler := handlers.NewPostHandler()
	commentHandler := handlers.NewCommentHandler()
	v1 := app.Group("/api/v1")
	{
		v1.GET("/ping", handlers.PingHandler)
		v1.POST("/auth/user-login", auth.UserLoginHandler)
		v1.GET("/courses", handlers.NewCourseHandler().GetCoursesHandler)
		v1.POST("/auth/user-register", auth.UserRegisterHandler)

		v1.Use(filter.UserAuthChecker())
		v1.GET("/users/echo", handlers.UserEchoHandler)
		v1.GET("/users/profile", handlers.UserProfileHandler)
		v1.PUT("/users/profile", handlers.UpdateUserProfileHandler)

		posts := v1.Group("/posts") // 应用用户认证中间件
		{
			posts.GET("/:id", postHandler.GetPostByID)
			posts.GET("", postHandler.GetPosts)
			posts.POST("", postHandler.CreatePost)                           // POST /api/v1/posts (创建帖子)
			posts.PUT("/:id", postHandler.UpdatePost)                        // PUT /api/v1/posts/:id (更新帖子)
			posts.DELETE("/:id", postHandler.DeletePost)                     // DELETE /api/v1/posts/:id (删除帖子)
			posts.POST("/:id/toggle-like", postHandler.ToggleLikePost)       // POST /api/v1/posts/:id/toggle-like
			posts.POST("/:id/toggle-collect", postHandler.ToggleCollectPost) // POST /api/v1/posts/:id/toggle-collect

		} // POST /api/v1/posts/:id/toggle-collect
		v1.GET("/posts/comments/:postId", commentHandler.GetCommentsByPostID) // GET /api/v1/posts/:id/comments (获取帖子的评论)
		comments := v1.Group("/comments")
		{
			comments.POST("", commentHandler.AddComment)                                // POST /api/v1/posts/:id/comments (创建帖子的评论)
			comments.DELETE("/:commentId", commentHandler.DeleteComment)                // DELETE /api/v1/posts/:id/comments/:comment_id (删除帖子的评论)
			comments.POST("/:com画mentId/toggle-like", commentHandler.ToggleLikeComment) // POST /api/v1/posts/:id/toggle-collect
		} // 应用用户认证中间件

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
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS", "PUT", "PATCH"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))
}
