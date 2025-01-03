package controller_v1_event

import (
	application_communication "first-project/src/application/communication/emailService"
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/dto"
	"first-project/src/enums"
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminEventController struct {
	constants    *bootstrap.Constants
	eventService application_interfaces.EventService
	emailService *application_communication.EmailService
}

func NewAdminEventController(
	constants *bootstrap.Constants,
	eventService application_interfaces.EventService,
	emailService *application_communication.EmailService,
) *AdminEventController {
	return &AdminEventController{
		constants:    constants,
		eventService: eventService,
		emailService: emailService,
	}
}

func (adminEventController *AdminEventController) GetEventsList(c *gin.Context) {
	pagination := controller.GetPagination(c, &adminEventController.constants.Context)
	allowedStatus := []enums.EventStatus{enums.Published, enums.Draft, enums.Completed, enums.Cancelled}
	events := adminEventController.eventService.GetEventsList(allowedStatus, pagination.Page, pagination.PageSize)

	controller.Response(c, 200, "", events)
}

func (adminEventController *AdminEventController) SearchEvents(c *gin.Context) {
	type searchEventsParams struct {
		Query string `form:"query"`
	}
	param := controller.Validated[searchEventsParams](c, &adminEventController.constants.Context)
	pagination := controller.GetPagination(c, &adminEventController.constants.Context)
	allowedStatus := []enums.EventStatus{enums.Published, enums.Draft, enums.Completed, enums.Cancelled}
	events := adminEventController.eventService.SearchEvents(param.Query, pagination.Page, pagination.PageSize, allowedStatus)

	controller.Response(c, 200, "", events)
}

func (adminEventController *AdminEventController) FilterEvents(c *gin.Context) {
	type filterEventsParams struct {
		Categories []string `form:"categories"`
	}
	param := controller.Validated[filterEventsParams](c, &adminEventController.constants.Context)
	pagination := controller.GetPagination(c, &adminEventController.constants.Context)
	allowedStatus := []enums.EventStatus{enums.Published, enums.Draft, enums.Completed, enums.Cancelled}
	events := adminEventController.eventService.FilterEventsByCategories(param.Categories, pagination.Page, pagination.PageSize, allowedStatus)

	controller.Response(c, 200, "", events)
}

func (adminEventController *AdminEventController) GetTicketDetails(c *gin.Context) {
	type getEventParams struct {
		TicketID uint `uri:"ticketID" validate:"required"`
	}
	param := controller.Validated[getEventParams](c, &adminEventController.constants.Context)
	ticketDetails := adminEventController.eventService.GetTicketDetails(param.TicketID)
	controller.Response(c, 200, "", ticketDetails)
}

func (adminEventController *AdminEventController) GetDiscountDetails(c *gin.Context) {
	type getEventParams struct {
		DiscountID uint `uri:"discountID" validate:"required"`
	}
	param := controller.Validated[getEventParams](c, &adminEventController.constants.Context)
	discountDetails := adminEventController.eventService.GetDiscountDetails(param.DiscountID)
	controller.Response(c, 200, "", discountDetails)
}

func (adminEventController *AdminEventController) GetMediaDetails(c *gin.Context) {
	type getEventParams struct {
		MediaID uint `uri:"mediaID" validate:"required"`
	}
	param := controller.Validated[getEventParams](c, &adminEventController.constants.Context)
	mediaDetails := adminEventController.eventService.GetEventMediaDetails(param.MediaID)
	controller.Response(c, 200, "", mediaDetails)
}

func (adminEventController *AdminEventController) GetEventDetails(c *gin.Context) {
	type getEventParams struct {
		EventID uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[getEventParams](c, &adminEventController.constants.Context)
	allowedStatus := []enums.EventStatus{enums.Published, enums.Draft, enums.Completed, enums.Cancelled}
	eventDetails := adminEventController.eventService.GetEventDetails(allowedStatus, param.EventID)
	controller.Response(c, 200, "", eventDetails)
}

func (adminEventController *AdminEventController) GetAllTicketDetails(c *gin.Context) {
	type getEventParams struct {
		EventID uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[getEventParams](c, &adminEventController.constants.Context)
	ticketDetails := adminEventController.eventService.GetEventTickets(param.EventID, []bool{true, false})
	controller.Response(c, 200, "", ticketDetails)
}

func (adminEventController *AdminEventController) GetAllDiscountDetails(c *gin.Context) {
	type getEventParams struct {
		EventID uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[getEventParams](c, &adminEventController.constants.Context)
	discountDetails := adminEventController.eventService.GetEventDiscounts(param.EventID)
	controller.Response(c, 200, "", discountDetails)
}

func (adminEventController *AdminEventController) GetEventMedia(c *gin.Context) {
	type getEventParams struct {
		EventID uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[getEventParams](c, &adminEventController.constants.Context)
	mediaDetails := adminEventController.eventService.GetListEventMedia(param.EventID)
	controller.Response(c, 200, "", mediaDetails)
}

func (adminEventController *AdminEventController) CreateEvent(c *gin.Context) {
	type createEventParams struct {
		Name        string                `form:"name" validate:"required,max=50"`
		Status      string                `form:"status"`
		Description string                `form:"description"`
		BasePrice   float64               `form:"basePrice" validate:"required"`
		FromDate    time.Time             `form:"fromDate" validate:"required"`
		ToDate      time.Time             `form:"toDate" validate:"required,gtfield=FromDate"`
		MinCapacity uint                  `form:"minCapacity" validate:"required,min=1"`
		MaxCapacity uint                  `form:"maxCapacity" validate:"required,gtfield=MinCapacity"`
		VenueType   string                `form:"venueType" validate:"required"`
		Location    string                `form:"location"`
		Banner      *multipart.FileHeader `form:"banner"`
		Categories  []string              `form:"categories"`
	}
	param := controller.Validated[createEventParams](c, &adminEventController.constants.Context)
	adminEventController.eventService.ValidateEventCreationDetails(
		param.Name, param.VenueType, param.Location, param.FromDate, param.ToDate,
	)

	eventDetails := dto.CreateEventRequest{
		Name:        param.Name,
		Status:      param.Status,
		Categories:  param.Categories,
		Description: param.Description,
		BasePrice:   param.BasePrice,
		FromDate:    param.FromDate,
		ToDate:      param.ToDate,
		MinCapacity: param.MinCapacity,
		MaxCapacity: param.MaxCapacity,
		VenueType:   param.VenueType,
		Location:    param.Location,
		Banner:      param.Banner,
	}

	event := adminEventController.eventService.CreateEvent(eventDetails)

	trans := controller.GetTranslator(c, adminEventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.createEvent")
	controller.Response(c, 200, message, event.ID)
}

func (adminEventController *AdminEventController) AddEventTicket(c *gin.Context) {
	type addEventTicketParams struct {
		Name           string    `json:"name" validate:"required,max=50"`
		Description    string    `json:"description"`
		Price          float64   `json:"price" validate:"required"`
		Quantity       uint      `json:"quantity" validate:"required"`
		SoldCount      uint      `json:"soldCount"`
		IsAvailable    bool      `json:"isAvailable" validate:"required"`
		AvailableFrom  time.Time `json:"availableFrom" validate:"required"`
		AvailableUntil time.Time `json:"availableUntil" validate:"required,gtfield=AvailableFrom"`
		EventID        uint      `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[addEventTicketParams](c, &adminEventController.constants.Context)

	ticketDetails := dto.CreateTicketRequest{
		Name:           param.Name,
		Description:    param.Description,
		Price:          param.Price,
		Quantity:       param.Quantity,
		SoldCount:      param.SoldCount,
		IsAvailable:    param.IsAvailable,
		AvailableFrom:  param.AvailableFrom,
		AvailableUntil: param.AvailableUntil,
		EventID:        param.EventID,
	}
	adminEventController.eventService.CreateEventTicket(ticketDetails)

	trans := controller.GetTranslator(c, adminEventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.addTicket")
	controller.Response(c, 200, message, nil)
}

func (adminEventController *AdminEventController) AddEventDiscount(c *gin.Context) {
	type addEventDiscountParams struct {
		Code       string    `json:"code" validate:"required,max=50"`
		Type       string    `json:"type" validate:"required"`
		Value      float64   `json:"value" validate:"required"`
		ValidFrom  time.Time `json:"validFrom" validate:"required"`
		ValidUntil time.Time `json:"validUntil" validate:"required,gtfield=ValidFrom"`
		Quantity   uint      `json:"quantity" validate:"required"`
		UsedCount  uint      `json:"usedCount"`
		MinTickets uint      `json:"minTickets"`
		EventID    uint      `uri:"eventID" validate:"required"`
	}

	param := controller.Validated[addEventDiscountParams](c, &adminEventController.constants.Context)

	discountDetails := dto.CreateDiscountRequest{
		Code:       param.Code,
		Type:       param.Type,
		Value:      param.Value,
		ValidFrom:  param.ValidFrom,
		ValidUntil: param.ValidUntil,
		Quantity:   param.Quantity,
		UsedCount:  param.UsedCount,
		MinTickets: param.MinTickets,
		EventID:    param.EventID,
	}
	adminEventController.eventService.CreateEventDiscount(discountDetails)

	trans := controller.GetTranslator(c, adminEventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.addDiscount")
	controller.Response(c, 200, message, nil)
}

func (adminEventController *AdminEventController) AddEventOrganizer(c *gin.Context) {
	type addEventOrganizerParams struct {
		Name        string                `form:"name" validate:"required,max=50"`
		Email       string                `form:"email" validate:"required,email"`
		Description string                `form:"description"`
		Profile     *multipart.FileHeader `form:"profile"`
		EventID     uint                  `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[addEventOrganizerParams](c, &adminEventController.constants.Context)
	adminEventController.eventService.CreateEventOrganizer(param.EventID, param.Name, param.Email, param.Description, param.Profile)

	eventName := adminEventController.eventService.FetchEventByID(param.EventID).Name
	emailTemplateData := struct {
		Name      string
		EventName string
	}{
		Name:      param.Name,
		EventName: eventName,
	}
	templatePath := controller.GetTemplatePath(c, adminEventController.constants.Context.Translator)
	adminEventController.emailService.SendEmail(
		param.Email, "Accept invitation", "acceptInvitation/"+templatePath, emailTemplateData)

	trans := controller.GetTranslator(c, adminEventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.organizerRegistration")
	controller.Response(c, 200, message, nil)
}

func (adminEventController *AdminEventController) UploadEventMedia(c *gin.Context) {
	type eventMedia struct {
		Name    string                `form:"name" validate:"required,max=50"`
		Media   *multipart.FileHeader `form:"media" validate:"required"`
		EventID uint                  `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[eventMedia](c, &adminEventController.constants.Context)
	adminEventController.eventService.CreateEventMedia(param.EventID, param.Name, param.Media)

	trans := controller.GetTranslator(c, adminEventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.uploadMedia")
	controller.Response(c, 200, message, nil)
}

func (adminEventController *AdminEventController) PublishEvent(c *gin.Context) {
	type publishEventParams struct {
		EventID uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[publishEventParams](c, &adminEventController.constants.Context)
	adminEventController.eventService.ChangeEventStatus(param.EventID, "Published")

	trans := controller.GetTranslator(c, adminEventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.publishEvent")
	controller.Response(c, 200, message, nil)
}

func (adminEventController *AdminEventController) UnpublishEvent(c *gin.Context) {
	type publishEventParams struct {
		EventID uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[publishEventParams](c, &adminEventController.constants.Context)
	adminEventController.eventService.ChangeEventStatus(param.EventID, "Draft")

	trans := controller.GetTranslator(c, adminEventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.unpublishEvent")
	controller.Response(c, 200, message, nil)
}

func (adminEventController *AdminEventController) UpdateEvent(c *gin.Context) {
	type updateEventParams struct {
		Name        *string               `form:"name" validate:"omitempty,max=50"`
		Status      *string               `form:"status"`
		Description *string               `form:"description"`
		FromDate    *time.Time            `form:"fromDate" validate:"omitempty"`
		ToDate      *time.Time            `form:"toDate" validate:"omitempty,gtfield=FromDate"`
		BasePrice   *float64              `form:"basePrice"`
		MinCapacity *uint                 `form:"minCapacity" validate:"omitempty,min=1"`
		MaxCapacity *uint                 `form:"maxCapacity" validate:"omitempty,gtfield=MinCapacity"`
		VenueType   *string               `form:"eventType" validate:"omitempty"`
		Location    *string               `form:"address"`
		Banner      *multipart.FileHeader `form:"banner"`
		Categories  *[]string             `form:"categories"`
		EventID     uint                  `uri:"eventID" binding:"required"`
	}

	param := controller.Validated[updateEventParams](c, &adminEventController.constants.Context)

	eventDetails := dto.UpdateEventRequest{
		ID:          param.EventID,
		Name:        param.Name,
		Status:      param.Status,
		Description: param.Description,
		BasePrice:   param.BasePrice,
		FromDate:    param.FromDate,
		ToDate:      param.ToDate,
		MinCapacity: param.MinCapacity,
		MaxCapacity: param.MaxCapacity,
		VenueType:   param.VenueType,
		Location:    param.Location,
		Categories:  param.Categories,
		Banner:      param.Banner,
	}

	adminEventController.eventService.UpdateEvent(eventDetails)

	trans := controller.GetTranslator(c, adminEventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.updateEvent")
	controller.Response(c, 200, message, nil)
}

func (adminEventController *AdminEventController) DeleteEvent(c *gin.Context) {
	type deleteEventParams struct {
		EventID uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[deleteEventParams](c, &adminEventController.constants.Context)
	adminEventController.eventService.DeleteEvent(param.EventID)

	trans := controller.GetTranslator(c, adminEventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteEvent")
	controller.Response(c, 200, message, nil)
}

func (adminEventController *AdminEventController) UpdateEventTicket(c *gin.Context) {
	type EditEventTicketParams struct {
		Name           *string    `json:"name"`
		Description    *string    `json:"description"`
		Price          *float64   `json:"price"`
		Quantity       *uint      `json:"quantity" `
		SoldCount      *uint      `json:"soldCount"`
		IsAvailable    *bool      `json:"isAvailable" `
		AvailableFrom  *time.Time `json:"availableFrom" `
		AvailableUntil *time.Time `json:"availableUntil" `
		TicketID       uint       `uri:"ticketID" validate:"required"`
	}
	param := controller.Validated[EditEventTicketParams](c, &adminEventController.constants.Context)
	ticketDetails := dto.UpdateTicketRequest{
		Name:           param.Name,
		Description:    param.Description,
		Price:          param.Price,
		Quantity:       param.Quantity,
		SoldCount:      param.SoldCount,
		IsAvailable:    param.IsAvailable,
		AvailableFrom:  param.AvailableFrom,
		AvailableUntil: param.AvailableUntil,
		TicketID:       param.TicketID,
	}
	adminEventController.eventService.UpdateEventTicket(ticketDetails)

	trans := controller.GetTranslator(c, adminEventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.updateTicket")
	controller.Response(c, 200, message, nil)
}

func (adminEventController *AdminEventController) DeleteTicket(c *gin.Context) {
	type deleteTicketParams struct {
		TicketID uint `uri:"ticketID" validate:"required"`
	}
	param := controller.Validated[deleteTicketParams](c, &adminEventController.constants.Context)
	adminEventController.eventService.DeleteTicket(param.TicketID)

	trans := controller.GetTranslator(c, adminEventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteTicket")
	controller.Response(c, 200, message, nil)
}

func (adminEventController *AdminEventController) UpdateEventDiscount(c *gin.Context) {
	type updateEventDiscountParams struct {
		Code       *string    `json:"code"`
		Type       *string    `json:"type"`
		Value      *float64   `json:"value"`
		ValidFrom  *time.Time `json:"validFrom"`
		ValidUntil *time.Time `json:"validUntil"`
		Quantity   *uint      `json:"quantity"`
		UsedCount  *uint      `json:"usedCount"`
		MinTickets *uint      `json:"minTickets"`
		DiscountID uint       `uri:"discountID" validate:"required"`
	}
	param := controller.Validated[updateEventDiscountParams](c, &adminEventController.constants.Context)

	discountDetails := dto.UpdateDiscountRequest{
		Code:           param.Code,
		Type:           param.Type,
		Value:          param.Value,
		AvailableFrom:  param.ValidFrom,
		AvailableUntil: param.ValidUntil,
		Quantity:       param.Quantity,
		UsedCount:      param.UsedCount,
		MinTickets:     param.MinTickets,
		DiscountID:     param.DiscountID,
	}
	adminEventController.eventService.UpdateEventDiscount(discountDetails)

	trans := controller.GetTranslator(c, adminEventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.updateDiscount")
	controller.Response(c, 200, message, nil)
}

func (adminEventController *AdminEventController) DeleteDiscount(c *gin.Context) {
	type deleteDiscountParams struct {
		DiscountID uint `uri:"discountID" validate:"required"`
	}
	param := controller.Validated[deleteDiscountParams](c, &adminEventController.constants.Context)
	adminEventController.eventService.DeleteDiscount(param.DiscountID)

	trans := controller.GetTranslator(c, adminEventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteDiscount")
	controller.Response(c, 200, message, nil)
}

func (adminEventController *AdminEventController) DeleteOrganizer(c *gin.Context) {
	type deleteOrganizerParams struct {
		OrganizerID uint `uri:"organizerID" validate:"required"`
	}
	param := controller.Validated[deleteOrganizerParams](c, &adminEventController.constants.Context)
	adminEventController.eventService.DeleteOrganizer(param.OrganizerID)

	trans := controller.GetTranslator(c, adminEventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteOrganizer")
	controller.Response(c, 200, message, nil)
}

func (adminEventController *AdminEventController) UpdateEventMedia(c *gin.Context) {
	type eventMedia struct {
		Name    *string               `form:"name" validate:"omitempty,max=50"`
		Media   *multipart.FileHeader `form:"media"`
		MediaID uint                  `uri:"mediaID" validate:"required"`
	}
	param := controller.Validated[eventMedia](c, &adminEventController.constants.Context)
	adminEventController.eventService.UpdateEventMedia(param.MediaID, param.Name, param.Media)

	trans := controller.GetTranslator(c, adminEventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.updateMedia")
	controller.Response(c, 200, message, nil)
}

func (adminEventController *AdminEventController) DeleteEventMedia(c *gin.Context) {
	type eventMedia struct {
		MediaID uint `uri:"mediaId" validate:"required"`
	}
	param := controller.Validated[eventMedia](c, &adminEventController.constants.Context)
	adminEventController.eventService.DeleteEventMedia(param.MediaID)

	trans := controller.GetTranslator(c, adminEventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteMedia")
	controller.Response(c, 200, message, nil)
}
