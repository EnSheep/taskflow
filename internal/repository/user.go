package repository

import (
	"myproject/config"
	models "myproject/internal/model"
)

// CreateUser 创建用户
func CreateUser(user *models.User) error {
	return config.DB.Create(user).Error
}

// GetUserByUsername 根据用户名查询用户
func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
