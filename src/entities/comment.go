package entities

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	EventID uint   `gorm:"not null;index"`
	UserID  uint   `gorm:"not null;index"`
	Content string `gorm:"type:text;not null"`
	// ParentID *uint  // For nested comments
	// IsModerated bool   `gorm:"default:false"`
}
