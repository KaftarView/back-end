package application

import (
	application_aws "first-project/src/application/aws"
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"
	"first-project/src/dto"
	"first-project/src/entities"
	"first-project/src/enums"
	"first-project/src/exceptions"
	repository_database "first-project/src/repository/database"
	repository_database_interfaces "first-project/src/repository/database/interfaces"
	"math/rand/v2"
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

type eventService struct {
	constants          *bootstrap.Constants
	awsS3Service       *application_aws.S3service
	categoryService    application_interfaces.CategoryService
	eventRepository    repository_database_interfaces.EventRepository
	commentRepository  repository_database_interfaces.CommentRepository
	purchaseRepository repository_database_interfaces.PurchaseRepository
	db                 *gorm.DB
}

func NewEventService(
	constants *bootstrap.Constants,
	awsService *application_aws.S3service,
	categoryService application_interfaces.CategoryService,
	eventRepository repository_database_interfaces.EventRepository,
	commentRepository repository_database_interfaces.CommentRepository,
	purchaseRepository repository_database_interfaces.PurchaseRepository,
	db *gorm.DB,
) *eventService {
	return &eventService{
		constants:          constants,
		awsS3Service:       awsService,
		categoryService:    categoryService,
		eventRepository:    eventRepository,
		commentRepository:  commentRepository,
		purchaseRepository: purchaseRepository,
		db:                 db,
	}
}

func (eventService *eventService) FetchEventByID(eventID uint) *entities.Event {
	var notFoundError exceptions.NotFoundError
	event, eventExist := eventService.eventRepository.FindEventByID(eventService.db, eventID)
	if !eventExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Event
		panic(notFoundError)
	}
	return event
}

func (eventService *eventService) fetchTicketByID(ticketID uint) *entities.Ticket {
	var notFoundError exceptions.NotFoundError
	ticket, ticketExist := eventService.eventRepository.FindEventTicketByID(eventService.db, ticketID)
	if !ticketExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Ticket
		panic(notFoundError)
	}
	return ticket
}

func (eventService *eventService) fetchDiscountByID(discountID uint) *entities.Discount {
	var notFoundError exceptions.NotFoundError
	discount, eventExist := eventService.eventRepository.FindDiscountByDiscountID(eventService.db, discountID)
	if !eventExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Event
		panic(notFoundError)
	}
	return discount
}

func (eventService *eventService) fetchMediaByID(mediaID uint) *entities.Media {
	var notFoundError exceptions.NotFoundError
	media, mediaExist := eventService.eventRepository.FindMediaByID(eventService.db, mediaID)
	if !mediaExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Media
		panic(notFoundError)
	}
	return media
}

