package dto

import (
	"mime/multipart"
	"time"
)

type CreateEventRequest struct {
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
	Banner      *multipart.FileHeader
}

type UpdateEventRequest struct {
	ID          uint
	Name        *string
	Status      *string
	Description *string
	BasePrice   *float64
	FromDate    *time.Time
	ToDate      *time.Time
	MinCapacity *uint
	MaxCapacity *uint
	VenueType   *string
	Location    *string
	Categories  *[]string
	Banner      *multipart.FileHeader
}

type EventDetailsResponse struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	BasePrice   float64   `json:"basePrice"`
	MinCapacity uint      `json:"minCapacity"`
	MaxCapacity uint      `json:"maxCapacity"`
	FromDate    time.Time `json:"fromDate"`
	ToDate      time.Time `json:"toDate"`
	VenueType   string    `json:"venueType"`
	Categories  []string  `json:"categories"`
	Location    string    `json:"location"`
	Banner      string    `json:"banner"`
}

type CreateTicketRequest struct {
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

type UpdateTicketRequest struct {
	Name           *string
	Description    *string
	Price          *float64
	Quantity       *uint
	SoldCount      *uint
	IsAvailable    *bool
	AvailableFrom  *time.Time
	AvailableUntil *time.Time
	TicketID       uint
}

type TicketDetailsResponse struct {
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Price          float64   `json:"price"`
	Quantity       uint      `json:"quantity"`
	RemainTickets  uint      `json:"remainTickets"`
	IsAvailable    bool      `json:"isAvailable"`
	AvailableFrom  time.Time `json:"availableFrom"`
	AvailableUntil time.Time `json:"availableUntil"`
}

type CreateDiscountRequest struct {
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

type UpdateDiscountRequest struct {
	Code           *string
	Type           *string
	Value          *float64
	AvailableFrom  *time.Time
	AvailableUntil *time.Time
	Quantity       *uint
	UsedCount      *uint
	MinTickets     *uint
	DiscountID     uint
}

type DiscountDetailsResponse struct {
	ID         uint      `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	Code       string    `json:"code"`
	Type       string    `json:"type"`
	Value      float64   `json:"value"`
	ValidFrom  time.Time `json:"validFrom"`
	ValidUntil time.Time `json:"validUntil"`
	Quantity   uint      `json:"quantity"`
	UsedCount  uint      `json:"usedCount"`
	MinTickets uint      `json:"minTickets"`
}

type MediaDetailsResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	Size      int64     `json:"mediaSize"`
	Type      string    `json:"mediaType"`
	MediaPath string    `json:"mediaPath"`
}

type ReserveTicketRequest struct {
	ID       uint
	Quantity uint
}

type ReserveTicketResponse struct {
	ID         uint    `json:"id"`
	FinalPrice float64 `json:"finalPrice"`
}

type EventAttendeesResponse struct {
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	Ticket       string  `json:"ticket"`
	CountTickets uint    `json:"countTickets"`
	Price        float64 `json:"price"`
}

type OrganizerDetailsResponse struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Profile     string `json:"profile"`
	Description string `json:"description"`
}
