package entity

import (
	"time"
)

type Wallet struct {
	Walletid  int32     `gorm:"primaryKey;column:wallet_id"`
	UserID    uint      `gorm:"not null"`
	Balance   float64   `gorm:"default:0.00"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
