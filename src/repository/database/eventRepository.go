package repository_database

import (
	"first-project/src/entities"
	"first-project/src/enums"
	"strings"
	"time"

	"gorm.io/gorm"
)

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *eventRepository {
	return &eventRepository{
		db: db,
	}
}

const queryByID = "id = ?"
const queryByEventID = "event_id = ?"
const queryByStatusIn = "status IN ?"

func (repo *eventRepository) FindDuplicatedEvent(name, venueType, location string, fromDate, toDate time.Time) (*entities.Event, bool) {
	var existingEvent entities.Event
	query := repo.db.Where("name = ? AND status != ?", name, enums.Cancelled)

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

func (repo *eventRepository) CreateNewEvent(event *entities.Event) *entities.Event {
	result := repo.db.Create(event)
	if result.Error != nil {
		panic(result.Error)
	}
	return event
}

func (repo *eventRepository) FindEventByID(eventID uint) (*entities.Event, bool) {
	var event entities.Event
	result := repo.db.First(&event, queryByID, eventID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &event, true
}

func (repo *eventRepository) FindEventCategoriesByEvent(event *entities.Event) []entities.Category {
	if err := repo.db.Model(event).Association("Categories").Find(&event.Categories); err != nil {
		panic(err)
	}
	return event.Categories
}

func (repo *eventRepository) FindTicketsByEventID(eventID uint, availability []bool) ([]*entities.Ticket, bool) {
	var tickets []*entities.Ticket
	result := repo.db.Where("event_id = ? AND is_available IN ?", eventID, availability).Find(&tickets)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return tickets, true
}

func (repo *eventRepository) FindDiscountsByEventID(eventID uint) ([]*entities.Discount, bool) {
	var discounts []*entities.Discount
	result := repo.db.Where(queryByEventID, eventID).Find(&discounts)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return discounts, true
}
func (repo *eventRepository) FindDiscountByDiscountID(discountID uint) (*entities.Discount, bool) {
	var discount entities.Discount
	result := repo.db.First(&discount, queryByID, discountID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &discount, true
}

func (repo *eventRepository) FindEventTicketByName(ticketName string, eventID uint) (*entities.Ticket, bool) {
	var ticket entities.Ticket
	result := repo.db.First(&ticket, "name = ? AND event_id = ?", ticketName, eventID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &ticket, true
}

func (repo *eventRepository) FindEventTicketByID(ticketID uint) (*entities.Ticket, bool) {
	var ticket entities.Ticket
	result := repo.db.First(&ticket, queryByID, ticketID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &ticket, true
}

func (repo *eventRepository) FindEventMediaByName(mediaName string, eventID uint) (*entities.Media, bool) {
	var media entities.Media
	result := repo.db.First(&media, "name = ? AND event_id = ?", mediaName, eventID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &media, true
}

func (repo *eventRepository) FindAllEventMedia(eventID uint) ([]*entities.Media, bool) {
	var media []*entities.Media
	result := repo.db.Where(queryByEventID, eventID).Find(&media)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return media, true
}

func (repo *eventRepository) FindAllEventOrganizers(eventID uint) ([]*entities.Organizer, bool) {
	var organizers []*entities.Organizer
	result := repo.db.Where(queryByEventID, eventID).Find(&organizers)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return organizers, true
}

func (repo *eventRepository) CreateNewTicket(ticket *entities.Ticket) *entities.Ticket {
	result := repo.db.Create(ticket)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticket
}

func (repo *eventRepository) FindEventDiscountByCode(discountCode string, eventID uint) (*entities.Discount, bool) {
	var discount entities.Discount
	result := repo.db.First(&discount, "code = ? AND event_id = ?", discountCode, eventID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &discount, true
}

func (repo *eventRepository) CreateNewDiscount(discount *entities.Discount) *entities.Discount {
	result := repo.db.Create(discount)
	if result.Error != nil {
		panic(result.Error)
	}
	return discount
}

func (repo *eventRepository) UpdateEventCategories(eventID uint, categories []entities.Category) {
	err := repo.db.Model(&entities.Event{ID: eventID}).Association("Categories").Replace(categories)
	if err != nil {
		panic(err)
	}
}

func (repo *eventRepository) UpdateEvent(event *entities.Event) {
	if err := repo.db.Save(event).Error; err != nil {
		panic(err)
	}
}

func (repo *eventRepository) FindOrganizerByID(organizerID uint) (*entities.Organizer, bool) {
	var organizer entities.Organizer
	result := repo.db.First(&organizer, queryByID, organizerID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &organizer, true
}

func (repo *eventRepository) FindOrganizerByEmail(eventID uint, email string) (*entities.Organizer, bool) {
	var organizer entities.Organizer
	result := repo.db.Where("email = ? AND event_id = ?", email, eventID).First(&organizer)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &organizer, true
}

func (repo *eventRepository) CreateOrganizerForEventID(eventID uint, name, email, description string) *entities.Organizer {
	organizer := &entities.Organizer{
		Name:        name,
		Email:       email,
		Description: description,
		EventID:     eventID,
	}
	result := repo.db.Create(organizer)
	if result.Error != nil {
		panic(result.Error)
	}
	return organizer
}

func (repo *eventRepository) FindEventsByStatus(allowedStatus []enums.EventStatus, offset, pageSize int) ([]*entities.Event, bool) {
	var events []*entities.Event
	result := repo.db.Where(queryByStatusIn, allowedStatus).Offset(offset).Limit(pageSize).Find(&events)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return events, true
}

func (repo *eventRepository) DeleteEvent(eventID uint) {
	err := repo.db.Unscoped().Delete(&entities.Event{}, eventID).Error
	if err != nil {
		panic(err)
	}
}

func (repo *eventRepository) DeleteTicket(ticketID uint) {
	err := repo.db.Unscoped().Delete(&entities.Ticket{}, ticketID).Error
	if err != nil {
		panic(err)
	}
}

func (repo *eventRepository) DeleteDiscount(discountID uint) {
	err := repo.db.Unscoped().Delete(&entities.Discount{}, discountID).Error
	if err != nil {
		panic(err)
	}
}

func (repo *eventRepository) DeleteOrganizer(organizerID uint) {
	err := repo.db.Unscoped().Delete(&entities.Organizer{}, organizerID).Error
	if err != nil {
		panic(err)
	}
}

func (repo *eventRepository) ChangeStatusByEvent(event *entities.Event, newStatus enums.EventStatus) {
	event.Status = newStatus
	repo.db.Save(event)
}

func (repo *eventRepository) CreateNewMedia(media *entities.Media) *entities.Media {
	result := repo.db.Create(media)
	if result.Error != nil {
		panic(result.Error)
	}
	return media
}

func (repo *eventRepository) FindMediaByID(mediaID uint) (*entities.Media, bool) {
	var media entities.Media
	result := repo.db.Where(queryByID, mediaID).First(&media)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &media, true
}

func (repo *eventRepository) UpdateEventMedia(media *entities.Media) {
	if err := repo.db.Save(media).Error; err != nil {
		panic(err)
	}
}

func (repo *eventRepository) DeleteMedia(mediaID uint) {
	result := repo.db.Unscoped().Delete(&entities.Media{}, mediaID)
	if result.Error != nil {
		panic(result.Error)
	}
}

func (repo *eventRepository) UpdateEventTicket(ticket *entities.Ticket) *entities.Ticket {
	result := repo.db.Save(ticket)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticket
}

func (repo *eventRepository) UpdateEventDiscount(discount *entities.Discount) {
	result := repo.db.Save(discount)
	if result.Error != nil {
		panic(result.Error)
	}
}

func (repo *eventRepository) UpdateEventOrganizer(organizer *entities.Organizer) {
	result := repo.db.Save(organizer)
	if result.Error != nil {
		panic(result.Error)
	}
}

func (repo *eventRepository) FullTextSearch(query string, allowedStatus []enums.EventStatus, offset, pageSize int) []*entities.Event {
	var events []*entities.Event

	repo.db.Exec(`ALTER TABLE events ADD FULLTEXT INDEX idx_name_description (name, description)`)
	searchQuery := "+" + strings.Join(strings.Fields(query), "* +") + "*"

	result := repo.db.Model(&entities.Event{}).
		Where("MATCH(name, description) AGAINST(? IN BOOLEAN MODE)", searchQuery).
		Where(queryByStatusIn, allowedStatus).
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

func (repo *eventRepository) FindEventsByCategoryName(categories []string, offset, pageSize int, allowedStatus []enums.EventStatus) []*entities.Event {
	var events []*entities.Event

	result := repo.db.
		Distinct("events.*").
		Joins("JOIN event_categories ON events.id = event_categories.event_id").
		Joins("JOIN categories ON categories.id = event_categories.category_id").
		Where("categories.name IN ?", categories).
		Where(queryByStatusIn, allowedStatus).
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
