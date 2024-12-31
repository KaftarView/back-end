package entities

import (
	"time"

	"gorm.io/gorm"
)

type Councilor struct {
	gorm.Model
	FirstName    string    `gorm:"type:varchar(50);not null"`
	LastName     string    `gorm:"type:varchar(50);not null"`
	ProfilePath  string    `gorm:"text"`
	Semester     int       `gorm:"not null"`
	Description  string    `gorm:"text"`
	PromotedDate time.Time `gorm:"not null"`
	UserID       uint      `gorm:"not null"`
	User         User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
