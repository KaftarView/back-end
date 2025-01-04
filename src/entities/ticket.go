package entities

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	Name           string    `gorm:"type:varchar(50);not null"`
	Description    string    `gorm:"type:text"`
	Price          float64   `gorm:"type:decimal(10,2);not null"`
	Quantity       uint      `gorm:"not null"`
	SoldCount      uint      `gorm:"default:0"`
	IsAvailable    bool      `gorm:"default:true"`
	AvailableFrom  time.Time `gorm:"not null"`
	AvailableUntil time.Time `gorm:"not null"`
	// Purchasable    Purchasable `gorm:"foreignKey:ID"`
	EventID uint  `gorm:"not null"`
	Event   Event `gorm:"foreignKey:EventID"`
}
