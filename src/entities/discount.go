package entities

import (
	"first-project/src/enums"
	"time"

	"gorm.io/gorm"
)

type Discount struct {
	gorm.Model
	Code       string             `gorm:"type:varchar(50);uniqueIndex"`
	Type       enums.DiscountType `gorm:"type:int;not null"`
	Value      float64            `gorm:"not null"`
	ValidFrom  time.Time          `gorm:"not null"`
	ValidUntil time.Time          `gorm:"not null"`
	Quantity   uint               `gorm:"not null"`
	UsedCount  uint               `gorm:"default:0"`
	MinTickets uint               `gorm:"default:1"`
	EventID    uint               `gorm:"not null"`
	Event      Event              `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
