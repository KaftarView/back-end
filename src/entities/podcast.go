package entities

import "gorm.io/gorm"

type Podcast struct {
	gorm.Model
	ID          uint       `gorm:"primarykey"`
	Name        string     `gorm:"type:varchar(50);not null"`
	Description string     `gorm:"type:text"`
	BannerPath  string     `gorm:"type:text"`
	PublisherID uint       `gorm:"not null;index"`
	Publisher   User       `gorm:"foreignKey:PublisherID"`
	Episodes    []Episode  `gorm:"foreignKey:PodcastID"`
	Categories  []Category `gorm:"many2many:event_categories"`
	Subscribers []User     `gorm:"many2many:podcast_subscribers;"`
}
