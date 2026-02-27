package repository

import (
	"myproject/config"
	models "myproject/internal/model"
)

// CreateTask 创建任务
func CreateTask(task *models.Task) error {
	return config.DB.Create(task).Error
}

// GetTasksByUserAndDate 根据用户和日期获取任务列表
func GetTasksByUserAndDate(userID uint, date string) ([]models.Task, error) {
	var tasks []models.Task
	if err := config.DB.
		Where("user_id = ? AND DATE(created_at) = ?", userID, date).
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetTaskByID 获取单个任务
func GetTaskByID(id string, userID uint) (*models.Task, error) {
	var task models.Task
	if err := config.DB.
		Where("id = ? AND user_id = ?", id, userID).
		First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// UpdateTask 更新任务
func UpdateTask(task *models.Task, updates map[string]interface{}) error {
	return config.DB.Model(task).Updates(updates).Error
}

// DeleteTask 删除任务
func DeleteTask(task *models.Task) error {
	return config.DB.Delete(task).Error
}