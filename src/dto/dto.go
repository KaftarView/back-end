package dto

import "time"

type CreateEventDetails struct {
	Name        string
	Status      string
	Categories  []string
	Description string
	BasePrice   float64
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

type CreateDiscountDetails struct {
	Code       string
	Type       string
	Value      float64
	ValidFrom  time.Time
	ValidUntil time.Time
	Quantity   uint
	UsedCount  uint
	MinTickets uint
	EventID    uint
}

type EventDetailsResponse struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	BasePrice   float64   `json:"base-price"`
	MinCapacity uint      `json:"min_capacity"`
	MaxCapacity uint      `json:"max_capacity"`
	FromDate    time.Time `json:"from_date"`
	ToDate      time.Time `json:"to_date"`
	VenueType   string    `json:"venue_type"`
	Categories  []string  `json:"categories"`
	Location    string    `json:"location"`
	Banner      string    `json:"banner"`
}

type TicketDetailsResponse struct {
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Price          float64   `json:"price"`
	Quantity       uint      `json:"quantity"`
	IsAvailable    bool      `json:"is_available"`
	AvailableFrom  time.Time `json:"available_from"`
	AvailableUntil time.Time `json:"available_until"`
}

type DiscountDetailsResponse struct {
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	Code           string    `json:"code"`
	Type           string    `json:"type"`
	Value          float64   `json:"value"`
	AvailableFrom  time.Time `json:"available_from"`
	AvailableUntil time.Time `json:"available_until"`
	Quantity       uint      `json:"quantity"`
	UsedCount      uint      `json:"used_count"`
	MinTickets     uint      `json:"min_tickets"`
}
