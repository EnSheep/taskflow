package middleware

import (
	"myproject/config"
	models "myproject/internal/model"
	"myproject/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取 Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证token"})
            c.Abort()
            return
        }
        
        // 提取 token (格式: Bearer <token>)
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "token格式错误"})
            c.Abort()
            return
        }
        
        // 解析 token
        claims, err := utils.ParseToken(parts[1])
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
            c.Abort()
            return
        }
        
        // 获取用户信息
        var user models.User
        if err := config.DB.First(&user, claims.UserID).Error; err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
            c.Abort()
            return
        }
        
        // 将用户信息存入上下文
        c.Set("user", user)
        c.Next()
    }
}