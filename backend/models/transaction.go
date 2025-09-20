package models

import (
	"time"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Description     string
	Amount          float64
	Type            string
	TransactionDate time.Time
	UserID          uint
	CategoryID      uint
}