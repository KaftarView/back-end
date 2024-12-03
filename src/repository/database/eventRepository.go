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

const queryByIDAndEventID = "id = ? AND event_id = ?"

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

func (repo *EventRepository) FindOrganizerByEventIDAndEmailAndVerified(eventID uint, email string, verified bool) (entities.Organizer, bool) {
	var organizer entities.Organizer
	result := repo.db.First(&organizer, "event_id = ? AND email = ? AND verified = ?", eventID, email, verified)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return organizer, false
		}
		panic(result.Error)
	}
	return organizer, true
}

func (repo *EventRepository) UpdateOrganizerToken(organizer entities.Organizer, token string) {
	organizer.Token = token
	repo.db.Save(&organizer)
}

func (repo *EventRepository) FindOrganizerByIDAndEventIDAndVerified(organizerID, eventID uint, verified bool) (entities.Organizer, bool) {
	var organizer entities.Organizer
	result := repo.db.First(&organizer, "id = ? AND event_id = ? AND verified = ?", organizerID, eventID, verified)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return organizer, false
		}
		panic(result.Error)
	}
	return organizer, true
}

func (repo *EventRepository) FindEventCategoriesByEvent(event entities.Event) entities.Event {
	if err := repo.db.Model(&event).Association("Categories").Find(&event.Categories); err != nil {
		panic(err)
	}
	return event
}

func (repo *EventRepository) FindTicketsByEventID(eventID uint) ([]entities.Ticket, bool) {
	var tickets []entities.Ticket
	result := repo.db.Where("event_id = ?", eventID).Find(&tickets)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return tickets, false
		}
		panic(result.Error)
	}
	return tickets, true
}

func (repo *EventRepository) FindDiscountsByEventID(eventID uint) ([]entities.Discount, bool) {
	var discounts []entities.Discount
	result := repo.db.Where("event_id = ?", eventID).Find(&discounts)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return discounts, false
		}
		panic(result.Error)
	}
	return discounts, true
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

func (repo *EventRepository) FindEventMediaByName(mediaName string, eventID uint) (entities.Media, bool) {
	var media entities.Media
	result := repo.db.First(&media, "name = ? AND event_id = ?", mediaName, eventID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return media, false
		}
		panic(result.Error)
	}
	return media, true
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

func (repo *EventRepository) FindActiveOrVerifiedOrganizerByEmail(eventID uint, email string) (entities.Organizer, bool) {
	var organizer entities.Organizer
	result := repo.db.Where("email = ? AND event_id = ?", email, eventID).First(&organizer)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return organizer, false
		}
		panic(result.Error)
	}
	if organizer.Verified || time.Since(organizer.UpdatedAt) < 8*time.Hour {
		return organizer, true
	}
	repo.db.Delete(&organizer)
	return organizer, false
}

func (repo *EventRepository) CreateOrganizerForEventID(eventID uint, name, email, description, token string, verified bool) entities.Organizer {
	organizer := entities.Organizer{
		Name:        name,
		Email:       email,
		Description: description,
		Token:       token,
		Verified:    verified,
		EventID:     eventID,
	}
	result := repo.db.Create(&organizer)
	if result.Error != nil {
		panic(result.Error)
	}
	return organizer
}

func (repo *EventRepository) ActivateOrganizer(organizer entities.Organizer) {
	organizer.Verified = true
	organizer.Token = ""
	if err := repo.db.Save(&organizer).Error; err != nil {
		panic(err)
	}
}

func (repo *EventRepository) FindEventsByStatus(allowedStatus []enums.EventStatus) ([]entities.Event, bool) {
	var events []entities.Event
	result := repo.db.Where("status IN ?", allowedStatus).Find(&events)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return events, false
		}
		panic(result.Error)
	}
	return events, true
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

func (repo *EventRepository) DeleteTicket(eventID, ticketID uint) bool {
	var ticket entities.Ticket
	result := repo.db.Where(queryByIDAndEventID, ticketID, eventID).First(&ticket)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false
		}
		panic(result.Error)
	}
	if err := repo.db.Delete(&ticket).Error; err != nil {
		panic(fmt.Errorf("failed to delete ticket: %w", err))
	}
	return true
}

func (repo *EventRepository) DeleteDiscount(eventID, discountID uint) bool {
	var discount entities.Discount
	result := repo.db.Where(queryByIDAndEventID, discountID, eventID).First(&discount)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false
		}
		panic(result.Error)
	}
	if err := repo.db.Delete(&discount).Error; err != nil {
		panic(fmt.Errorf("failed to delete discount: %w", err))
	}
	return true
}

func (repo *EventRepository) UpdateEventBannerByEventID(mediaPath string, eventID uint) {
	var event entities.Event
	if err := repo.db.Model(&event).Where("id = ?", eventID).Update("banner_path", mediaPath).Error; err != nil {
		panic(err)
	}
}

func (repo *EventRepository) ChangeStatusByEvent(event entities.Event, newStatus enums.EventStatus) {
	event.Status = newStatus
	repo.db.Save(&event)
}

func (repo *EventRepository) CreateNewMedia(media entities.Media) entities.Media {
	result := repo.db.Create(&media)
	if result.Error != nil {
		panic(result.Error)
	}
	return media
}

func (repo *EventRepository) FindMediaByIDAndEventID(mediaID, eventID uint) (entities.Media, bool) {
	var media entities.Media
	result := repo.db.Where(queryByIDAndEventID, mediaID, eventID).First(&media)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return media, false
		}
		panic(result.Error)
	}
	return media, true
}

func (repo *EventRepository) DeleteMedia(mediaID uint) bool {
	result := repo.db.Delete(&entities.Media{}, mediaID)
	if result.Error != nil {
		panic(result.Error)
	}
	if result.RowsAffected == 0 {
		return false
	}
	return true
}
