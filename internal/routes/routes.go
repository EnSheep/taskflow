package routes

import (
	"myproject/internal/handler"
	"myproject/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes 配置所有路由
	func SetupRoutes(r *gin.Engine) {
		// Swagger 文档
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// 跨域配置
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "http://127.0.0.1:5173"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: true,
		}))

		// 公开路由
		SetupPublicRoutes(r)
		
		// 需要认证的路由
		SetupPrivateRoutes(r)
	}

	func SetupPublicRoutes(r *gin.Engine) {
		// 添加路由
		r.POST("/register", handler.Register)
		r.POST("/login", handler.Login)
	}

	func SetupPrivateRoutes(r *gin.Engine) {
		// 添加路由
		tasks := r.Group("/tasks").Use(middleware.AuthMiddleware())
		{
			tasks.GET("", handler.GetTasks)
			tasks.POST("", handler.CreateTask)
			// tasks.GET("/:id", handler.GetTask)
			tasks.PUT("/:id", handler.UpdateTask)
			tasks.DELETE("/:id", handler.DeleteTask)
		}
	}
