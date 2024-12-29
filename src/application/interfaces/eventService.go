package application_interfaces

import (
	"first-project/src/dto"
	"first-project/src/entities"
	"first-project/src/enums"
	"mime/multipart"
	"time"
)

type EventService interface {
	ChangeEventStatus(eventID uint, newStatus string)
	CreateEvent(eventDetails dto.CreateEventRequest) *entities.Event
	CreateEventDiscount(discountDetails dto.CreateDiscountRequest) *entities.Discount
	CreateEventMedia(eventID uint, mediaName string, mediaFile *multipart.FileHeader)
	CreateEventOrganizer(eventID uint, name string, email string, description string, profile *multipart.FileHeader)
	CreateEventTicket(ticketDetails dto.CreateTicketRequest) *entities.Ticket
	DeleteDiscount(discountID uint)
	DeleteEvent(eventID uint)
	DeleteEventMedia(mediaID uint)
	DeleteOrganizer(organizerID uint)
	DeleteTicket(ticketID uint)
	FilterEventsByCategories(categories []string, page int, pageSize int, allowedStatus []enums.EventStatus) []dto.EventDetailsResponse
	GetDiscountDetails(discountID uint) dto.DiscountDetailsResponse
	GetEventByID(eventID uint) *entities.Event
	GetEventDetails(allowedStatus []enums.EventStatus, eventID uint) dto.EventDetailsResponse
	GetEventDiscounts(eventID uint) []dto.DiscountDetailsResponse
	GetEventMediaDetails(mediaID uint) dto.MediaDetailsResponse
	GetEventTickets(eventID uint, availability []bool) []dto.TicketDetailsResponse
	GetEventsList(allowedStatus []enums.EventStatus, page int, pageSize int) []dto.EventDetailsResponse
	GetListEventMedia(eventID uint) []dto.MediaDetailsResponse
	GetTicketDetails(ticketID uint) dto.TicketDetailsResponse
	SearchEvents(query string, page int, pageSize int, allowedStatus []enums.EventStatus) []dto.EventDetailsResponse
	UpdateEvent(updateDetails dto.UpdateEventRequest)
	UpdateEventDiscount(discountDetails dto.UpdateDiscountRequest)
	UpdateEventMedia(mediaID uint, name *string, file *multipart.FileHeader)
	UpdateEventTicket(ticketDetails dto.UpdateTicketRequest)
	ValidateEventCreationDetails(name string, venueType string, location string, fromDate time.Time, toDate time.Time)
	ValidateNewEventDiscountDetails(discountCode string, eventID uint)
	ValidateNewEventTicketDetails(ticketName string, eventID uint)
}
