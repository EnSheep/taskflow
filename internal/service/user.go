package service

import (
	"errors"

	models "myproject/internal/model"
	"myproject/internal/repository"
	"myproject/utils"
)

var (
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrHashPasswordFailed  = errors.New("hash password failed")
	ErrGenerateTokenFailed = errors.New("generate token failed")
)

// Register 用户注册业务
func Register(username, password, email string) (*models.User, error) {
	user := models.User{
		Username: username,
		Password: password,
		Email:    email,
	}

	// 密码加密
	if err := user.HashPassword(); err != nil {
		return nil, ErrHashPasswordFailed
	}

	// 保存到数据库
	if err := repository.CreateUser(&user); err != nil {
		// 这里简化处理为用户已存在，后续可以根据错误类型更精细区分
		return nil, ErrUserAlreadyExists
	}

	return &user, nil
}

// Login 用户登录业务
func Login(username, password string) (string, *models.User, error) {
	// 查找用户
	user, err := repository.GetUserByUsername(username)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	// 校验密码
	if !user.CheckPassword(password) {
		return "", nil, ErrInvalidCredentials
	}

	// 生成 token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", nil, ErrGenerateTokenFailed
	}

	return token, user, nil
}