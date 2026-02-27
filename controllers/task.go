package controllers

import (
	"myproject/config"
	"myproject/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateTaskInput struct {
    // Title       string `json:"title" binding:"required"`
    Description string `json:"description"`
}

type UpdateTaskInput struct {
    // Title       string `json:"title"`
    Description string `json:"description"`
    Status      string `json:"status" binding:"omitempty,oneof=pending done"`
}

// CreateTask 创建任务
// @Summary      创建新任务
// @Description  为当前用户创建新任务
// @Tags         任务
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateTaskInput true "任务信息"
// @Success      200  {object}  map[string]interface{}  "创建成功"
// @Failure      400  {object}  map[string]interface{}  "请求参数错误"
// @Failure      401  {object}  map[string]interface{}  "未认证"
// @Router       /api/tasks [post]
func CreateTask(c *gin.Context) {
    user, _ := c.Get("user")
    currentUser := user.(models.User)
    
    var input CreateTaskInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    task := models.Task{
        Description: input.Description,
        UserID:      currentUser.ID,
    }
    
    if err := config.DB.Create(&task).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "创建成功",
        "task": task,
    })
}

// GetTasks 获取任务列表
// @Summary      获取所有任务
// @Description  获取当前用户的所有任务
// @Tags         任务
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}  "任务列表"
// @Failure      401  {object}  map[string]interface{}  "未认证"
// @Router       /api/tasks [get]
func GetTasks(c *gin.Context) {
    user, _ := c.Get("user")
    currentUser := user.(models.User)
    date := c.Query("date") // 可选的日期参数，格式为 "YYYY-MM-DD"
    if date == "" {
        date = time.Now().Format("2006-01-02")
    }

    var tasks []models.Task
    if err := config.DB.Where("user_id = ? AND DATE(created_at) = ?", currentUser.ID, date).Find(&tasks).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "tasks": tasks,
    })
}

// GetTask 获取单个任务
// @Summary      获取任务详情
// @Description  根据ID获取任务详情
// @Tags         任务
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "任务ID"
// @Success      200  {object}  map[string]interface{}  "任务详情"
// @Failure      401  {object}  map[string]interface{}  "未认证"
// @Failure      404  {object}  map[string]interface{}  "任务不存在"
// @Router       /api/tasks/{id} [get]
func GetTask(c *gin.Context) {
    user, _ := c.Get("user")
    currentUser := user.(models.User)
    
    id := c.Param("id")
    var task models.Task
    
    if err := config.DB.Where("id = ? AND user_id = ?", id, currentUser.ID).First(&task).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"task": task})
}

// UpdateTask 更新任务
// @Summary      更新任务
// @Description  更新任务信息
// @Tags         任务
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path    string           true  "任务ID"
// @Param        request body UpdateTaskInput true "更新信息"
// @Success      200     {object} map[string]interface{} "更新成功"
// @Failure      400     {object} map[string]interface{} "请求参数错误"
// @Failure      401     {object} map[string]interface{} "未认证"
// @Failure      404     {object} map[string]interface{} "任务不存在"
// @Router       /api/tasks/{id} [put]
func UpdateTask(c *gin.Context) {
    user, _ := c.Get("user")
    currentUser := user.(models.User)
    
    id := c.Param("id")
    var task models.Task
    
    if err := config.DB.Where("id = ? AND user_id = ?", id, currentUser.ID).First(&task).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
        return
    }
    
    var input UpdateTaskInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // 更新字段
    updates := make(map[string]interface{})
    // if input.Title != "" {
    //     updates["title"] = input.Title
    // }
    if input.Description != "" {
        updates["description"] = input.Description
    }
    if input.Status != "" {
        updates["status"] = input.Status
    }
    
    config.DB.Model(&task).Updates(updates)
    
    c.JSON(http.StatusOK, gin.H{
        "message": "更新成功",
        "task": task,
    })
}

// DeleteTask 删除任务
// @Summary      删除任务
// @Description  删除指定任务
// @Tags         任务
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "任务ID"
// @Success      200  {object}  map[string]interface{}  "删除成功"
// @Failure      401  {object}  map[string]interface{}  "未认证"
// @Failure      404  {object}  map[string]interface{}  "任务不存在"
// @Router       /api/tasks/{id} [delete]
func DeleteTask(c *gin.Context) {
    user, _ := c.Get("user")
    currentUser := user.(models.User)
    
    id := c.Param("id")
    var task models.Task
    
    if err := config.DB.Where("id = ? AND user_id = ?", id, currentUser.ID).First(&task).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
        return
    }
    
    config.DB.Delete(&task)
    
    c.JSON(http.StatusOK, gin.H{
        "message": "删除成功",
    })
}