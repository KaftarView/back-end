package entities

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Status      string    `gorm:"type:enum('draft','published','cancelled','completed');default:'draft'"`
	Category    string    `gorm:"type:varchar(50);index"`
	Description string    `gorm:"type:text"`
	FromDate    time.Time `gorm:"not null;index"`
	ToDate      time.Time `gorm:"not null"`
	MinCapacity uint      `gorm:"not null"`
	MaxCapacity uint      `gorm:"not null"`
	VenueType   string    `gorm:"type:enum('online','physical','hybrid');not null"`
	Location    string    `gorm:"type:text"`

	Organizers     []Organizer     `gorm:"many2many:event_organizers"`
	Attendees      []User          `gorm:"many2many:event_attendees"`
	Tickets        []Ticket        `gorm:"has_many:event_tickets"`
	Discounts      []Discount      `gorm:"has_many:event_discounts"`
	Comments       []Comment       `gorm:"has_many:event_comments"`
	Communications []Communication `gorm:"has_many:event_communications"`
}
