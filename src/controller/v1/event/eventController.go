package controller_v1_event

import (
	"first-project/src/application"
	application_aws "first-project/src/application/aws"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/dto"
	"first-project/src/enums"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"
)

type EventController struct {
	constants    *bootstrap.Constants
	eventService *application.EventService
	awsService   *application_aws.S3service
}

func NewEventController(
	constants *bootstrap.Constants,
	eventService *application.EventService,
	awsService *application_aws.S3service,
) *EventController {
	return &EventController{
		constants:    constants,
		eventService: eventService,
		awsService:   awsService,
	}
}

func (eventController *EventController) GetEventsListForAdmin(c *gin.Context) {
	allowedStatus := []enums.EventStatus{enums.Published, enums.Draft, enums.Completed, enums.Cancelled}
	events := eventController.eventService.GetEventsList(allowedStatus)
	controller.Response(c, 200, "", events)
}

func (eventController *EventController) GetEventDetailsForAdmin(c *gin.Context) {
	type getEventParams struct {
		EventID uint `uri:"id" validate:"required"`
	}
	param := controller.Validated[getEventParams](c, &eventController.constants.Context)
	allowedStatus := []enums.EventStatus{enums.Published, enums.Draft, enums.Completed, enums.Cancelled}
	eventDetails := eventController.eventService.GetEventDetails(allowedStatus, param.EventID)
	controller.Response(c, 200, "", eventDetails)
}

func (eventController *EventController) GetTicketDetails(c *gin.Context) {
	type getEventParams struct {
		EventID uint `uri:"id" validate:"required"`
	}
	param := controller.Validated[getEventParams](c, &eventController.constants.Context)
	ticketDetails := eventController.eventService.GetEventTickets(param.EventID)
	controller.Response(c, 200, "", ticketDetails)
}

func (eventController *EventController) GetDiscountDetails(c *gin.Context) {
	type getEventParams struct {
		EventID uint `uri:"id" validate:"required"`
	}
	param := controller.Validated[getEventParams](c, &eventController.constants.Context)
	discountDetails := eventController.eventService.GetEventDiscounts(param.EventID)
	controller.Response(c, 200, "", discountDetails)
}

func (eventController *EventController) CreateEvent(c *gin.Context) {
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
		Categories  []string              `form:"category"`
	}
	param := controller.Validated[createEventParams](c, &eventController.constants.Context)
	eventController.eventService.ValidateEventCreationDetails(
		param.Name, param.VenueType, param.Location, param.FromDate, param.ToDate,
	)

	eventDetails := dto.CreateEventDetails{
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
	}

	event := eventController.eventService.CreateEvent(eventDetails)
	objectPath := fmt.Sprintf("events/%d/banners/%s", event.ID, param.Banner.Filename)
	eventController.awsService.UploadObject(enums.BannersBucket, objectPath, param.Banner)
	eventController.eventService.SetBannerPath(objectPath, event.ID)

	trans := controller.GetTranslator(c, eventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.createEvent")
	controller.Response(c, 200, message, event.ID)
}

