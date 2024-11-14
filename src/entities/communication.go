package entities

import (
	"time"

	"gorm.io/gorm"
)

type Communication struct {
	gorm.Model
	EventID    uint      `gorm:"not null;index"`
	Type       string    `gorm:"type:enum('announcement','update','reminder');not null"`
	Title      string    `gorm:"type:varchar(200);not null"`
	Content    string    `gorm:"type:text;not null"`
	SentAt     time.Time `gorm:"not null"`
	SendToAll  bool      `gorm:"default:true"`
	Recipients []User    `gorm:"many2many:communication_recipients"`
}
