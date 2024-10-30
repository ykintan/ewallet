package models

import "time"

type User struct {
	UserID    int32     `gorm:"primaryKey;column:user_id"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Email     string    `gorm:"unique;not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
