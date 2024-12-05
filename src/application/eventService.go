package application

import (
	"encoding/base64"
	"first-project/src/bootstrap"
	"first-project/src/dto"
	"first-project/src/entities"
	"first-project/src/enums"
	"first-project/src/exceptions"
	repository_database "first-project/src/repository/database"
	"time"
)

type EventService struct {
	constants         *bootstrap.Constants
	eventRepository   *repository_database.EventRepository
	commentRepository *repository_database.CommentRepository
}

func NewEventService(
	constants *bootstrap.Constants,
	eventRepository *repository_database.EventRepository,
	commentRepository *repository_database.CommentRepository,
) *EventService {
	return &EventService{
		constants:         constants,
		eventRepository:   eventRepository,
		commentRepository: commentRepository,
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
	commentable := eventService.commentRepository.CreateNewCommentable()

	eventDetailsModel := entities.Event{
		ID:          commentable.CID,
		Name:        eventDetails.Name,
		Status:      enumStatus,
		Categories:  categories,
		Description: eventDetails.Description,
		BasePrice:   eventDetails.BasePrice,
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

func (eventService *EventService) SetBannerPath(mediaPath string, eventID uint) {
	eventService.eventRepository.UpdateEventBannerByEventID(mediaPath, eventID)
}

func (eventService *EventService) ValidateNewEventTicketDetails(ticketName string, eventID uint) {
	var conflictError exceptions.ConflictError
	var notFoundError exceptions.NotFoundError
	_, eventExist := eventService.eventRepository.FindEventByID(eventID)
	if !eventExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Event
		panic(notFoundError)
	}
	_, mediaExist := eventService.eventRepository.FindEventTicketByName(ticketName, eventID)
	if mediaExist {
		conflictError.AppendError(
			eventService.constants.ErrorField.Media,
			eventService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}
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

func (eventService *EventService) ValidateNewEventDiscountDetails(discountCode string, eventID uint) {
	var conflictError exceptions.ConflictError
	var notFoundError exceptions.NotFoundError
	_, eventExist := eventService.eventRepository.FindEventByID(eventID)
	if !eventExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Event
		panic(notFoundError)
	}
	_, discountExist := eventService.eventRepository.FindEventDiscountByCode(discountCode, eventID)
	if discountExist {
		conflictError.AppendError(
			eventService.constants.ErrorField.Discount,
			eventService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}
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

func (eventService *EventService) UpdateOrCreateEventOrganizer(eventID uint, name, email, description, token string) uint {
	var notFoundError exceptions.NotFoundError
	var conflictError exceptions.ConflictError
	_, eventExist := eventService.eventRepository.FindEventByID(eventID)
	if !eventExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Event
		panic(notFoundError)
	}
	_, organizerExist := eventService.eventRepository.FindActiveOrVerifiedOrganizerByEmail(eventID, email)
	if organizerExist {
		conflictError.AppendError(
			eventService.constants.ErrorField.Organizer,
			eventService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}
	organizer, organizerExist := eventService.eventRepository.FindOrganizerByEventIDAndEmailAndVerified(eventID, email, false)
	if organizerExist {
		eventService.eventRepository.UpdateOrganizerToken(organizer, token)
		return organizer.ID
	}
	organizer = eventService.eventRepository.CreateOrganizerForEventID(eventID, name, email, description, token, false)
	return organizer.ID
}

func (eventService *EventService) GetEventName(eventID uint) string {
	event, _ := eventService.eventRepository.FindEventByID(eventID)
	return event.Name
}

func (eventService *EventService) ActivateUser(encodedOrganizerID, encodedEventID, token string) {
	decodedOrganizerID, err := base64.StdEncoding.DecodeString(encodedOrganizerID)
	if err != nil {
		panic(err)
	}
	decodedEventID, err := base64.StdEncoding.DecodeString(encodedEventID)
	if err != nil {
		panic(err)
	}
	organizerID := uint(decodedOrganizerID[0])
	eventID := uint(decodedEventID[0])
	var registrationError exceptions.UserRegistrationError
	var notFoundError exceptions.NotFoundError
	_, organizerExist := eventService.eventRepository.FindOrganizerByIDAndEventIDAndVerified(organizerID, eventID, true)
	if organizerExist {
		registrationError.AppendError(
			eventService.constants.ErrorField.Organizer,
			eventService.constants.ErrorTag.AlreadyVerified)
		panic(registrationError)
	}
	organizer, organizerExist := eventService.eventRepository.FindOrganizerByIDAndEventIDAndVerified(organizerID, eventID, false)
	if !organizerExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Organizer
		panic(notFoundError)
	}
	if organizer.Token != token {
		registrationError.AppendError(
			eventService.constants.ErrorField.Organizer,
			eventService.constants.ErrorTag.InvalidToken)
		panic(registrationError)
	}
	if time.Since(organizer.UpdatedAt) > 8*time.Hour {
		registrationError.AppendError(
			eventService.constants.ErrorField.Token,
			eventService.constants.ErrorTag.ExpiredToken)
		panic(registrationError)
	}
	eventService.eventRepository.ActivateOrganizer(organizer)
}

func (eventService *EventService) GetEventsList(allowedStatus []enums.EventStatus) []dto.EventDetailsResponse {
	events, _ := eventService.eventRepository.FindEventsByStatus(allowedStatus)
	eventsDetails := make([]dto.EventDetailsResponse, len(events))
	for i, event := range events {
		eventsDetails[i] = dto.EventDetailsResponse{
			ID:          event.ID,
			CreatedAt:   event.CreatedAt,
			Name:        event.Name,
			Status:      event.Status.String(),
			Description: event.Description,
			FromDate:    event.FromDate,
			ToDate:      event.ToDate,
			VenueType:   event.VenueType.String(),
			Banner:      event.BannerPath,
		}
	}
	return eventsDetails
}

func (eventService *EventService) GetEventDetails(allowedStatus []enums.EventStatus, eventID uint) dto.EventDetailsResponse {
	var notFoundError exceptions.NotFoundError
	event, eventExist := eventService.eventRepository.FindEventByID(eventID)
	if !eventExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Event
		panic(notFoundError)
	}
	isAllowStatus := false
	for _, status := range allowedStatus {
		if event.Status == status {
			isAllowStatus = true
		}
	}
	if !isAllowStatus {
		notFoundError.ErrorField = eventService.constants.ErrorField.Event
		panic(notFoundError)
	}
	event = eventService.eventRepository.FindEventCategoriesByEvent(event)
	categoryNames := make([]string, len(event.Categories))
	for i, category := range event.Categories {
		categoryNames[i] = category.Name
	}

	comments := eventService.commentRepository.GetCommentsByEventID(eventID)
	var commentDetails []dto.CommentDetails
	for _, comment := range comments {
		commentDetails = append(commentDetails, dto.CommentDetails{
			Content:     comment.Content,
			IsModerated: comment.IsModerated,
			AuthorName:  comment.Author.Name,
		})
	}

	eventDetails := dto.EventDetailsResponse{
		ID:          event.ID,
		CreatedAt:   event.CreatedAt,
		Name:        event.Name,
		Status:      event.Status.String(),
		Description: event.Description,
		BasePrice:   event.BasePrice,
		MinCapacity: event.MinCapacity,
		MaxCapacity: event.MaxCapacity,
		FromDate:    event.FromDate,
		ToDate:      event.ToDate,
		VenueType:   event.VenueType.String(),
		Location:    event.Location,
		Categories:  categoryNames,
		Banner:      event.BannerPath,
		Comments:    commentDetails,
	}
	return eventDetails
}

func (eventService *EventService) GetEventTickets(eventID uint) []dto.TicketDetailsResponse {
	var notFoundError exceptions.NotFoundError
	_, eventExist := eventService.eventRepository.FindEventByID(eventID)
	if !eventExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Event
		panic(notFoundError)
	}

	tickets, ticketExist := eventService.eventRepository.FindTicketsByEventID(eventID)
	if !ticketExist {
		return []dto.TicketDetailsResponse{}
	}
	ticketsDetails := make([]dto.TicketDetailsResponse, len(tickets))
	for i, ticket := range tickets {
		ticketsDetails[i] = dto.TicketDetailsResponse{
			ID:             ticket.ID,
			CreatedAt:      ticket.CreatedAt,
			Name:           ticket.Name,
			Description:    ticket.Description,
			Price:          ticket.Price,
			Quantity:       ticket.Quantity,
			IsAvailable:    ticket.IsAvailable,
			AvailableFrom:  ticket.AvailableFrom,
			AvailableUntil: ticket.AvailableUntil,
		}
	}
	return ticketsDetails
}

func (eventService *EventService) GetEventDiscounts(eventID uint) []dto.DiscountDetailsResponse {
	var notFoundError exceptions.NotFoundError
	_, eventExist := eventService.eventRepository.FindEventByID(eventID)
	if !eventExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Event
		panic(notFoundError)
	}

	discounts, discountExist := eventService.eventRepository.FindDiscountsByEventID(eventID)
	if !discountExist {
		return []dto.DiscountDetailsResponse{}
	}
	discountsDetails := make([]dto.DiscountDetailsResponse, len(discounts))
	for i, discount := range discounts {
		discountsDetails[i] = dto.DiscountDetailsResponse{
			ID:             discount.ID,
			CreatedAt:      discount.CreatedAt,
			Code:           discount.Code,
			Type:           discount.Type.String(),
			Value:          discount.Value,
			AvailableFrom:  discount.ValidFrom,
			AvailableUntil: discount.ValidUntil,
			Quantity:       discount.Quantity,
			UsedCount:      discount.UsedCount,
			MinTickets:     discount.MinTickets,
		}
	}
	return discountsDetails
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

func (eventService *EventService) DeleteTicket(eventID, ticketID uint) {
	var notFoundError exceptions.NotFoundError
	eventExist := eventService.eventRepository.DeleteTicket(eventID, ticketID)
	if !eventExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Ticket
		panic(notFoundError)
	}
}

func (eventService *EventService) DeleteDiscount(eventID, discountID uint) {
	var notFoundError exceptions.NotFoundError
	eventExist := eventService.eventRepository.DeleteDiscount(eventID, discountID)
	if !eventExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Discount
		panic(notFoundError)
	}
}

