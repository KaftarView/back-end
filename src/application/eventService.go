package application

import (
	application_aws "first-project/src/application/aws"
	"first-project/src/bootstrap"
	"first-project/src/dto"
	"first-project/src/entities"
	"first-project/src/enums"
	"first-project/src/exceptions"
	repository_database "first-project/src/repository/database"
	"fmt"
	"mime/multipart"
	"time"
)

type EventService struct {
	constants         *bootstrap.Constants
	awsS3Service      *application_aws.S3service
	eventRepository   *repository_database.EventRepository
	commentRepository *repository_database.CommentRepository
}

func NewEventService(
	constants *bootstrap.Constants,
	awsService *application_aws.S3service,
	eventRepository *repository_database.EventRepository,
	commentRepository *repository_database.CommentRepository,
) *EventService {
	return &EventService{
		constants:         constants,
		awsS3Service:      awsService,
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

func (eventService *EventService) CreateEvent(eventDetails dto.RequestEventDetails) *entities.Event {
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

	bannerPath := fmt.Sprintf("banners/events/%d/images/%s", commentable.CID, eventDetails.Banner.Filename)
	eventService.awsS3Service.UploadObject(enums.BannersBucket, bannerPath, eventDetails.Banner)

	eventDetailsModel := &entities.Event{
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
		BannerPath:  bannerPath,
	}
	event := eventService.eventRepository.CreateNewEvent(eventDetailsModel)
	return event
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

func (eventService *EventService) CreateEventTicket(ticketDetails dto.CreateTicketDetails) *entities.Ticket {
	ticketDetailsModel := &entities.Ticket{
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

func (eventService *EventService) UpdateEventTicket(ticketDetails dto.EditTicketDetails) {
	var notFoundError exceptions.NotFoundError
	var conflictError exceptions.ConflictError
	ticket, ticketExist := eventService.eventRepository.FindEventTicketByID(ticketDetails.TicketID)
	if !ticketExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Ticket
		panic(notFoundError)
	}

	if ticketDetails.Name != nil {
		_, ticketExist := eventService.eventRepository.FindEventTicketByName(*ticketDetails.Name, ticket.EventID)
		if ticketExist {
			conflictError.AppendError(
				eventService.constants.ErrorField.Media,
				eventService.constants.ErrorTag.AlreadyExist)
			panic(conflictError)
		}
		ticket.Name = *ticketDetails.Name
	}
	if ticketDetails.Description != nil {
		ticket.Description = *ticketDetails.Description
	}
	if ticketDetails.Price != nil {
		ticket.Price = *ticketDetails.Price
	}
	if ticketDetails.Quantity != nil {
		ticket.Quantity = *ticketDetails.Quantity
	}
	if ticketDetails.SoldCount != nil {
		ticket.SoldCount = *ticketDetails.SoldCount
	}
	if ticketDetails.IsAvailable != nil {
		ticket.IsAvailable = *ticketDetails.IsAvailable
	}
	if ticketDetails.AvailableFrom != nil {
		ticket.AvailableFrom = *ticketDetails.AvailableFrom
	}
	if ticketDetails.AvailableUntil != nil {
		ticket.AvailableUntil = *ticketDetails.AvailableUntil
	}

	eventService.eventRepository.UpdateEventTicket(ticket)
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

func (eventService *EventService) CreateEventDiscount(discountDetails dto.CreateDiscountDetails) *entities.Discount {
	var enumDiscountType enums.DiscountType
	discountTypes := enums.GetAllDiscountTypes()
	for _, discountType := range discountTypes {
		if discountType.String() == discountDetails.Type {
			enumDiscountType = discountType
		}
	}

	discountDetailsModel := &entities.Discount{
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

func (eventService *EventService) UpdateEventDiscount(discountDetails dto.EditDiscountDetails) {
	var notFoundError exceptions.NotFoundError
	var conflictError exceptions.ConflictError

	discount, discountExist := eventService.eventRepository.FindDiscountByDiscountID(discountDetails.DiscountID)
	if !discountExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Discount
		panic(notFoundError)
	}

	if discountDetails.Code != nil {
		_, discountExist := eventService.eventRepository.FindEventDiscountByCode(*discountDetails.Code, discount.EventID)
		if discountExist {
			conflictError.AppendError(
				eventService.constants.ErrorField.Discount,
				eventService.constants.ErrorTag.AlreadyExist)
			panic(conflictError)
		}
		discount.Code = *discountDetails.Code
	}
	if discountDetails.Type != nil {
		enumDiscountType := discount.Type
		discountTypes := enums.GetAllDiscountTypes()
		for _, discountType := range discountTypes {
			if discountType.String() == *discountDetails.Type {
				enumDiscountType = discountType
				break
			}
		}
		discount.Type = enumDiscountType
	}
	if discountDetails.Value != nil {
		discount.Value = *discountDetails.Value
	}
	if discountDetails.AvailableFrom != nil {
		discount.ValidFrom = *discountDetails.AvailableFrom
	}

	if discountDetails.AvailableUntil != nil {
		discount.ValidUntil = *discountDetails.AvailableUntil
	}
	if discountDetails.Quantity != nil {
		discount.Quantity = *discountDetails.Quantity
	}
	if discountDetails.UsedCount != nil {
		discount.UsedCount = *discountDetails.UsedCount
	}

	eventService.eventRepository.UpdateEventDiscount(discount)
}

func updateBasicDetails(event *entities.Event, updateDetails dto.UpdateEventDetails) {
	if updateDetails.Description != nil {
		event.Description = *updateDetails.Description
	}
	if updateDetails.FromDate != nil {
		event.FromDate = *updateDetails.FromDate
	}
	if updateDetails.ToDate != nil {
		event.ToDate = *updateDetails.ToDate
	}
	if updateDetails.BasePrice != nil {
		event.BasePrice = *updateDetails.BasePrice
	}
	if updateDetails.MinCapacity != nil {
		event.MinCapacity = *updateDetails.MinCapacity
	}
	if updateDetails.MaxCapacity != nil {
		event.MaxCapacity = *updateDetails.MaxCapacity
	}
}

func (eventService *EventService) updateEventBanner(event *entities.Event, banner *multipart.FileHeader) {
	if banner != nil {
		eventService.awsS3Service.DeleteObject(enums.BannersBucket, event.BannerPath)
		bannerPath := fmt.Sprintf("profiles/events/%d/images/%s", event.ID, banner.Filename)
		eventService.awsS3Service.UploadObject(enums.BannersBucket, bannerPath, banner)
		event.BannerPath = bannerPath
	}
}

func (eventService *EventService) UpdateEvent(updateDetails dto.UpdateEventDetails) {
	event, eventExist := eventService.eventRepository.FindEventByID(updateDetails.ID)
	if !eventExist {
		var notFoundError exceptions.NotFoundError
		notFoundError.ErrorField = eventService.constants.ErrorField.Event
		panic(notFoundError)
	}

	eventService.updateEventBanner(event, updateDetails.Banner)
	updateBasicDetails(event, updateDetails)

	if updateDetails.Status != nil {
		statusEnum := event.Status
		for _, status := range enums.GetAllEventStatus() {
			if status.String() == *updateDetails.Status {
				statusEnum = status
				break
			}
		}
		event.Status = statusEnum
	}

	if updateDetails.VenueType != nil {
		venueEnum := event.VenueType
		for _, venue := range enums.GetAllEventVenues() {
			if venue.String() == *updateDetails.VenueType {
				venueEnum = venue
				break
			}
		}
		event.VenueType = venueEnum
	}

	if updateDetails.Categories != nil {
		categories := eventService.eventRepository.FindCategoriesByNames(*updateDetails.Categories)
		event.Categories = categories
	}

	if updateDetails.Name != nil {
		eventService.ValidateEventCreationDetails(*updateDetails.Name, event.VenueType.String(), event.Location, event.FromDate, event.ToDate)
		event.Name = *updateDetails.Name
	}

	eventService.eventRepository.UpdateEvent(event)
}

func (eventService *EventService) CreateEventOrganizer(eventID uint, name, email, description string, profile *multipart.FileHeader) {
	var notFoundError exceptions.NotFoundError
	var conflictError exceptions.ConflictError
	_, eventExist := eventService.eventRepository.FindEventByID(eventID)
	if !eventExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Event
		panic(notFoundError)
	}

	_, organizerExist := eventService.eventRepository.FindOrganizerByEmail(eventID, email)
	if organizerExist {
		conflictError.AppendError(
			eventService.constants.ErrorField.Organizer,
			eventService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}

	profilePath := fmt.Sprintf("profiles/organizers/%d/images/%s", eventID, profile.Filename)
	eventService.awsS3Service.UploadObject(enums.ProfileBucket, profilePath, profile)

	eventService.eventRepository.CreateOrganizerForEventID(eventID, name, email, description, profilePath)
}

func (eventService *EventService) GetEventByID(eventID uint) *entities.Event {
	event, _ := eventService.eventRepository.FindEventByID(eventID)
	return event
}

func (eventService *EventService) GetEventsList(allowedStatus []enums.EventStatus) []dto.EventDetailsResponse {
	events, _ := eventService.eventRepository.FindEventsByStatus(allowedStatus)
	eventsDetails := make([]dto.EventDetailsResponse, len(events))
	for i, event := range events {
		banner := eventService.awsS3Service.GetPresignedURL(enums.BannersBucket, event.BannerPath, 8*time.Hour)
		eventsDetails[i] = dto.EventDetailsResponse{
			ID:          event.ID,
			CreatedAt:   event.CreatedAt,
			Name:        event.Name,
			Status:      event.Status.String(),
			Description: event.Description,
			FromDate:    event.FromDate,
			ToDate:      event.ToDate,
			VenueType:   event.VenueType.String(),
			Banner:      banner,
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

	banner := eventService.awsS3Service.GetPresignedURL(enums.BannersBucket, event.BannerPath, 8*time.Hour)

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
		Banner:      banner,
	}
	return eventDetails
}

func (eventService *EventService) GetEventTickets(eventID uint, availability []bool) []dto.TicketDetailsResponse {
	var notFoundError exceptions.NotFoundError
	_, eventExist := eventService.eventRepository.FindEventByID(eventID)
	if !eventExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Event
		panic(notFoundError)
	}

	tickets, ticketExist := eventService.eventRepository.FindTicketsByEventID(eventID, availability)
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

func (eventService *EventService) GetTicketDetails(ticketID uint) dto.TicketDetailsResponse {
	ticket, ticketExist := eventService.eventRepository.FindEventTicketByID(ticketID)
	if !ticketExist {
		return dto.TicketDetailsResponse{}
	}
	ticketDetails := dto.TicketDetailsResponse{
		ID:             ticketID,
		CreatedAt:      ticket.CreatedAt,
		Name:           ticket.Name,
		Description:    ticket.Description,
		Price:          ticket.Price,
		Quantity:       ticket.Quantity,
		IsAvailable:    ticket.IsAvailable,
		AvailableFrom:  ticket.AvailableFrom,
		AvailableUntil: ticket.AvailableUntil,
	}

	return ticketDetails
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

func (eventService *EventService) GetDiscountDetails(discountID uint) dto.DiscountDetailsResponse {
	var notFoundError exceptions.NotFoundError
	discount, discountExist := eventService.eventRepository.FindDiscountByDiscountID(discountID)
	if !discountExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Discount
		panic(notFoundError)
	}
	discountDetails := dto.DiscountDetailsResponse{
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

	return discountDetails
}

func (eventService *EventService) GetListEventMedia(eventID uint) []*entities.Media {
	var notFoundError exceptions.NotFoundError
	media, mediaExist := eventService.eventRepository.FindAllEventMedia(eventID)
	if !mediaExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Media
		panic(notFoundError)
	}
	return media
}

func (eventService *EventService) GetListOfCategories() []string {
	categoryNames := eventService.eventRepository.FindAllCategories()
	return categoryNames
}

func (eventService *EventService) DeleteEvent(eventID uint) {
	var notFoundError exceptions.NotFoundError
	event, eventExist := eventService.eventRepository.FindEventByID(eventID)
	if !eventExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Event
		panic(notFoundError)
	}
	eventMedia, _ := eventService.eventRepository.FindAllEventMedia(eventID)
	eventOrganizers, _ := eventService.eventRepository.FindAllEventOrganizers(eventID)

	eventService.eventRepository.DeleteEvent(eventID)

	eventService.awsS3Service.DeleteObject(enums.BannersBucket, event.BannerPath)
	for _, organizer := range eventOrganizers {
		eventService.awsS3Service.DeleteObject(enums.ProfileBucket, organizer.ProfilePath)
	}
	for _, media := range eventMedia {
		eventService.awsS3Service.DeleteObject(enums.SessionsBucket, media.Path)
	}
}

func (eventService *EventService) DeleteTicket(ticketID uint) {
	var notFoundError exceptions.NotFoundError
	_, ticketExist := eventService.eventRepository.FindEventTicketByID(ticketID)
	if !ticketExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Ticket
		panic(notFoundError)
	}
	eventService.eventRepository.DeleteTicket(ticketID)
}

func (eventService *EventService) DeleteDiscount(discountID uint) {
	var notFoundError exceptions.NotFoundError
	_, discountExist := eventService.eventRepository.FindDiscountByDiscountID(discountID)
	if !discountExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Discount
		panic(notFoundError)
	}
	eventService.eventRepository.DeleteDiscount(discountID)
}

func (eventService *EventService) DeleteOrganizer(organizerID uint) {
	var notFoundError exceptions.NotFoundError
	organizer, organizerExist := eventService.eventRepository.FindOrganizerByID(organizerID)
	if !organizerExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Organizer
		panic(notFoundError)
	}

	eventService.awsS3Service.DeleteObject(enums.ProfileBucket, organizer.ProfilePath)
	eventService.eventRepository.DeleteOrganizer(organizerID)
}

func (eventService *EventService) GetEventMediaDetails(mediaID, eventID uint) *entities.Media {
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
	media, mediaExist := eventService.eventRepository.FindMediaByID(mediaID)
	if !mediaExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Media
		panic(notFoundError)
	}
	eventService.awsS3Service.DeleteObject(enums.SessionsBucket, media.Path)
	eventService.eventRepository.DeleteMedia(mediaID)
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

func (eventService *EventService) CreateEventMedia(eventID uint, mediaName string, mediaFile *multipart.FileHeader) {
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
	mediaPath := fmt.Sprintf("media/events/%d/resources/%s", eventID, mediaFile.Filename)
	eventService.awsS3Service.UploadObject(enums.SessionsBucket, mediaPath, mediaFile)

	mediaModel := &entities.Media{
		Name:    mediaName,
		Path:    mediaPath,
		EventID: eventID,
	}
	eventService.eventRepository.CreateNewMedia(mediaModel)
}

func (eventService *EventService) SearchEvents(query string, page, pageSize int, allowedStatus []enums.EventStatus) []dto.EventDetailsResponse {
	var events []*entities.Event
	offset := (page - 1) * pageSize
	if query != "" {
		events = eventService.eventRepository.FullTextSearch(query, allowedStatus, offset, pageSize)
	} else {
		events, _ = eventService.eventRepository.FindEventsByStatus(allowedStatus)
	}
	eventsDetails := make([]dto.EventDetailsResponse, len(events))
	for i, event := range events {
		banner := eventService.awsS3Service.GetPresignedURL(enums.BannersBucket, event.BannerPath, 8*time.Hour)
		eventsDetails[i] = dto.EventDetailsResponse{
			ID:          event.ID,
			CreatedAt:   event.CreatedAt,
			Name:        event.Name,
			Status:      event.Status.String(),
			Description: event.Description,
			FromDate:    event.FromDate,
			ToDate:      event.ToDate,
			VenueType:   event.VenueType.String(),
			Banner:      banner,
		}
	}
	return eventsDetails
}
