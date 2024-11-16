package dto

import "time"

type CreateEventDetails struct {
	Name        string
	Status      string
	Category    string
	Description string
	FromDate    time.Time
	ToDate      time.Time
	MinCapacity uint
	MaxCapacity uint
	VenueType   string
	Location    string
}
