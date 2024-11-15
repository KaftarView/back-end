package entities

import (
	"first-project/src/enums"
	"time"

	"gorm.io/gorm"
)

type News struct {
	gorm.Model
	Title       string             `gorm:"type:varchar(255);not null"`
	Description string             `gorm:"type:text"`
	Content     string             `gorm:"type:text;not null"`
	ImageURL    string             `gorm:"type:varchar(255)"`
	Category    enums.CategoryType `gorm:"type:tinyint;not null;index"`
	Author      string             `gorm:"type:varchar(100);not null"`
	PublishedAt time.Time          `gorm:"type:timestamp;not null;index"`
}
