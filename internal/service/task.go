package service

import (
	"errors"
	"time"

	models "myproject/internal/model"
	"myproject/internal/repository"
)

var (
	ErrTaskNotFound   = errors.New("task not found")
	ErrCreateTaskFail = errors.New("create task failed")
	ErrQueryTaskFail  = errors.New("query task failed")
	ErrUpdateTaskFail = errors.New("update task failed")
	ErrDeleteTaskFail = errors.New("delete task failed")
)

// CreateTask 创建任务
func CreateTask(user models.User, description string) (*models.Task, error) {
	task := models.Task{
		Description: description,
		UserID:      user.ID,
	}

	if err := repository.CreateTask(&task); err != nil {
		return nil, ErrCreateTaskFail
	}

	return &task, nil
}

// GetTasks 获取任务列表（按日期）
func GetTasks(user models.User, date string) ([]models.Task, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	tasks, err := repository.GetTasksByUserAndDate(user.ID, date)
	if err != nil {
		return nil, ErrQueryTaskFail
	}
	return tasks, nil
}

// GetTask 获取单个任务
func GetTask(user models.User, id string) (*models.Task, error) {
	task, err := repository.GetTaskByID(id, user.ID)
	if err != nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

// UpdateTask 更新任务
func UpdateTask(user models.User, id string, updates map[string]interface{}) (*models.Task, error) {
	task, err := repository.GetTaskByID(id, user.ID)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	if err := repository.UpdateTask(task, updates); err != nil {
		return nil, ErrUpdateTaskFail
	}

	return task, nil
}

// DeleteTask 删除任务
func DeleteTask(user models.User, id string) error {
	task, err := repository.GetTaskByID(id, user.ID)
	if err != nil {
		return ErrTaskNotFound
	}

	if err := repository.DeleteTask(task); err != nil {
		return ErrDeleteTaskFail
	}

	return nil
}