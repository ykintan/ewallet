package entity

import (
	"time"
)

type Transaction struct {
	TransactionID   uint      `gorm:"primaryKey;autoIncrement"`
	WalletID        int       `gorm:"not null;index"`
	Amount          float64   `gorm:"type:decimal(10,2);not null"`
	TransactionType string    `gorm:"type:varchar(20);not null"`
	CreatedAt       time.Time `gorm:"default:current_timestamp"`
	WalletIDSource  int       `gorm:"column:wallet_id_source"`
}
