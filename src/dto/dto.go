package dto

import "time"

type CreateEventDetails struct {
	Name        string
	Status      string
	Categories  []string
	Description string
	FromDate    time.Time
	ToDate      time.Time
	MinCapacity uint
	MaxCapacity uint
	VenueType   string
	Location    string
}

type CreateTicketDetails struct {
	Name           string
	Description    string
	Price          float64
	Quantity       uint
	SoldCount      uint
	IsAvailable    bool
	AvailableFrom  time.Time
	AvailableUntil time.Time
	EventID        uint
}
