package main

import (
	"myproject/config"
	_ "myproject/docs"
	models "myproject/internal/model"
	"myproject/internal/routes"

	"github.com/gin-gonic/gin"
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

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务
	r.Run(":8080")
}