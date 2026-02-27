package controllers

import (
	"myproject/config"
	"myproject/models"
	"myproject/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
    Username string `json:"username" binding:"required,min=3"`
    Password string `json:"password" binding:"required,min=6"`
    Email    string `json:"email" binding:"required,email"`
}

type LoginInput struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

// Register 用户注册
// @Summary      用户注册
// @Description  创建新用户账号
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        request body RegisterInput true "注册信息"
// @Success      200  {object}  map[string]interface{}  "注册成功"
// @Failure      400  {object}  map[string]interface{}  "请求参数错误"
// @Failure      500  {object}  map[string]interface{}  "服务器错误"
// @Router       /register [post]
func Register(c *gin.Context) {
    var input RegisterInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // 创建用户
    user := models.User{
        Username: input.Username,
        Password: input.Password,
        Email:    input.Email,
    }
    
    // 密码加密
    if err := user.HashPassword(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
        return
    }
    
    // 保存到数据库
    if err := config.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或邮箱已存在"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "注册成功",
        "user": gin.H{
            "id": user.ID,
            "username": user.Username,
            "email": user.Email,
        },
    })
}

// Login 用户登录
// @Summary      用户登录
// @Description  用户登录获取JWT令牌
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        request body LoginInput true "登录信息"
// @Success      200  {object}  map[string]interface{}  "登录成功返回token"
// @Failure      400  {object}  map[string]interface{}  "请求参数错误"
// @Failure      401  {object}  map[string]interface{}  "用户名或密码错误"
// @Router       /login [post]
func Login(c *gin.Context) {
    var input LoginInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // 查找用户
    var user models.User
    if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
        return
    }
    
    // 验证密码
    if !user.CheckPassword(input.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
        return
    }
    
    // 生成token
    token, err := utils.GenerateToken(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "登录成功",
        "token": token,
        "user": gin.H{
            "id": user.ID,
            "username": user.Username,
            "email": user.Email,
        },
    })
}