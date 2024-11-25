package application

import (
	"first-project/src/bootstrap"
	"first-project/src/dto"
	"first-project/src/entities"
	"first-project/src/enums"
	"first-project/src/exceptions"
	repository_database "first-project/src/repository/database"
	"time"
)

type EventService struct {
	constants       *bootstrap.Constants
	eventRepository *repository_database.EventRepository
}

func NewEventService(constants *bootstrap.Constants, eventRepository *repository_database.EventRepository) *EventService {
	return &EventService{
		constants:       constants,
		eventRepository: eventRepository,
	}
}

func (eventService *EventService) ValidateEventCreationDetails(
	name, venueType, location string, fromDate, toDate time.Time,
) {
	var conflictError exceptions.ConflictError
	_, eventExist := eventService.eventRepository.FindDuplicatedEvent(name, venueType, location, fromDate, toDate)
	if eventExist {
		conflictError.AppendError(
			eventService.constants.ErrorField.Tittle,
			eventService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}
}

func (eventService *EventService) CreateEvent(eventDetails dto.CreateEventDetails) entities.Event {
	enumStatus := enums.Draft
	eventStatuses := enums.GetAllEventStatus()
	for _, eventStatus := range eventStatuses {
		if eventStatus.String() == eventDetails.Status {
			enumStatus = eventStatus
		}
	}

	enumVenue := enums.Online
	eventVenues := enums.GetAllEventVenues()
	for _, eventVenue := range eventVenues {
		if eventVenue.String() == eventDetails.VenueType {
			enumVenue = eventVenue
		}
	}

	categories := eventService.eventRepository.FindCategoriesByNames(eventDetails.Categories)

	eventDetailsModel := entities.Event{
		Name:        eventDetails.Name,
		Status:      enumStatus,
		Categories:  categories,
		Description: eventDetails.Description,
		FromDate:    eventDetails.FromDate,
		ToDate:      eventDetails.ToDate,
		MinCapacity: eventDetails.MinCapacity,
		MaxCapacity: eventDetails.MaxCapacity,
		VenueType:   enumVenue,
		Location:    eventDetails.Location,
	}
	event := eventService.eventRepository.CreateNewEvent(eventDetailsModel)
	return event
}

func (eventService *EventService) ValidateNewEventTicketDetails(ticketName string, eventID uint) entities.Ticket {
	var conflictError exceptions.ConflictError
	var notFoundError exceptions.NotFoundError
	_, eventExist := eventService.eventRepository.FindEventByID(eventID)
	if !eventExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Event
		panic(notFoundError)
	}
	ticket, ticketExist := eventService.eventRepository.FindEventTicketByName(ticketName, eventID)
	if ticketExist {
		conflictError.AppendError(
			eventService.constants.ErrorField.Ticket,
			eventService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}
	return ticket
}

func (eventService *EventService) CreateEventTicket(ticketDetails dto.CreateTicketDetails) entities.Ticket {
	ticketDetailsModel := entities.Ticket{
		Name:           ticketDetails.Name,
		Description:    ticketDetails.Description,
		Price:          ticketDetails.Price,
		Quantity:       ticketDetails.Quantity,
		SoldCount:      ticketDetails.SoldCount,
		IsAvailable:    ticketDetails.IsAvailable,
		AvailableFrom:  ticketDetails.AvailableFrom,
		AvailableUntil: ticketDetails.AvailableUntil,
		EventID:        ticketDetails.EventID,
	}
	ticket := eventService.eventRepository.CreateNewTicket(ticketDetailsModel)
	return ticket
}

func (eventService *EventService) ValidateNewEventDiscountDetails(discountCode string, eventID uint) entities.Discount {
	var conflictError exceptions.ConflictError
	var notFoundError exceptions.NotFoundError
	_, eventExist := eventService.eventRepository.FindEventByID(eventID)
	if !eventExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Event
		panic(notFoundError)
	}
	discount, discountExist := eventService.eventRepository.FindEventDiscountByCode(discountCode, eventID)
	if discountExist {
		conflictError.AppendError(
			eventService.constants.ErrorField.Discount,
			eventService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}
	return discount
}

func (eventService *EventService) CreateEventDiscount(discountDetails dto.CreateDiscountDetails) entities.Discount {
	var enumDiscountType enums.DiscountType
	discountTypes := enums.GetAllDiscountTypes()
	for _, discountType := range discountTypes {
		if discountType.String() == discountDetails.Type {
			enumDiscountType = discountType
		}
	}

	discountDetailsModel := entities.Discount{
		Code:       discountDetails.Code,
		Type:       enumDiscountType,
		Value:      discountDetails.Value,
		ValidFrom:  discountDetails.ValidFrom,
		ValidUntil: discountDetails.ValidUntil,
		Quantity:   discountDetails.Quantity,
		UsedCount:  discountDetails.UsedCount,
		MinTickets: discountDetails.MinTickets,
		EventID:    discountDetails.EventID,
	}
	discount := eventService.eventRepository.CreateNewDiscount(discountDetailsModel)
	return discount
}

func (eventService *EventService) GetListOfPublishedEvents() []entities.Event {
	allowedStatus := []enums.EventStatus{enums.Published}
	events := eventService.eventRepository.FindEventsByStatus(allowedStatus)
	return events
}

func (eventService *EventService) GetPublicEventDetails(eventID uint) entities.Event {
	var notFoundError exceptions.NotFoundError
	event, eventExist := eventService.eventRepository.FindEventByID(eventID)
	if !eventExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Event
		panic(notFoundError)
	}
	isValidStatus := false
	allowedStatus := []enums.EventStatus{enums.Published}
	for _, status := range allowedStatus {
		if event.Status == status {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		notFoundError.ErrorField = eventService.constants.ErrorField.Event
		panic(notFoundError)
	}
	event = eventService.eventRepository.FetchEventDetailsAfterFetching(event)
	return event
}

func (eventService *EventService) GetListOfCategories() []string {
	categoryNames := eventService.eventRepository.FindAllCategories()
	return categoryNames
}

func (eventService *EventService) DeleteEvent(eventID uint) {
	var notFoundError exceptions.NotFoundError
	eventExist := eventService.eventRepository.DeleteEvent(eventID)
	if !eventExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Event
		panic(notFoundError)
	}
}
