package repository_database_interfaces

import (
	"first-project/src/entities"
	"first-project/src/enums"
	"time"

	"gorm.io/gorm"
)

type EventRepository interface {
	ChangeStatusByEvent(db *gorm.DB, event *entities.Event, newStatus enums.EventStatus)
	CreateNewDiscount(db *gorm.DB, discount *entities.Discount) error
	CreateNewEvent(db *gorm.DB, event *entities.Event) error
	CreateNewMedia(db *gorm.DB, media *entities.Media) error
	CreateNewTicket(db *gorm.DB, ticket *entities.Ticket) error
	CreateOrganizerForEventID(db *gorm.DB, organizer *entities.Organizer) error
	CreateReservation(db *gorm.DB, reservation *entities.Reservation) error
	DeleteDiscount(db *gorm.DB, discountID uint) error
	DeleteEvent(db *gorm.DB, eventID uint) error
	DeleteMedia(db *gorm.DB, mediaID uint) error
	DeleteOrganizer(db *gorm.DB, organizerID uint) error
	DeleteTicket(db *gorm.DB, ticketID uint) error
	FindAllEventMedia(db *gorm.DB, eventID uint) ([]*entities.Media, bool)
	FindAllEventOrganizers(db *gorm.DB, eventID uint) ([]*entities.Organizer, bool)
	FindDiscountByDiscountID(db *gorm.DB, discountID uint) (*entities.Discount, bool)
	FindDiscountsByEventID(db *gorm.DB, eventID uint) ([]*entities.Discount, bool)
	FindDuplicatedEvent(db *gorm.DB, name string, venueType string, location string, fromDate time.Time, toDate time.Time) (*entities.Event, bool)
	FindEventByID(db *gorm.DB, eventID uint) (*entities.Event, bool)
	FindEventCategoriesByEvent(db *gorm.DB, event *entities.Event) []entities.Category
	FindEventDiscountByCode(db *gorm.DB, discountCode string, eventID uint) (*entities.Discount, bool)
	FindEventDiscountByCodeForUpdate(db *gorm.DB, discountCode string, eventID uint) (*entities.Discount, bool)
	FindEventMediaByName(db *gorm.DB, mediaName string, eventID uint) (*entities.Media, bool)
	FindEventTicketByID(db *gorm.DB, ticketID uint) (*entities.Ticket, bool)
	FindEventTicketByIDForUpdate(db *gorm.DB, ticketID uint) (*entities.Ticket, bool)
	FindEventTicketByName(db *gorm.DB, ticketName string, eventID uint) (*entities.Ticket, bool)
	FindEventsByCategoryName(db *gorm.DB, categories []string, offset int, pageSize int, allowedStatus []enums.EventStatus) []*entities.Event
	FindEventsByStatus(db *gorm.DB, allowedStatus []enums.EventStatus, offset int, pageSize int) ([]*entities.Event, bool)
	FindMediaByID(db *gorm.DB, mediaID uint) (*entities.Media, bool)
	FindOrganizerByEmail(db *gorm.DB, eventID uint, email string) (*entities.Organizer, bool)
	FindOrganizerByID(db *gorm.DB, organizerID uint) (*entities.Organizer, bool)
	FindTicketsByEventID(db *gorm.DB, eventID uint, availability []bool) ([]*entities.Ticket, bool)
	FullTextSearch(db *gorm.DB, query string, allowedStatus []enums.EventStatus, offset int, pageSize int) []*entities.Event
	UpdateEvent(db *gorm.DB, event *entities.Event) error
	UpdateEventCategories(db *gorm.DB, eventID uint, categories []entities.Category) error
	UpdateEventDiscount(db *gorm.DB, discount *entities.Discount) error
	UpdateEventMedia(db *gorm.DB, media *entities.Media) error
	UpdateEventTicket(db *gorm.DB, ticket *entities.Ticket) error
	UpdateEventOrganizer(db *gorm.DB, organizer *entities.Organizer) error
}
