package entities

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	AuthorID    uint   `gorm:"not null;index"`
	Author      User   `gorm:"foreignKey:AuthorID"`
	Content     string `gorm:"type:text;not null"`
	IsModerated bool   `gorm:"default:false"`

	CommentableID uint        `gorm:"not null"`
	Commentable   Commentable `gorm:"foreignKey:CommentableID"`
}
