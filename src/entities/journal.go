package entities

import "gorm.io/gorm"

type Journal struct {
	gorm.Model
	Name            string `gorm:"type:varchar(50);not null"`
	Description     string `gorm:"type:text"`
	BannerPath      string `gorm:"type:text"`
	JournalFilePath string `gorm:"type:text"`
}
