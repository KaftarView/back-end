package repository_database

import (
	"first-project/src/entities"
	"first-project/src/enums"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EventRepository struct{}

func NewEventRepository() *EventRepository {
	return &EventRepository{}
}

const queryByID = "id = ?"
const queryByEventID = "event_id = ?"
const queryByStatusIn = "status IN ?"

func (repo *EventRepository) FindDuplicatedEvent(db *gorm.DB, name, venueType, location string, fromDate, toDate time.Time) (*entities.Event, bool) {
	var existingEvent entities.Event
	query := db.Where("name = ? AND status != ?", name, enums.Cancelled)

	timeOverlapCondition := "(" +
		"(from_date BETWEEN ? AND ?) OR " +
		"(to_date BETWEEN ? AND ?) OR " +
		"(? BETWEEN from_date AND to_date) OR " +
		"(? BETWEEN from_date AND to_date)" +
		")"
	query = query.Where(
		timeOverlapCondition,
		fromDate, toDate,
		fromDate, toDate,
		fromDate,
		toDate,
	)

	result := query.First(&existingEvent)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &existingEvent, true
}

func (repo *EventRepository) CreateNewEvent(db *gorm.DB, event *entities.Event) error {
	return db.Create(event).Error
}

func (repo *EventRepository) FindEventByID(db *gorm.DB, eventID uint) (*entities.Event, bool) {
	var event entities.Event
	result := db.First(&event, queryByID, eventID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &event, true
}

func (repo *EventRepository) FindEventCategoriesByEvent(db *gorm.DB, event *entities.Event) []entities.Category {
	if err := db.Model(event).Association("Categories").Find(&event.Categories); err != nil {
		panic(err)
	}
	return event.Categories
}

func (repo *EventRepository) FindAvailableTicketsByEventID(db *gorm.DB, eventID uint) ([]*entities.Ticket, bool) {
	var tickets []*entities.Ticket
	now := time.Now()

	result := db.Where(&entities.Ticket{
		EventID:     eventID,
		IsAvailable: true,
	}).Where("available_from <= ?", now).
		Where("available_until >= ?", now).
		Order("created_at DESC").
		Find(&tickets)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return tickets, true
}

func (repo *EventRepository) FindAllTicketsByEventID(db *gorm.DB, eventID uint) ([]*entities.Ticket, bool) {
	var tickets []*entities.Ticket
	result := db.Where(queryByEventID, eventID).Order("created_at DESC").Find(&tickets)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return tickets, true
}

func (repo *EventRepository) FindDiscountsByEventID(db *gorm.DB, eventID uint) ([]*entities.Discount, bool) {
	var discounts []*entities.Discount
	result := db.Where(queryByEventID, eventID).Order("created_at DESC").Order("created_at DESC").Find(&discounts)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return discounts, true
}
func (repo *EventRepository) FindDiscountByDiscountID(db *gorm.DB, discountID uint) (*entities.Discount, bool) {
	var discount entities.Discount
	result := db.First(&discount, queryByID, discountID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &discount, true
}

func (repo *EventRepository) FindEventTicketByName(db *gorm.DB, ticketName string, eventID uint) (*entities.Ticket, bool) {
	var ticket entities.Ticket
	result := db.First(&ticket, "name = ? AND event_id = ?", ticketName, eventID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &ticket, true
}

func (repo *EventRepository) FindEventTicketByID(db *gorm.DB, ticketID uint) (*entities.Ticket, bool) {
	var ticket entities.Ticket
	result := db.First(&ticket, queryByID, ticketID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &ticket, true
}

func (repo *EventRepository) FindEventTicketByIDForUpdate(db *gorm.DB, ticketID uint) (*entities.Ticket, bool) {
	var ticket entities.Ticket
	result := db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&ticket, queryByID, ticketID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &ticket, true
}

func (repo *EventRepository) FindEventMediaByName(db *gorm.DB, mediaName string, eventID uint) (*entities.Media, bool) {
	var media entities.Media
	result := db.First(&media, "name = ? AND event_id = ?", mediaName, eventID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &media, true
}

func (repo *EventRepository) FindAllEventMedia(db *gorm.DB, eventID uint) ([]*entities.Media, bool) {
	var media []*entities.Media
	result := db.Where(queryByEventID, eventID).Order("created_at DESC").Find(&media)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return media, true
}

func (repo *EventRepository) FindAllEventOrganizers(db *gorm.DB, eventID uint) ([]*entities.Organizer, bool) {
	var organizers []*entities.Organizer
	result := db.Where(queryByEventID, eventID).Order("created_at DESC").Find(&organizers)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return organizers, true
}

func (repo *EventRepository) CreateNewTicket(db *gorm.DB, ticket *entities.Ticket) error {
	return db.Create(ticket).Error
}

func (repo *EventRepository) FindEventDiscountByCode(db *gorm.DB, discountCode string, eventID uint) (*entities.Discount, bool) {
	var discount entities.Discount
	result := db.First(&discount, "code = ? AND event_id = ?", discountCode, eventID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &discount, true
}

func (repo *EventRepository) FindEventDiscountByCodeForUpdate(db *gorm.DB, discountCode string, eventID uint) (*entities.Discount, bool) {
	var discount entities.Discount
	result := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&discount, "code = ? AND event_id = ?", discountCode, eventID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &discount, true
}

func (repo *EventRepository) CreateNewDiscount(db *gorm.DB, discount *entities.Discount) error {
	return db.Create(discount).Error
}

func (repo *EventRepository) UpdateEventCategories(db *gorm.DB, eventID uint, categories []entities.Category) error {
	return db.Model(&entities.Event{ID: eventID}).Association("Categories").Replace(categories)
}

func (repo *EventRepository) UpdateEvent(db *gorm.DB, event *entities.Event) error {
	return db.Save(event).Error
}

func (repo *EventRepository) FindOrganizerByID(db *gorm.DB, organizerID uint) (*entities.Organizer, bool) {
	var organizer entities.Organizer
	result := db.First(&organizer, queryByID, organizerID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &organizer, true
}

func (repo *EventRepository) FindOrganizerByEmail(db *gorm.DB, eventID uint, email string) (*entities.Organizer, bool) {
	var organizer entities.Organizer
	result := db.Where("email = ? AND event_id = ?", email, eventID).First(&organizer)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &organizer, true
}

func (repo *EventRepository) CreateOrganizer(db *gorm.DB, organizer *entities.Organizer) error {
	return db.Create(organizer).Error
}

func (repo *EventRepository) FindEventsByStatus(db *gorm.DB, allowedStatus []enums.EventStatus, offset, pageSize int) ([]*entities.Event, bool) {
	var events []*entities.Event
	result := db.Where(queryByStatusIn, allowedStatus).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&events)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return events, true
}

func (repo *EventRepository) DeleteEvent(db *gorm.DB, eventID uint) error {
	return db.Delete(&entities.Event{}, eventID).Error
}

func (repo *EventRepository) DeleteTicket(db *gorm.DB, ticketID uint) error {
	return db.Delete(&entities.Ticket{}, ticketID).Error
}

func (repo *EventRepository) DeleteDiscount(db *gorm.DB, discountID uint) error {
	return db.Unscoped().Delete(&entities.Discount{}, discountID).Error
}

func (repo *EventRepository) DeleteOrganizer(db *gorm.DB, organizerID uint) error {
	return db.Unscoped().Delete(&entities.Organizer{}, organizerID).Error
}

func (repo *EventRepository) ChangeStatusByEvent(db *gorm.DB, event *entities.Event, newStatus enums.EventStatus) {
	event.Status = newStatus
	db.Save(event)
}

func (repo *EventRepository) CreateNewMedia(db *gorm.DB, media *entities.Media) error {
	return db.Create(media).Error
}

func (repo *EventRepository) FindMediaByID(db *gorm.DB, mediaID uint) (*entities.Media, bool) {
	var media entities.Media
	result := db.Where(queryByID, mediaID).First(&media)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &media, true
}

func (repo *EventRepository) FindOrdersEventID(db *gorm.DB, eventID uint) []*entities.Order {
	var orders []*entities.Order

	result := db.Preload("User").
		Preload("Items").
		Where(queryByEventID, eventID).
		Find(&orders)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return []*entities.Order{}
		}
		panic(result.Error)
	}
	return orders
}

func (repo *EventRepository) IsUserAttendingEvent(db *gorm.DB, userID, eventID uint) bool {
	var count int64

	err := db.Model(&entities.Order{}).
		Where("user_id = ? AND event_id = ?", userID, eventID).
		Count(&count).Error

	if err != nil {
		panic(err)
	}
	return count > 0
}

func (repo *EventRepository) UpdateEventMedia(db *gorm.DB, media *entities.Media) error {
	return db.Save(media).Error
}

func (repo *EventRepository) DeleteMedia(db *gorm.DB, mediaID uint) error {
	return db.Unscoped().Delete(&entities.Media{}, mediaID).Error
}

func (repo *EventRepository) UpdateEventTicket(db *gorm.DB, ticket *entities.Ticket) error {
	return db.Save(ticket).Error
}

func (repo *EventRepository) UpdateEventDiscount(db *gorm.DB, discount *entities.Discount) error {
	return db.Save(discount).Error
}

func (repo *EventRepository) UpdateEventOrganizer(db *gorm.DB, organizer *entities.Organizer) error {
	return db.Save(organizer).Error
}

func (repo *EventRepository) FullTextSearch(db *gorm.DB, query string, allowedStatus []enums.EventStatus, offset, pageSize int) []*entities.Event {
	var events []*entities.Event

	db.Exec(`ALTER TABLE events ADD FULLTEXT INDEX idx_name_description (name, description)`)
	searchQuery := "+" + strings.Join(strings.Fields(query), "* +") + "*"

	result := db.Model(&entities.Event{}).
		Where("MATCH(name, description) AGAINST(? IN BOOLEAN MODE)", searchQuery).
		Where(queryByStatusIn, allowedStatus).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&events)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return events
}

func (repo *EventRepository) FindEventsByCategoryName(db *gorm.DB, categories []string, offset, pageSize int, allowedStatus []enums.EventStatus) []*entities.Event {
	var events []*entities.Event

	result := db.
		Distinct("events.*").
		Joins("JOIN event_categories ON events.id = event_categories.event_id").
		Joins("JOIN categories ON categories.id = event_categories.category_id").
		Where("categories.name IN ?", categories).
		Where(queryByStatusIn, allowedStatus).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&events)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}

	return events
}
