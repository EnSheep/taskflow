package main

import (
	"myproject/config"
	"myproject/controllers"
	_ "myproject/docs"
	"myproject/middleware"
	"myproject/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           TaskFlow API
// @version         1.0
// @description     任务管理API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization




func main() {
	// 连接数据库
	config.ConnectDB()

	// 自动迁移
	config.DB.AutoMigrate(&models.User{}, &models.Task{})

	// 创建路由
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 允许前端跨域（Vite 开发服务器等）
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// 公开路由
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// 需要认证的任务路由
	tasks := r.Group("/tasks").Use(middleware.AuthMiddleware())
	{
		tasks.GET("", controllers.GetTasks)
		tasks.POST("", controllers.CreateTask)
		// tasks.GET("/:id", controllers.GetTask)
		tasks.PUT("/:id", controllers.UpdateTask)
		tasks.DELETE("/:id", controllers.DeleteTask)
	}

	// 启动服务
	r.Run(":8080")
}