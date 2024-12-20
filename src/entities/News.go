package entities

import (
	"gorm.io/gorm"
)

type News struct {
	gorm.Model
	ID          uint        `gorm:"primarykey"`
	Title       string      `gorm:"type:varchar(50);not null"`
	Description string      `gorm:"type:text" `
	Content     string      `gorm:"type:text;not null" `
	Content2    string      `gorm:"type:text"`
	BannerPath  string      `gorm:"type:text"`
	Banner2Path string      `gorm:"type:text"`
	AuthorID    uint        `gorm:"not null;index"`
	Author      User        `gorm:"foreignKey:AuthorID"`
	Commentable Commentable `gorm:"foreignKey:ID;constraint:OnDelete:CASCADE"`
	Categories  []Category  `gorm:"many2many:news_categories;constraint:OnDelete:CASCADE"`
}