func (eventService *EventService) GetEventMediaPath(mediaID, eventID uint) entities.Media {
	var notFoundError exceptions.NotFoundError
	media, mediaExist := eventService.eventRepository.FindMediaByIDAndEventID(mediaID, eventID)
	if !mediaExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Media
		panic(notFoundError)
	}
	return media
}

func (eventService *EventService) DeleteEventMedia(mediaID uint) {
	var notFoundError exceptions.NotFoundError
	eventExist := eventService.eventRepository.DeleteMedia(mediaID)
	if !eventExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Media
		panic(notFoundError)
	}
}

func (eventService *EventService) ChangeEventStatus(eventID uint, newStatus string) {
	var notFoundError exceptions.NotFoundError
	var conflictError exceptions.ConflictError
	event, eventExist := eventService.eventRepository.FindEventByID(eventID)
	if !eventExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Event
		panic(notFoundError)
	}
	var enumNewStatus enums.EventStatus
	eventStatuses := enums.GetAllEventStatus()
	for _, eventStatus := range eventStatuses {
		if eventStatus.String() == newStatus {
			enumNewStatus = eventStatus
		}
	}
	if event.Status == enumNewStatus {
		conflictError.AppendError(
			eventService.constants.ErrorField.EventStatus,
			eventService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}
	eventService.eventRepository.ChangeStatusByEvent(event, enumNewStatus)
}

func (eventService *EventService) ValidateNewEventMediaDetails(eventID uint, mediaName string) {
	var notFoundError exceptions.NotFoundError
	var conflictError exceptions.ConflictError
	_, eventExist := eventService.eventRepository.FindEventByID(eventID)
	if !eventExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Event
		panic(notFoundError)
	}
	_, mediaExist := eventService.eventRepository.FindEventMediaByName(mediaName, eventID)
	if mediaExist {
		conflictError.AppendError(
			eventService.constants.ErrorField.Ticket,
			eventService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}
}

func (eventService *EventService) CreateEventMedia(mediaName, mediaPath string, eventID uint) entities.Media {
	eventMediaModel := entities.Media{
		Name:    mediaName,
		Path:    mediaPath,
		EventID: eventID,
	}
	media := eventService.eventRepository.CreateNewMedia(eventMediaModel)
	return media
}