func (eventService *eventService) validateUniqueTicketName(name string, eventID uint) {
	var conflictError exceptions.ConflictError
	_, ticketExist := eventService.eventRepository.FindEventTicketByName(eventService.db, name, eventID)
	if ticketExist {
		conflictError.AppendError(
			eventService.constants.ErrorField.Ticket,
			eventService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}
}

func (eventService *eventService) validateUniqueDiscountCode(code string, eventID uint) {
	var conflictError exceptions.ConflictError
	_, discountExist := eventService.eventRepository.FindEventDiscountByCode(eventService.db, code, eventID)
	if discountExist {
		conflictError.AppendError(
			eventService.constants.ErrorField.Discount,
			eventService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}
}

func (eventService *eventService) validateUniqueMediaName(name string, eventID uint) {
	var conflictError exceptions.ConflictError
	_, mediaExist := eventService.eventRepository.FindEventMediaByName(eventService.db, name, eventID)
	if mediaExist {
		conflictError.AppendError(
			eventService.constants.ErrorField.Media,
			eventService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}
}

func (eventService *eventService) setEventBannerPath(event *entities.Event, banner *multipart.FileHeader) {
	bannerPath := eventService.constants.S3Service.GetEventBannerKey(event.ID, banner.Filename)
	eventService.awsS3Service.UploadObject(enums.EventsBucket, bannerPath, banner)
	event.BannerPath = bannerPath
}

func (eventService *eventService) setMediaFilePath(media *entities.Media, file *multipart.FileHeader) {
	mediaPath := eventService.constants.S3Service.GetEventSessionKey(media.EventID, media.ID, file.Filename)
	eventService.awsS3Service.UploadObject(enums.EventsBucket, mediaPath, file)
	media.Path = mediaPath
}

func (eventService *eventService) ValidateEventCreationDetails(
	name, venueType, location string, fromDate, toDate time.Time,
) {
	var conflictError exceptions.ConflictError
	_, eventExist := eventService.eventRepository.FindDuplicatedEvent(eventService.db, name, venueType, location, fromDate, toDate)
	if eventExist {
		conflictError.AppendError(
			eventService.constants.ErrorField.Tittle,
			eventService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}
}

func (eventService *eventService) CreateEvent(eventDetails dto.CreateEventRequest) *entities.Event {
	enumStatus := enums.Draft
	eventStatuses := enums.GetAllEventStatus()
	for _, eventStatus := range eventStatuses {
		if eventStatus.String() == eventDetails.Status {
			enumStatus = eventStatus
			break
		}
	}

	enumVenue := enums.Online
	eventVenues := enums.GetAllEventVenues()
	for _, eventVenue := range eventVenues {
		if eventVenue.String() == eventDetails.VenueType {
			enumVenue = eventVenue
			break
		}
	}

	categories := eventService.categoryService.GetCategoriesByName(eventDetails.Categories)

	var event *entities.Event
	err := repository_database.ExecuteInTransaction(eventService.db, func(tx *gorm.DB) error {
		commentable := eventService.commentRepository.CreateNewCommentable(tx)

		bannerPath := eventService.constants.S3Service.GetEventBannerKey(commentable.CID, eventDetails.Banner.Filename)
		eventService.awsS3Service.UploadObject(enums.EventsBucket, bannerPath, eventDetails.Banner)

		event = &entities.Event{
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
		if err := eventService.eventRepository.CreateNewEvent(tx, event); err != nil {
			panic(err)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return event
}

func (eventService *eventService) CreateEventTicket(ticketDetails dto.CreateTicketRequest) *entities.Ticket {
	eventService.FetchEventByID(ticketDetails.EventID)
	eventService.validateUniqueTicketName(ticketDetails.Name, ticketDetails.EventID)

	ticket := &entities.Ticket{
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
	err := repository_database.ExecuteInTransaction(eventService.db, func(tx *gorm.DB) error {
		if err := eventService.eventRepository.CreateNewTicket(tx, ticket); err != nil {
			panic(err)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
	return ticket
}

func (eventService *eventService) UpdateEventTicket(ticketID uint, ticketDetails dto.CreateTicketRequest) {
	ticket := eventService.fetchTicketByID(ticketID)

	err := repository_database.ExecuteInTransaction(eventService.db, func(tx *gorm.DB) error {
		if ticket.Name != ticketDetails.Name {
			eventService.validateUniqueTicketName(ticketDetails.Name, ticket.EventID)
		}
		ticket.Name = ticketDetails.Name
		ticket.Description = ticketDetails.Description
		ticket.Price = ticketDetails.Price
		ticket.Quantity = ticketDetails.Quantity
		ticket.SoldCount = ticketDetails.SoldCount
		ticket.IsAvailable = ticketDetails.IsAvailable
		ticket.AvailableFrom = ticketDetails.AvailableFrom
		ticket.AvailableUntil = ticketDetails.AvailableUntil

		if err := eventService.eventRepository.UpdateEventTicket(tx, ticket); err != nil {
			panic(err)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (eventService *eventService) CreateEventDiscount(discountDetails dto.CreateDiscountRequest) *entities.Discount {
	eventService.FetchEventByID(discountDetails.EventID)
	eventService.validateUniqueDiscountCode(discountDetails.Code, discountDetails.EventID)

	var enumDiscountType enums.DiscountType
	discountTypes := enums.GetAllDiscountTypes()
	for _, discountType := range discountTypes {
		if discountType.String() == discountDetails.Type {
			enumDiscountType = discountType
			break
		}
	}

	discount := &entities.Discount{
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
	err := repository_database.ExecuteInTransaction(eventService.db, func(tx *gorm.DB) error {
		if err := eventService.eventRepository.CreateNewDiscount(tx, discount); err != nil {
			panic(err)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
	return discount
}

func (eventService *eventService) UpdateEventDiscount(discountID uint, discountDetails dto.CreateDiscountRequest) {
	discount := eventService.fetchDiscountByID(discountID)

	err := repository_database.ExecuteInTransaction(eventService.db, func(tx *gorm.DB) error {
		if discountDetails.Code != discount.Code {
			eventService.validateUniqueDiscountCode(discountDetails.Code, discount.EventID)
		}

		discount.Code = discountDetails.Code
		discount.Value = discountDetails.Value
		discount.ValidFrom = discountDetails.ValidFrom
		discount.ValidUntil = discountDetails.ValidUntil
		discount.Quantity = discountDetails.Quantity
		discount.UsedCount = discountDetails.UsedCount

		if err := eventService.eventRepository.UpdateEventDiscount(tx, discount); err != nil {
			panic(err)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (eventService *eventService) UpdateEvent(updateDetails dto.UpdateEventRequest) {
	event := eventService.FetchEventByID(updateDetails.ID)

	err := repository_database.ExecuteInTransaction(eventService.db, func(tx *gorm.DB) error {
		updateField(updateDetails.Description, &event.Description)
		updateField(updateDetails.FromDate, &event.FromDate)
		updateField(updateDetails.ToDate, &event.ToDate)
		updateField(updateDetails.BasePrice, &event.BasePrice)
		updateField(updateDetails.MinCapacity, &event.MinCapacity)
		updateField(updateDetails.MaxCapacity, &event.MaxCapacity)

		if updateDetails.Banner != nil {
			eventService.awsS3Service.DeleteObject(enums.EventsBucket, event.BannerPath)
			eventService.setEventBannerPath(event, updateDetails.Banner)
		}
		event.Status = updateEnumField(event.Status, updateDetails.Status, enums.GetAllEventStatus)
		event.VenueType = updateEnumField(event.VenueType, updateDetails.VenueType, enums.GetAllEventVenues)

		if updateDetails.Categories != nil {
			categories := eventService.categoryService.GetCategoriesByName(*updateDetails.Categories)
			if err := eventService.eventRepository.UpdateEventCategories(tx, updateDetails.ID, categories); err != nil {
				panic(err)
			}
		}

		if updateDetails.Name != nil {
			eventService.ValidateEventCreationDetails(*updateDetails.Name, event.VenueType.String(), event.Location, event.FromDate, event.ToDate)
			event.Name = *updateDetails.Name
		}

		if err := eventService.eventRepository.UpdateEvent(tx, event); err != nil {
			panic(err)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (eventService *eventService) CreateEventOrganizer(eventID uint, name, email, description string, profile *multipart.FileHeader) {
	var conflictError exceptions.ConflictError
	eventService.FetchEventByID(eventID)

	_, organizerExist := eventService.eventRepository.FindOrganizerByEmail(eventService.db, eventID, email)
	if organizerExist {
		conflictError.AppendError(
			eventService.constants.ErrorField.Organizer,
			eventService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}

	organizer := &entities.Organizer{
		Name:        name,
		Email:       email,
		Description: description,
		EventID:     eventID,
	}
	err := repository_database.ExecuteInTransaction(eventService.db, func(tx *gorm.DB) error {
		if err := eventService.eventRepository.CreateOrganizerForEventID(tx, organizer); err != nil {
			panic(err)
		}

		profilePath := eventService.constants.S3Service.GetOrganizerProfileKey(organizer.ID, profile.Filename)
		eventService.awsS3Service.UploadObject(enums.ProfilesBucket, profilePath, profile)
		organizer.ProfilePath = profilePath

		if err := eventService.eventRepository.UpdateEventOrganizer(tx, organizer); err != nil {
			panic(err)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (eventService *eventService) GetEventsList(allowedStatus []enums.EventStatus, page, pageSize int) []dto.EventDetailsResponse {
	offset := (page - 1) * pageSize
	events, _ := eventService.eventRepository.FindEventsByStatus(eventService.db, allowedStatus, offset, pageSize)
	eventsDetails := make([]dto.EventDetailsResponse, len(events))
	for i, event := range events {
		categories := eventService.eventRepository.FindEventCategoriesByEvent(eventService.db, event)
		categoryNames := make([]string, len(categories))
		for i, category := range categories {
			categoryNames[i] = category.Name
		}
		banner := eventService.awsS3Service.GetPresignedURL(enums.EventsBucket, event.BannerPath, 8*time.Hour)
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
			BasePrice:   event.BasePrice,
			Categories:  categoryNames,
		}
	}
	return eventsDetails
}

func (eventService *eventService) GetEventDetails(allowedStatus []enums.EventStatus, eventID uint) dto.EventDetailsResponse {
	var notFoundError exceptions.NotFoundError
	event, eventExist := eventService.eventRepository.FindEventByID(eventService.db, eventID)
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
	categories := eventService.eventRepository.FindEventCategoriesByEvent(eventService.db, event)
	categoryNames := make([]string, len(categories))
	for i, category := range event.Categories {
		categoryNames[i] = category.Name
	}

	banner := eventService.awsS3Service.GetPresignedURL(enums.EventsBucket, event.BannerPath, 8*time.Hour)

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

func (eventService *eventService) GetAvailableEventTickets(eventID uint) []dto.TicketDetailsResponse {
	eventService.FetchEventByID(eventID)

	tickets, ticketExist := eventService.eventRepository.FindAvailableTicketsByEventID(eventService.db, eventID)
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
			RemainTickets:  ticket.Quantity - ticket.SoldCount,
			Quantity:       ticket.Quantity,
			IsAvailable:    ticket.IsAvailable,
			AvailableFrom:  ticket.AvailableFrom,
			AvailableUntil: ticket.AvailableUntil,
		}
	}
	return ticketsDetails
}

func (eventService *eventService) GetAllEventTickets(eventID uint) []dto.TicketDetailsResponse {
	eventService.FetchEventByID(eventID)

	tickets, ticketExist := eventService.eventRepository.FindAllTicketsByEventID(eventService.db, eventID)
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
			RemainTickets:  ticket.Quantity - ticket.SoldCount,
			Quantity:       ticket.Quantity,
			IsAvailable:    ticket.IsAvailable,
			AvailableFrom:  ticket.AvailableFrom,
			AvailableUntil: ticket.AvailableUntil,
		}
	}
	return ticketsDetails
}

func (eventService *eventService) GetTicketDetails(ticketID uint) dto.TicketDetailsResponse {
	ticket, ticketExist := eventService.eventRepository.FindEventTicketByID(eventService.db, ticketID)
	if !ticketExist {
		return dto.TicketDetailsResponse{}
	}
	ticketDetails := dto.TicketDetailsResponse{
		ID:             ticketID,
		CreatedAt:      ticket.CreatedAt,
		Name:           ticket.Name,
		Description:    ticket.Description,
		Price:          ticket.Price,
		RemainTickets:  ticket.Quantity - ticket.SoldCount,
		Quantity:       ticket.Quantity,
		IsAvailable:    ticket.IsAvailable,
		AvailableFrom:  ticket.AvailableFrom,
		AvailableUntil: ticket.AvailableUntil,
	}

	return ticketDetails
}

func (eventService *eventService) GetEventDiscounts(eventID uint) []dto.DiscountDetailsResponse {
	eventService.FetchEventByID(eventID)

	discounts, discountExist := eventService.eventRepository.FindDiscountsByEventID(eventService.db, eventID)
	if !discountExist {
		return []dto.DiscountDetailsResponse{}
	}
	discountsDetails := make([]dto.DiscountDetailsResponse, len(discounts))
	for i, discount := range discounts {
		discountsDetails[i] = dto.DiscountDetailsResponse{
			ID:         discount.ID,
			CreatedAt:  discount.CreatedAt,
			Code:       discount.Code,
			Type:       discount.Type.String(),
			Value:      discount.Value,
			ValidFrom:  discount.ValidFrom,
			ValidUntil: discount.ValidUntil,
			Quantity:   discount.Quantity,
			UsedCount:  discount.UsedCount,
			MinTickets: discount.MinTickets,
		}
	}
	return discountsDetails
}

func (eventService *eventService) GetDiscountDetails(discountID uint) dto.DiscountDetailsResponse {
	discount := eventService.fetchDiscountByID(discountID)

	discountDetails := dto.DiscountDetailsResponse{
		ID:         discount.ID,
		CreatedAt:  discount.CreatedAt,
		Code:       discount.Code,
		Type:       discount.Type.String(),
		Value:      discount.Value,
		ValidFrom:  discount.ValidFrom,
		ValidUntil: discount.ValidUntil,
		Quantity:   discount.Quantity,
		UsedCount:  discount.UsedCount,
		MinTickets: discount.MinTickets,
	}

	return discountDetails
}

func (eventService *eventService) GetEventAttendees(eventID uint) []dto.EventAttendeesResponse {
	eventService.FetchEventByID(eventID)
	orders := eventService.eventRepository.FindOrdersEventID(eventService.db, eventID)
	var attendees []dto.EventAttendeesResponse
	for _, order := range orders {
		for _, item := range order.Items {
			attendant := dto.EventAttendeesResponse{
				Name:         order.User.Name,
				Email:        order.User.Email,
				Ticket:       item.TicketName,
				CountTickets: item.TicketQuantity,
				Price:        item.TicketPrice,
			}
			attendees = append(attendees, attendant)
		}
	}
	return attendees
}

func (eventService *eventService) DeleteEvent(eventID uint) {
	event := eventService.FetchEventByID(eventID)

	eventMedia, _ := eventService.eventRepository.FindAllEventMedia(eventService.db, eventID)
	eventOrganizers, _ := eventService.eventRepository.FindAllEventOrganizers(eventService.db, eventID)

	err := repository_database.ExecuteInTransaction(eventService.db, func(tx *gorm.DB) error {
		if err := eventService.eventRepository.DeleteEvent(tx, eventID); err != nil {
			panic(err)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	eventService.awsS3Service.DeleteObject(enums.EventsBucket, event.BannerPath)
	for _, organizer := range eventOrganizers {
		eventService.awsS3Service.DeleteObject(enums.ProfilesBucket, organizer.ProfilePath)
	}
	for _, media := range eventMedia {
		eventService.awsS3Service.DeleteObject(enums.EventsBucket, media.Path)
	}
}

func (eventService *eventService) DeleteTicket(ticketID uint) {
	eventService.fetchTicketByID(ticketID)

	err := repository_database.ExecuteInTransaction(eventService.db, func(tx *gorm.DB) error {
		if err := eventService.eventRepository.DeleteTicket(tx, ticketID); err != nil {
			panic(err)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (eventService *eventService) DeleteDiscount(discountID uint) {
	eventService.fetchDiscountByID(discountID)

	err := repository_database.ExecuteInTransaction(eventService.db, func(tx *gorm.DB) error {
		if err := eventService.eventRepository.DeleteDiscount(tx, discountID); err != nil {
			panic(err)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (eventService *eventService) DeleteOrganizer(organizerID uint) {
	var notFoundError exceptions.NotFoundError
	organizer, organizerExist := eventService.eventRepository.FindOrganizerByID(eventService.db, organizerID)
	if !organizerExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Organizer
		panic(notFoundError)
	}

	err := repository_database.ExecuteInTransaction(eventService.db, func(tx *gorm.DB) error {
		eventService.awsS3Service.DeleteObject(enums.ProfilesBucket, organizer.ProfilePath)
		if err := eventService.eventRepository.DeleteOrganizer(tx, organizerID); err != nil {
			return nil
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (eventService *eventService) GetEventMediaDetails(mediaID uint) dto.MediaDetailsResponse {
	media := eventService.fetchMediaByID(mediaID)

	mediaPath := eventService.awsS3Service.GetPresignedURL(enums.EventsBucket, media.Path, 8*time.Hour)
	mediaDetails := dto.MediaDetailsResponse{
		ID:        mediaID,
		Name:      media.Name,
		CreatedAt: media.CreatedAt,
		Size:      media.Size,
		Type:      media.Type,
		MediaPath: mediaPath,
	}

	return mediaDetails
}

func (eventService *eventService) GetListEventMedia(eventID uint) []dto.MediaDetailsResponse {
	eventService.FetchEventByID(eventID)

	allEventMedia, _ := eventService.eventRepository.FindAllEventMedia(eventService.db, eventID)
	allMediaDetails := make([]dto.MediaDetailsResponse, len(allEventMedia))
	for i, media := range allEventMedia {
		mediaPath := eventService.awsS3Service.GetPresignedURL(enums.EventsBucket, media.Path, 8*time.Hour)
		allMediaDetails[i] = dto.MediaDetailsResponse{
			ID:        media.ID,
			Name:      media.Name,
			CreatedAt: media.CreatedAt,
			Size:      media.Size,
			Type:      media.Type,
			MediaPath: mediaPath,
		}
	}
	return allMediaDetails
}

func (eventService *eventService) GetAttendantEventMedia(eventID, userID uint) []dto.MediaDetailsResponse {
	eventService.FetchEventByID(eventID)
	isUserAttended := eventService.eventRepository.IsUserAttendingEvent(eventService.db, userID, eventID)
	if !isUserAttended {
		panic(exceptions.NewForbiddenError())
	}
	allEventMedia, _ := eventService.eventRepository.FindAllEventMedia(eventService.db, eventID)
	allMediaDetails := make([]dto.MediaDetailsResponse, len(allEventMedia))
	for i, media := range allEventMedia {
		mediaPath := eventService.awsS3Service.GetPresignedURL(enums.EventsBucket, media.Path, 8*time.Hour)
		allMediaDetails[i] = dto.MediaDetailsResponse{
			ID:        media.ID,
			Name:      media.Name,
			CreatedAt: media.CreatedAt,
			Size:      media.Size,
			Type:      media.Type,
			MediaPath: mediaPath,
		}
	}
	return allMediaDetails
}

func (eventService *eventService) UpdateEventMedia(mediaID uint, name *string, file *multipart.FileHeader) {
	media := eventService.fetchMediaByID(mediaID)

	err := repository_database.ExecuteInTransaction(eventService.db, func(tx *gorm.DB) error {
		if name != nil {
			eventService.validateUniqueMediaName(*name, media.EventID)
			media.Name = *name
		}
		if file != nil {
			eventService.awsS3Service.DeleteObject(enums.EventsBucket, media.Path)
			eventService.setMediaFilePath(media, file)
		}

		if err := eventService.eventRepository.UpdateEventMedia(tx, media); err != nil {
			panic(err)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (eventService *eventService) DeleteEventMedia(mediaID uint) {
	media := eventService.fetchMediaByID(mediaID)

	err := repository_database.ExecuteInTransaction(eventService.db, func(tx *gorm.DB) error {
		eventService.awsS3Service.DeleteObject(enums.EventsBucket, media.Path)
		if err := eventService.eventRepository.DeleteMedia(tx, mediaID); err != nil {
			panic(err)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (eventService *eventService) ChangeEventStatus(eventID uint, newStatus string) {
	event := eventService.FetchEventByID(eventID)
	updateEnumField(event.Status, &newStatus, enums.GetAllEventStatus)
}

func (eventService *eventService) CreateEventMedia(eventID uint, mediaName string, mediaFile *multipart.FileHeader) {
	eventService.FetchEventByID(eventID)
	eventService.validateUniqueMediaName(mediaName, eventID)

	media := &entities.Media{
		Name:    mediaName,
		Size:    mediaFile.Size,
		Type:    mediaFile.Header.Get("Content-Type"),
		EventID: eventID,
	}
	err := repository_database.ExecuteInTransaction(eventService.db, func(tx *gorm.DB) error {
		if err := eventService.eventRepository.CreateNewMedia(tx, media); err != nil {
			panic(err)
		}

		eventService.setMediaFilePath(media, mediaFile)

		if err := eventService.eventRepository.UpdateEventMedia(tx, media); err != nil {
			panic(err)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (eventService *eventService) SearchEvents(query string, page, pageSize int, allowedStatus []enums.EventStatus) []dto.EventDetailsResponse {
	var events []*entities.Event
	offset := (page - 1) * pageSize
	if query != "" {
		events = eventService.eventRepository.FullTextSearch(eventService.db, query, allowedStatus, offset, pageSize)
	} else {
		events, _ = eventService.eventRepository.FindEventsByStatus(eventService.db, allowedStatus, offset, pageSize)
	}
	eventsDetails := make([]dto.EventDetailsResponse, len(events))
	for i, event := range events {
		categories := eventService.eventRepository.FindEventCategoriesByEvent(eventService.db, event)
		categoryNames := make([]string, len(categories))
		for i, category := range categories {
			categoryNames[i] = category.Name
		}
		banner := eventService.awsS3Service.GetPresignedURL(enums.EventsBucket, event.BannerPath, 8*time.Hour)
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
			BasePrice:   event.BasePrice,
			Categories:  categoryNames,
		}
	}
	return eventsDetails
}

func (eventService *eventService) FilterEventsByCategories(categories []string, page, pageSize int, allowedStatus []enums.EventStatus) []dto.EventDetailsResponse {
	var eventsList []*entities.Event
	offset := (page - 1) * pageSize
	if len(categories) == 0 {
		eventsList, _ = eventService.eventRepository.FindEventsByStatus(eventService.db, allowedStatus, offset, pageSize)
	} else {
		eventsList = eventService.eventRepository.FindEventsByCategoryName(eventService.db, categories, offset, pageSize, allowedStatus)
	}

	eventsDetails := make([]dto.EventDetailsResponse, len(eventsList))
	for i, event := range eventsList {
		categories := eventService.eventRepository.FindEventCategoriesByEvent(eventService.db, event)
		categoryNames := make([]string, len(categories))
		for i, category := range categories {
			categoryNames[i] = category.Name
		}
		banner := ""
		if event.BannerPath != "" {
			banner = eventService.awsS3Service.GetPresignedURL(enums.EventsBucket, event.BannerPath, 8*time.Hour)
		}
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
			BasePrice:   event.BasePrice,
			Categories:  categoryNames,
		}
	}

	return eventsDetails
}

func (eventService *eventService) validateTicketAvailability(
	ticket *entities.Ticket, ticketExist bool, inputTicket dto.ReserveTicketRequest) error {
	var notFoundError exceptions.NotFoundError
	if !ticketExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Ticket
		return notFoundError
	}
	if !ticket.IsAvailable ||
		ticket.SoldCount+inputTicket.Quantity > ticket.Quantity ||
		time.Now().Before(ticket.AvailableFrom) ||
		time.Now().After(ticket.AvailableUntil) {
		appError := exceptions.NewAppError(
			eventService.constants.ErrorField.Ticket,
			eventService.constants.ErrorTag.NotAvailable)
		return appError
	}
	return nil
}

func (eventService *eventService) validateDiscountAvailability(
	discount *entities.Discount, discountExist bool, totalRequestedTickets uint) error {
	var notFoundError exceptions.NotFoundError
	if !discountExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Discount
		return notFoundError
	}
	if discount.UsedCount >= discount.Quantity ||
		time.Now().Before(discount.ValidFrom) ||
		time.Now().After(discount.ValidUntil) ||
		discount.MinTickets > totalRequestedTickets {
		appError := exceptions.NewAppError(
			eventService.constants.ErrorField.Discount,
			eventService.constants.ErrorTag.NotAvailable)
		return appError
	}
	return nil
}

func applyDiscount(totalPrice float64, discount *entities.Discount) float64 {
	if discount == nil {
		return totalPrice
	}

	var totalPriceAfterDiscount float64
	switch discount.Type {
	case enums.Percentage:
		totalPriceAfterDiscount = totalPrice * (1 - discount.Value/100)
	case enums.Fixed:
		totalPriceAfterDiscount = totalPrice - discount.Value
	default:
		totalPriceAfterDiscount = totalPrice
	}

	if totalPriceAfterDiscount < 0 {
		totalPriceAfterDiscount = 0
	}

	return totalPriceAfterDiscount
}

func (eventService *eventService) ReserveEventTicket(
	userID, eventID uint, discountCode *string, tickets []dto.ReserveTicketRequest) dto.ReserveTicketResponse {
	var totalRequestedTickets uint = 0
	var totalPrice float64 = 0
	var discountID *uint
	var discountType *enums.DiscountType
	var discountValue *float64
	var reservation *entities.Reservation

	err := repository_database.ExecuteInTransaction(eventService.db, func(tx *gorm.DB) error {
		reservationItems := make([]*entities.ReservationItem, len(tickets))
		for i, inputTicket := range tickets {
			ticket, ticketExist := eventService.eventRepository.FindEventTicketByIDForUpdate(tx, inputTicket.ID)
			if err := eventService.validateTicketAvailability(ticket, ticketExist, inputTicket); err != nil {
				panic(err)
			}
			reservationItems[i] = &entities.ReservationItem{
				TicketID:       ticket.ID,
				TicketName:     ticket.Name,
				TicketPrice:    ticket.Price,
				TicketQuantity: inputTicket.Quantity,
			}
			totalPrice += (float64(inputTicket.Quantity) * ticket.Price)
			totalRequestedTickets += inputTicket.Quantity
			ticket.SoldCount += inputTicket.Quantity
			eventService.eventRepository.UpdateEventTicket(tx, ticket)
		}

		if discountCode != nil {
			discount, discountExist := eventService.eventRepository.FindEventDiscountByCodeForUpdate(tx, *discountCode, eventID)
			if err := eventService.validateDiscountAvailability(discount, discountExist, totalRequestedTickets); err != nil {
				panic(err)
			}
			discount.UsedCount++
			discountID = &discount.ID
			discountType = &discount.Type
			discountValue = &discount.Value
			totalPrice = applyDiscount(totalPrice, discount)
			eventService.eventRepository.UpdateEventDiscount(tx, discount)
		}

		reservation = &entities.Reservation{
			UserID:        userID,
			EventID:       eventID,
			Expiration:    time.Now().Add(15 * time.Minute),
			Status:        enums.Pending,
			TotalPrice:    totalPrice,
			DiscountID:    discountID,
			DiscountType:  discountType,
			DiscountValue: discountValue,
			Items:         reservationItems,
		}
		if err := eventService.purchaseRepository.CreateReservation(tx, reservation); err != nil {
			panic(err)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
	reserveTicket := dto.ReserveTicketResponse{
		ID:         reservation.ID,
		FinalPrice: reservation.TotalPrice,
	}
	return reserveTicket
}

func (eventService *eventService) PurchaseEventTicket(userID, eventID, reservationID uint) {
	var notFoundError exceptions.NotFoundError
	var conflictError exceptions.ConflictError
	eventService.FetchEventByID(eventID)
	reservation, reservationExist := eventService.purchaseRepository.GetReservationByID(eventService.db, reservationID)
	if !reservationExist {
		notFoundError.ErrorField = eventService.constants.ErrorField.Reservation
		panic(notFoundError)
	}
	if reservation.Status == enums.Confirmed {
		conflictError.AppendError(
			eventService.constants.ErrorField.Reservation,
			eventService.constants.ErrorTag.AlreadyPurchased)
		panic(conflictError)
	}

	err := repository_database.ExecuteInTransaction(eventService.db, func(tx *gorm.DB) error {
		transaction := &entities.Transaction{
			ReservationID: reservationID,
			UserID:        userID,
			Amount:        reservation.TotalPrice,
			GatewayName:   "MOCKED",
			TrackingID:    "MOCKED",
		}
		if rand.Float64() < 0.1 {
			appError := exceptions.NewAppError(
				eventService.constants.ErrorField.Reservation,
				eventService.constants.ErrorTag.PurchaseFailed)
			transaction.Status = enums.Failed
			if err := eventService.purchaseRepository.CreateTransaction(tx, transaction); err != nil {
				panic(err)
			}
			panic(appError)
		}

		transaction.Status = enums.Success
		if err := eventService.purchaseRepository.CreateTransaction(tx, transaction); err != nil {
			panic(err)
		}

		reservation.Status = enums.Confirmed
		eventService.purchaseRepository.UpdateReservation(tx, reservation)

		orderItems := make([]*entities.OrderItem, len(reservation.Items))
		for i, item := range reservation.Items {
			orderItems[i] = &entities.OrderItem{
				OrderID:        item.ID,
				TicketID:       item.TicketID,
				TicketName:     item.TicketName,
				TicketPrice:    item.TicketPrice,
				TicketQuantity: item.TicketQuantity,
			}
		}
		order := &entities.Order{
			UserID:           userID,
			EventID:          eventID,
			ReservationID:    &reservationID,
			TotalPrice:       reservation.TotalPrice,
			DiscountID:       reservation.DiscountID,
			DiscountType:     reservation.DiscountType,
			DiscountValue:    reservation.DiscountValue,
			PaymentMethod:    "MOCKED",
			PaymentReference: "MOCKED",
			Items:            orderItems,
		}
		if err := eventService.purchaseRepository.CreateOrder(tx, order); err != nil {
			panic(err)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (eventService *eventService) GetAllUserJoinedEvents(userID uint) []dto.EventDetailsResponse {
	orders := eventService.purchaseRepository.GetUserOrders(eventService.db, userID)
	eventsDetails := make([]dto.EventDetailsResponse, len(orders))
	for i, order := range orders {
		event := eventService.FetchEventByID(order.EventID)
		categories := eventService.eventRepository.FindEventCategoriesByEvent(eventService.db, event)
		categoryNames := make([]string, len(categories))
		for i, category := range categories {
			categoryNames[i] = category.Name
		}
		banner := eventService.awsS3Service.GetPresignedURL(enums.EventsBucket, event.BannerPath, 8*time.Hour)
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
			BasePrice:   event.BasePrice,
			Categories:  categoryNames,
		}
	}
	return eventsDetails
}
