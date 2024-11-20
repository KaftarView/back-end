package entities

import (
	"time"

	"gorm.io/gorm"
)

type Discount struct {
	gorm.Model
	Code       string    `gorm:"type:varchar(50);uniqueIndex"`
	Type       string    `gorm:"type:enum('percentage','fixed');not null"`
	Value      float64   `gorm:"not null"`
	ValidFrom  time.Time `gorm:"not null"`
	ValidUntil time.Time `gorm:"not null"`
	MaxUsage   uint      `gorm:"default:0"`
	UsedCount  uint      `gorm:"default:0"`
	MinTickets uint      `gorm:"default:1"`
	// EventID    uint      `gorm:"not null;index"`
	// Event      Event     `gorm:"-"`
	EventID uint  `gorm:"not null"`
	Event   Event `gorm:"foreignKey:EventID"`
}
