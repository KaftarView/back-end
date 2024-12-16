package entities

import "gorm.io/gorm"

type Podcast struct {
	gorm.Model
	ID          uint        `gorm:"primarykey"`
	Name        string      `gorm:"type:varchar(50);not null"`
	Description string      `gorm:"type:text"`
	BannerPath  string      `gorm:"type:text"`
	Commentable Commentable `gorm:"foreignKey:ID;constraint:OnDelete:CASCADE"`
	PublisherID uint        `gorm:"not null;index"`
	Publisher   User        `gorm:"foreignKey:PublisherID"`
	Episodes    []Episode   `gorm:"foreignKey:PodcastID;constraint:OnDelete:CASCADE"`
	Categories  []Category  `gorm:"many2many:podcast_categories;constraint:OnDelete:CASCADE"`
	Subscribers []User      `gorm:"many2many:podcast_subscribers;"`
}
