package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Username  string    `json:"username" gorm:"unique;not null"`
    Password  string    `json:"-"`  // 不返回给前端
    Email     string    `json:"email" gorm:"unique;not null"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// 密码加密
func (u *User) HashPassword() error {
    hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashed)
    return nil
}

// 验证密码
func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    return err == nil
}