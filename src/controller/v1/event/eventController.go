package controller_v1_event

import (
	"first-project/src/application"
	application_aws "first-project/src/application/aws"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/dto"
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"
)

type EventController struct {
	constants    *bootstrap.Constants
	eventService *application.EventService
	awsService   *application_aws.AWSS3
}

func NewEventController(
	constants *bootstrap.Constants,
	eventService *application.EventService,
	awsService *application_aws.AWSS3,
) *EventController {
	return &EventController{
		constants:    constants,
		eventService: eventService,
		awsService:   awsService,
	}
}

func (eventController *EventController) ListEvents(c *gin.Context) {
	// some code here ...
}

func (eventController *EventController) GetEvent(c *gin.Context) {
	// some code here ...
}

func (eventController *EventController) CreateEvent(c *gin.Context) {
	type createEventParams struct {
		Name        string                `form:"name" validate:"required,max=50"`
		Status      string                `form:"status"`
		Description string                `form:"description"`
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
		FromDate:    param.FromDate,
		ToDate:      param.ToDate,
		MinCapacity: param.MinCapacity,
		MaxCapacity: param.MaxCapacity,
		VenueType:   param.VenueType,
		Location:    param.Location,
	}

	event := eventController.eventService.CreateEvent(eventDetails)

	eventController.awsService.UploadObject(param.Banner, "Events/Banners", int(event.ID))

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
	// some code here ...
}

func (eventController *EventController) UploadEventMedia(c *gin.Context) {
	// some code here ...
}

func (eventController *EventController) DeleteEventMedia(c *gin.Context) {
	// some code here ...
}

func (eventController *EventController) PublishEvent(c *gin.Context) {
	// some code here ...
}

func (eventController *EventController) UnpublishEvent(c *gin.Context) {
	// some code here ...
}
func (eventController *EventController) ListPublicEvents(c *gin.Context) {
	events := eventController.eventService.GetListOfPublishedEvents()
	controller.Response(c, 200, "", events)
}

func (eventController *EventController) GetPublicEvent(c *gin.Context) {
	type getPublicEventParams struct {
		EventID uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[getPublicEventParams](c, &eventController.constants.Context)
	event := eventController.eventService.GetPublicEventDetails(param.EventID)
	controller.Response(c, 200, "", event)
}

func (ec *EventController) ListCategories(c *gin.Context) {
	categoryList := ec.eventService.GetListOfCategories()
	controller.Response(c, 200, "", categoryList)
}

func (eventController *EventController) SearchPublicEvents(c *gin.Context) {
	// some code here ...
}
