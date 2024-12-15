package entities

import "gorm.io/gorm"

type Episode struct {
	gorm.Model
	Name        string  `gorm:"type:varchar(50);not null"`
	Description string  `gorm:"type:text"`
	BannerPath  string  `gorm:"type:text"`
	AudioPath   string  `gorm:"type:text"`
	PublisherID uint    `gorm:"not null;index"`
	Publisher   User    `gorm:"foreignKey:PublisherID"`
	PodcastID   uint    `gorm:"not null"`
	Podcast     Podcast `gorm:"foreignKey:PodcastID;"`
}
