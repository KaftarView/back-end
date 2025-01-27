package entities

import (
	"gorm.io/gorm"
)

type Organizer struct {
	gorm.Model
	Name        string `gorm:"type:varchar(50);not null"`
	Email       string `gorm:"type:varchar(50);not null"`
	ProfilePath string `gorm:"type:text"`
	Description string `gorm:"type:text"`
	EventID     uint   `gorm:"not null;index"`
	Event       Event  `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
