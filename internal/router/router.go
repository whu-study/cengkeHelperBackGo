package router

import (
	"cengkeHelperBackGo/internal/filter"
	"cengkeHelperBackGo/internal/handlers"
	"cengkeHelperBackGo/internal/handlers/auth"
	"cengkeHelperBackGo/internal/handlers/course"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var app *gin.Engine

func Routers() *gin.Engine {
	postHandler := handlers.NewPostHandler()
	commentHandler := handlers.NewCommentHandler()
	courseHandler := course.NewCourseHandler()

	v1 := app.Group("/api/v1")
	{
		v1.GET("/ping", handlers.PingHandler)
		v1.POST("/auth/user-login", auth.UserLoginHandler)
		v1.POST("/auth/user-register", auth.UserRegisterHandler)
		v1.POST("/auth/user-logout", auth.UserLogoutHandler)
		v1.POST("/auth/send-email-code", auth.SendEmailCodeHandler) // 添加发送验证码接口
		v1.GET("/courses", courseHandler.GetCoursesHandler)
		v1.GET("/courses/current-time", courseHandler.GetCurrentCourseTimeHandler) // 新增：获取当前课程时间
		v1.GET("/courses/structured", courseHandler.GetStructuredCoursesHandler)   // 新增：获取结构化课程数据
		v1.GET("/courses/:courseId", courseHandler.GetCourseDetailHandler)
		v1.GET("/posts/comments/:postId", commentHandler.GetCommentsByPostID) // GET /api/v1/posts/:id/comments (获取帖子的评论)
		v1.GET("/posts", postHandler.GetPosts)
		v1.GET("/posts/active-users", postHandler.GetActiveUsersHandler)
		v1.GET("/community/stats", postHandler.GetCommunityStatsHandler)
		v1.GET("/community/overview", postHandler.GetCommunityOverviewHandler)
		v1.GET("/posts/:id", postHandler.GetPostByID)
		v1.GET("/courses/reviews/:courseId", courseHandler.GetCourseReviewsHandler)
		v1.Use(filter.UserAuthChecker())

		v1.GET("/users/echo", handlers.UserEchoHandler)
		v1.GET("/users/profile", handlers.UserProfileHandler)
		v1.PUT("/users/profile", handlers.UpdateUserProfileHandler)
		courses := v1.Group("/courses")
		{

			courses.POST("/reviews", courseHandler.SubmitCourseReviewHandler)

		}
		posts := v1.Group("/posts") // 应用用户认证中间件
		{

			posts.POST("", postHandler.CreatePost)                           // POST /api/v1/posts (创建帖子)
			posts.PUT("/:id", postHandler.UpdatePost)                        // PUT /api/v1/posts/:id (更新帖子)
			posts.DELETE("/:id", postHandler.DeletePost)                     // DELETE /api/v1/posts/:id (删除帖子)
			posts.POST("/:id/toggle-like", postHandler.ToggleLikePost)       // POST /api/v1/posts/:id/toggle-like
			posts.POST("/:id/toggle-collect", postHandler.ToggleCollectPost) // POST /api/v1/posts/:id/toggle-collect

		} // POST /api/v1/posts/:id/toggle-collect

		comments := v1.Group("/comments")
		{
			comments.POST("", commentHandler.AddComment)                               // POST /api/v1/posts/:id/comments (创建帖子的评论)
			comments.DELETE("/:commentId", commentHandler.DeleteComment)               // DELETE /api/v1/posts/:id/comments/:comment_id (删除帖子的评论)
			comments.POST("/:commentId/toggle-like", commentHandler.ToggleLikeComment) // POST /api/v1/posts/:id/toggle-collect
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
