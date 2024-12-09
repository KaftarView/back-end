package entities

import (
	"first-project/src/enums"

	"gorm.io/gorm"
)

type News struct {
	gorm.Model
	Title       string     `gorm:"type:varchar(255);not null"`
	Description string     `gorm:"type:text" `
	Content     string     `gorm:"type:text;not null" `
	Content2    string     `gorm:"type:text"`
	Author      string     `gorm:"type:varchar(100);not null" `
	Categories  []Category `gorm:"many2many:news_categories"`
	BannerPaths string     `gorm:"type:text"`
}

type NewsDTO struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Content     string     `json:"content" binding:"required"`
	Content2    string     `json:"content2"`
	Author      string     `json:"author" binding:"required"`
	Categories  []Category `gorm:"many2many:news_categories"`
}

type NewsResponse struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Content     string     `json:"content" binding:"required"`
	Content2    string     `json:"content2"`
	Author      string     `json:"author" binding:"required"`
	Categories  []Category `gorm:"many2many:news_categories"`
	BannerPath  string     `json:"bannerpath1"`
	BannerPath2 string     `json:"bannerpath2"`
}
