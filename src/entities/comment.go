package entities

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	UserID      uint   `gorm:"not null;index"`
	Content     string `gorm:"type:text;not null"`
	IsModerated bool   `gorm:"default:false"`
	// ParentID *uint  // use nested set model

	CommentableID uint        `gorm:"not null"`
	Commentable   Commentable `gorm:"foreignKey:CommentableID"`
}
