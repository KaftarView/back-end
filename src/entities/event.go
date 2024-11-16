package entities

import (
	"first-project/src/enums"
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name        string            `gorm:"type:varchar(50);not null"`
	Status      enums.EventStatus ``
	Category    string            `gorm:"type:varchar(50);index"`
	Description string            `gorm:"type:text"`
	FromDate    time.Time         `gorm:"not null;index"`
	ToDate      time.Time         `gorm:"not null"`
	MinCapacity uint              `gorm:"not null"`
	MaxCapacity uint              `gorm:"not null"`
	VenueType   enums.EventVenue  ``
	Location    string            `gorm:"type:text"`

	Organizers     []Organizer     `gorm:"many2many:event_organizers"`
	Attendees      []User          `gorm:"many2many:event_attendees"`
	Tickets        []Ticket        `gorm:"has_many:event_tickets"`
	Discounts      []Discount      `gorm:"has_many:event_discounts"`
	Comments       []Comment       `gorm:"has_many:event_comments"`
	Communications []Communication `gorm:"has_many:event_communications"`
}
