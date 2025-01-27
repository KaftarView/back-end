package entities

import (
	"gorm.io/gorm"
)

type Councilor struct {
	gorm.Model
	FirstName    string `gorm:"type:varchar(50);not null"`
	LastName     string `gorm:"type:varchar(50);not null"`
	ProfilePath  string `gorm:"text"`
	EnteringYear int    `gorm:"not null"`
	Description  string `gorm:"text"`
	PromotedYear int    `gorm:"not null;index"`
	UserID       uint   `gorm:"not null"`
	User         User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
