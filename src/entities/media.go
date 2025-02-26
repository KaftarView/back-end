package entities

import (
	"gorm.io/gorm"
)

type Media struct {
	gorm.Model
	Name    string `gorm:"type:varchar(50);not null"`
	Size    int64  `gorm:"not null"`
	Type    string
	Path    string `gorm:"type:text;not null"`
	EventID uint   `gorm:"not null;index"`
	Event   Event  `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
