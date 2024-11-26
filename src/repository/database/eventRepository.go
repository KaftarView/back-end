package repository_database

import (
	"first-project/src/entities"
	"first-project/src/enums"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (repo *EventRepository) FindDuplicatedEvent(name, venueType, location string, fromDate, toDate time.Time) (entities.Event, bool) {
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
			return existingEvent, false
		}
		panic(result.Error)
	}
	return existingEvent, true
}

func (repo *EventRepository) CreateNewEvent(event entities.Event) entities.Event {
	result := repo.db.Create(&event)
	if result.Error != nil {
		panic(result.Error)
	}
	return event
}

func (repo *EventRepository) FindCategoriesByNames(categoryNames []string) []entities.Category {
	var categories []entities.Category

	for _, categoryName := range categoryNames {
		var category entities.Category
		if err := repo.db.FirstOrCreate(&category, entities.Category{Name: categoryName}).Error; err != nil {
			panic(err)
		}
		categories = append(categories, category)
	}
	return categories
}

func (repo *EventRepository) FindEventByID(eventID uint) (entities.Event, bool) {
	var event entities.Event
	result := repo.db.First(&event, "id = ?", eventID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return event, false
		}
		panic(result.Error)
	}
	return event, true
}

func (repo *EventRepository) FindEventTicketByName(ticketName string, eventID uint) (entities.Ticket, bool) {
	var ticket entities.Ticket
	result := repo.db.First(&ticket, "name = ? AND event_id = ?", ticketName, eventID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return ticket, false
		}
		panic(result.Error)
	}
	return ticket, true
}

func (repo *EventRepository) CreateNewTicket(ticket entities.Ticket) entities.Ticket {
	result := repo.db.Create(&ticket)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticket
}

func (repo *EventRepository) FindEventDiscountByCode(discountCode string, eventID uint) (entities.Discount, bool) {
	var discount entities.Discount
	result := repo.db.First(&discount, "code = ? AND event_id = ?", discountCode, eventID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return discount, false
		}
		panic(result.Error)
	}
	return discount, true
}

func (repo *EventRepository) CreateNewDiscount(discount entities.Discount) entities.Discount {
	result := repo.db.Create(&discount)
	if result.Error != nil {
		panic(result.Error)
	}
	return discount
}

func (repo *EventRepository) FindEventsByStatus(allowedStatus []enums.EventStatus) []entities.Event {
	var events []entities.Event
	result := repo.db.Where("status IN ?", allowedStatus).Find(&events)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return events
		}
		panic(result.Error)
	}
	return events
}

func (repo *EventRepository) FetchEventDetailsAfterFetching(event entities.Event) entities.Event {
	err := repo.db.Model(&event).
		Preload("Tickets.Purchasable").
		Preload("Organizers").
		Preload("Categories").
		Preload("Commentable.Comments.User").
		Find(&event).Error

	if err != nil {
		panic(fmt.Errorf("failed to preload event details: %w", err))
	}
	return event
}

func (repo *EventRepository) FindAllCategories() []string {
	var categoryNames []string
	result := repo.db.Model(&entities.Category{}).Pluck("name", &categoryNames)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return []string{}
		}
		panic(result.Error)
	}
	return categoryNames
}

func (repo *EventRepository) DeleteEvent(eventID uint) bool {
	result := repo.db.Delete(&entities.Event{}, eventID)
	if result.Error != nil {
		panic(result.Error)
	}
	if result.RowsAffected == 0 {
		return false
	}
	return true
}

func (repo *EventRepository) UpdateEventBannerByEventID(mediaPath string, eventID uint) {
	var event entities.Event
	if err := repo.db.Model(&event).Where("id = ?", eventID).Update("banner_path", mediaPath).Error; err != nil {
		panic(err)
	}
}