func (eventController *EventController) AddEventTicket(c *gin.Context) {
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
	param := controller.Validated[addEventTicketParams](c, &eventController.constants.Context)
	eventController.eventService.ValidateNewEventTicketDetails(param.Name, param.EventID)

	ticketDetails := dto.CreateTicketDetails{
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
	eventController.eventService.CreateEventTicket(ticketDetails)

	trans := controller.GetTranslator(c, eventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.addTicket")
	controller.Response(c, 200, message, nil)
}

func (eventController *EventController) AddEventDiscount(c *gin.Context) {
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

	param := controller.Validated[addEventDiscountParams](c, &eventController.constants.Context)
	eventController.eventService.ValidateNewEventDiscountDetails(param.Code, param.EventID)

	discountDetails := dto.CreateDiscountDetails{
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
	eventController.eventService.CreateEventDiscount(discountDetails)

	trans := controller.GetTranslator(c, eventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.addDiscount")
	controller.Response(c, 200, message, nil)
}

func (eventController *EventController) UpdateEvent(c *gin.Context) {
	// some code here ...
}

func (eventController *EventController) DeleteEvent(c *gin.Context) {
	type deleteEventParams struct {
		EventID uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[deleteEventParams](c, &eventController.constants.Context)
	eventController.eventService.DeleteEvent(param.EventID)

	trans := controller.GetTranslator(c, eventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteEvent")
	controller.Response(c, 200, message, nil)
}

func (eventController *EventController) DeleteTicket(c *gin.Context) {
	type deleteTicketParams struct {
		TicketID uint `uri:"ticketID" validate:"required"`
		EventID  uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[deleteTicketParams](c, &eventController.constants.Context)
	eventController.eventService.DeleteTicket(param.EventID, param.TicketID)

	trans := controller.GetTranslator(c, eventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteTicket")
	controller.Response(c, 200, message, nil)
}

func (eventController *EventController) DeleteDiscount(c *gin.Context) {
	type deleteDiscountParams struct {
		DiscountID uint `uri:"discountID" validate:"required"`
		EventID    uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[deleteDiscountParams](c, &eventController.constants.Context)
	eventController.eventService.DeleteDiscount(param.EventID, param.DiscountID)

	trans := controller.GetTranslator(c, eventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteDiscount")
	controller.Response(c, 200, message, nil)
}

func (eventController *EventController) UploadEventMedia(c *gin.Context) {
	type eventMedia struct {
		Name    string                `form:"name" validate:"required,max=50"`
		Media   *multipart.FileHeader `form:"media" validate:"required"`
		EventID uint                  `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[eventMedia](c, &eventController.constants.Context)
	eventController.eventService.ValidateNewEventMediaDetails(param.EventID, param.Name)
	mediaPath := fmt.Sprintf("events/%d/sessions/%s", param.EventID, param.Media.Filename)
	eventController.awsService.UploadObject(enums.SessionsBucket, mediaPath, param.Media)
	eventController.eventService.CreateEventMedia(param.Name, mediaPath, param.EventID)

	trans := controller.GetTranslator(c, eventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.uploadMedia")
	controller.Response(c, 200, message, nil)
}

func (eventController *EventController) DeleteEventMedia(c *gin.Context) {
	type eventMedia struct {
		EventID uint `uri:"eventID" validate:"required"`
		MediaID uint `uri:"mediaId" validate:"required"`
	}
	param := controller.Validated[eventMedia](c, &eventController.constants.Context)
	media := eventController.eventService.GetEventMediaPath(param.MediaID, param.EventID)
	eventController.awsService.DeleteObject(enums.SessionsBucket, media.Path)
	eventController.eventService.DeleteEventMedia(param.MediaID)

	trans := controller.GetTranslator(c, eventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteMedia")
	controller.Response(c, 200, message, nil)
}

func (eventController *EventController) PublishEvent(c *gin.Context) {
	type publishEventParams struct {
		EventID uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[publishEventParams](c, &eventController.constants.Context)
	eventController.eventService.ChangeEventStatus(param.EventID, "Published")

	trans := controller.GetTranslator(c, eventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.publishEvent")
	controller.Response(c, 200, message, nil)
}

func (eventController *EventController) UnpublishEvent(c *gin.Context) {
	type publishEventParams struct {
		EventID uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[publishEventParams](c, &eventController.constants.Context)
	eventController.eventService.ChangeEventStatus(param.EventID, "Draft")

	trans := controller.GetTranslator(c, eventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.unpublishEvent")
	controller.Response(c, 200, message, nil)
}

func (eventController *EventController) ListPublicEvents(c *gin.Context) {
	allowedStatus := []enums.EventStatus{enums.Published}
	events := eventController.eventService.GetEventsList(allowedStatus)
	controller.Response(c, 200, "", events)
}

func (eventController *EventController) GetPublicEvent(c *gin.Context) {
	type getPublicEventParams struct {
		EventID uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[getPublicEventParams](c, &eventController.constants.Context)
	allowedStatus := []enums.EventStatus{enums.Published, enums.Completed}
	event := eventController.eventService.GetEventDetails(allowedStatus, param.EventID)
	controller.Response(c, 200, "", event)
}

func (ec *EventController) ListCategories(c *gin.Context) {
	categoryList := ec.eventService.GetListOfCategories()
	controller.Response(c, 200, "", categoryList)
}

func (eventController *EventController) SearchPublicEvents(c *gin.Context) {
	// some code here ...
}
