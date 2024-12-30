package repository_database_interfaces

import (
	"first-project/src/entities"
	"first-project/src/enums"
	"time"
)

type EventRepository interface {
	ChangeStatusByEvent(event *entities.Event, newStatus enums.EventStatus)
	CreateNewDiscount(discount *entities.Discount) *entities.Discount
	CreateNewEvent(event *entities.Event) *entities.Event
	CreateNewMedia(media *entities.Media) *entities.Media
	CreateNewTicket(ticket *entities.Ticket) *entities.Ticket
	CreateOrganizerForEventID(eventID uint, name string, email string, description string) *entities.Organizer
	DeleteDiscount(discountID uint)
	DeleteEvent(eventID uint)
	DeleteMedia(mediaID uint)
	DeleteOrganizer(organizerID uint)
	DeleteTicket(ticketID uint)
	FindAllEventMedia(eventID uint) ([]*entities.Media, bool)
	FindAllEventOrganizers(eventID uint) ([]*entities.Organizer, bool)
	FindDiscountByDiscountID(discountID uint) (*entities.Discount, bool)
	FindDiscountsByEventID(eventID uint) ([]*entities.Discount, bool)
	FindDuplicatedEvent(name string, venueType string, location string, fromDate time.Time, toDate time.Time) (*entities.Event, bool)
	FindEventByID(eventID uint) (*entities.Event, bool)
	FindEventCategoriesByEvent(event *entities.Event) []entities.Category
	FindEventDiscountByCode(discountCode string, eventID uint) (*entities.Discount, bool)
	FindEventMediaByName(mediaName string, eventID uint) (*entities.Media, bool)
	FindEventTicketByID(ticketID uint) (*entities.Ticket, bool)
	FindEventTicketByName(ticketName string, eventID uint) (*entities.Ticket, bool)
	FindEventsByCategoryName(categories []string, offset int, pageSize int, allowedStatus []enums.EventStatus) []*entities.Event
	FindEventsByStatus(allowedStatus []enums.EventStatus, offset int, pageSize int) ([]*entities.Event, bool)
	FindMediaByID(mediaID uint) (*entities.Media, bool)
	FindOrganizerByEmail(eventID uint, email string) (*entities.Organizer, bool)
	FindOrganizerByID(organizerID uint) (*entities.Organizer, bool)
	FindTicketsByEventID(eventID uint, availability []bool) ([]*entities.Ticket, bool)
	FullTextSearch(query string, allowedStatus []enums.EventStatus, offset int, pageSize int) []*entities.Event
	UpdateEvent(event *entities.Event)
	UpdateEventCategories(eventID uint, categories []entities.Category)
	UpdateEventDiscount(discount *entities.Discount)
	UpdateEventMedia(media *entities.Media)
	UpdateEventTicket(ticket *entities.Ticket) *entities.Ticket
	UpdateEventOrganizer(organizer *entities.Organizer)
}
