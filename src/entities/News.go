package entities

import (
	"first-project/src/enums"
	"time"

	"gorm.io/gorm"
)

type News struct {
	gorm.Model
	Title       string             `gorm:"type:varchar(255);not null" json:"title"`
	Description string             `gorm:"type:text" json:"description"`
	Content     string             `gorm:"type:text;not null" json:"content"`
	ImageURL    string             `gorm:"type:varchar(255)" json:"image_url"`
	Category    enums.CategoryType `gorm:"type:tinyint;not null;index" json:"category"`
	Author      string             `gorm:"type:varchar(100);not null" json:"author"`
	PublishedAt time.Time          `gorm:"type:timestamp;not null;index" json:"published_at"`
}

type NewsDTO struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Content     string `json:"content" binding:"required"`
	ImageURL    string `json:"image_url"`
	Category    string `json:"category" binding:"required"`
	Author      string `json:"author" binding:"required"`
	PublishedAt string `json:"published_at" binding:"required"`
}
