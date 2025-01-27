package entities

import (
	"first-project/src/enums"
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	ID             uint                   `gorm:"primarykey"`
	Name           string                 `gorm:"type:varchar(50);not null"`
	Status         enums.EventStatus      `gorm:"type:int;not null"`
	Description    string                 `gorm:"type:text"`
	BasePrice      float64                `gorm:"type:decimal(10,2)"`
	FromDate       time.Time              `gorm:"not null;index"`
	ToDate         time.Time              `gorm:"not null"`
	MinCapacity    uint                   `gorm:"not null"`
	MaxCapacity    uint                   `gorm:"not null"`
	VenueType      enums.EventVenue       `gorm:"type:int;not null"`
	Location       string                 `gorm:"type:text"`
	BannerPath     string                 `gorm:"type:text"`
	Communications map[string]interface{} `gorm:"-"`
	Commentable    Commentable            `gorm:"foreignKey:ID;constraint:OnDelete:CASCADE"`
	Tickets        []Ticket               `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
	Discounts      []Discount             `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
	Media          []Media                `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
	Organizers     []Organizer            `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
	Categories     []Category             `gorm:"many2many:event_categories;constraint:OnDelete:CASCADE"`
}
