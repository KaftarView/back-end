package repository_database

import (
	"first-project/src/entities"
	"first-project/src/enums"
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

func (repo *EventRepository) CheckVenueAvailability(venueType, location string, fromDate, toDate time.Time) (entities.Event, bool) {
	var event entities.Event
	if venueType == enums.Online.String() {
		return event, true
	}

	result := repo.db.First(&event).
		Where("location = ? AND status != ?", location, enums.Cancelled).
		Where("(from_date, to_date) OVERLAPS (?, ?)", fromDate, toDate)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return event, true
		}
		panic(result.Error)
	}
	return event, false
}

func (repo *EventRepository) CreateNewEvent(event entities.Event) entities.Event {
	result := repo.db.Create(&event)
	if result.Error != nil {
		panic(result.Error)
	}
	return event
}
