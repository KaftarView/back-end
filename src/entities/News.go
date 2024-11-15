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
