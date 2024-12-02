package entities

import (
	"first-project/src/enums"

	"gorm.io/gorm"
)

type News struct {
	gorm.Model
	Title       string             `gorm:"type:varchar(255);not null"`
	Description string             `gorm:"type:text" `
	Content     string             `gorm:"type:text;not null" `
	Category    enums.CategoryType `gorm:"type:tinyint;not null;index" `
	Author      string             `gorm:"type:varchar(100);not null" `
}

type NewsDTO struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Content     string `json:"content" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Author      string `json:"author" binding:"required"`
}
