package entities

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	EventID     uint    `gorm:"not null;index"`
	Name        string  `gorm:"type:varchar(100);not null"`
	Description string  `gorm:"type:text"`
	Price       float64 `gorm:"type:decimal(10,2);not null"`
	Quantity    uint    `gorm:"not null"`
	SoldCount   uint    `gorm:"default:0"`
	IsAvailable bool    `gorm:"default:true"`
	ValidFrom   *time.Time
	ValidUntil  *time.Time
}
