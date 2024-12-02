package entities

import (
	"gorm.io/gorm"
)

type News struct {
	gorm.Model
	Title       string     `gorm:"type:varchar(255);not null"`
	Description string     `gorm:"type:text" `
	Content     string     `gorm:"type:text;not null" `
	Author      string     `gorm:"type:varchar(100);not null" `
	Categories  []Category `gorm:"many2many:news_categories"`
	BannerPaths []string   `gorm:"type:text"`
}

type NewsDTO struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Content     string     `json:"content" binding:"required"`
	Author      string     `json:"author" binding:"required"`
	Categories  []Category `gorm:"many2many:news_categories"`
}
