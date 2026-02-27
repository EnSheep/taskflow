package models

import "time"

type Task struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    // Title       string    `json:"title" gorm:"not null"`
    Description string    `json:"description"`
    Status      string    `json:"status" gorm:"default:pending"` // pending, done
    UserID      uint      `json:"user_id" gorm:"index"`
    User        User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}