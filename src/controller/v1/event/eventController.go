package controller_v1_event

import (
	"first-project/src/application"
	application_aws "first-project/src/application/aws"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/dto"
	"fmt"
	"log"
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
		FromDate    time.Time             `form:"from-date" validate:"required"`
		ToDate      time.Time             `form:"to-date" validate:"required,gtfield=FromDate"`
		MinCapacity uint                  `form:"min-capacity" validate:"required,min=1"`
		MaxCapacity uint                  `form:"max-capacity" validate:"required,gtfield=MinCapacity"`
		VenueType   string                `form:"venue-type" validate:"required"`
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
	controller.Response(c, 200, message, nil)
}

func (eventController *EventController) AddEventTicket(c *gin.Context) {
	type addEventTicketParams struct {
		Name           string    `json:"name" validate:"required,max=50"`
		Description    string    `json:"description"`
		Price          float64   `json:"price" validate:"required"`
		Quantity       uint      `json:"quantity" validate:"required"`
		SoldCount      uint      `json:"sold-count"`
		IsAvailable    bool      `json:"is-available" validate:"required"`
		AvailableFrom  time.Time `json:"available-from" validate:"required"`
		AvailableUntil time.Time `json:"available-until" validate:"required,gtfield=AvailableFrom"`
		EventID        uint      `json:"event-id" validate:"required"`
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
		ValidFrom  time.Time `json:"valid_from" validate:"required"`
		ValidUntil time.Time `json:"valid_until" validate:"required,gtfield=ValidFrom"`
		Quantity   uint      `json:"quantity" validate:"required"`
		UsedCount  uint      `json:"used_count"`
		MinTickets uint      `json:"min_tickets"`
		EventID    uint      `json:"event_id" validate:"required"`
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

func (eventController *EventController) EditEvent(c *gin.Context) {
	type editEventParams struct {
		EventID uint `uri:"id" binding:"required"`
	}
	param := controller.Validated[editEventParams](c, &eventController.constants.Context)

	trans := controller.GetTranslator(c, eventController.constants.Context.Translator)
	event, found := eventController.eventService.GetEventById(param.EventID)
	if !found {
		message, _ := trans.T("errorMessage.notFoundError")
		controller.Response(c, 404, message, nil)
		return
	}

	type responseStruct struct {
		Name        string    `json:"name"`
		Status      string    `json:"status"`
		Description string    `json:"description"`
		FromDate    time.Time `json:"fromDate"`
		ToDate      time.Time `json:"toDate"`
		MinCapacity uint      `json:"minCapacity"`
		MaxCapacity uint      `json:"maxCapacity"`
		VenueType   string    `json:"eventType"`
		Categories  []string  `json:"category"`
		Address     string    `json:"address"`
	}

	response := responseStruct{
		Name:        event.Name,
		Status:      event.Status.String(),
		Description: event.Description,
		FromDate:    event.FromDate,
		ToDate:      event.ToDate,
		MinCapacity: event.MinCapacity,
		MaxCapacity: event.MaxCapacity,
		VenueType:   event.VenueType.String(),
		Categories:  []string{"Music", "Workshop", "Tech"},
		Address:     event.Location,
	}

	message, _ := trans.T("successMessage.getEvent")
	controller.Response(c, 200, message, response)
}

func (eventController *EventController) UpdateEvent(c *gin.Context) {
	type updateEventParams struct {
		Name        *string               `form:"name" validate:"omitempty,max=50"`
		Status      *string               `form:"status"`
		Description *string               `form:"description"`
		FromDate    *time.Time            `form:"fromDate" validate:"omitempty"`
		ToDate      *time.Time            `form:"toDate" validate:"omitempty,gtfield=FromDate"`
		MinCapacity *uint                 `form:"minCapacity" validate:"omitempty,min=1"`
		MaxCapacity *uint                 `form:"maxCapacity" validate:"omitempty,gtfield=MinCapacity"`
		VenueType   *string               `form:"eventType" validate:"omitempty"`
		Location    *string               `form:"address"`
		Banner      *multipart.FileHeader `form:"banner"`
		Categories  *[]string             `form:"category"`
		EventID     uint                  `uri:"id" binding:"required"`
	}

	param := controller.Validated[updateEventParams](c, &eventController.constants.Context)

	eventDetails := dto.UpdateEventDetails{
		ID:          param.EventID,
		Name:        param.Name,
		Status:      param.Status,
		Description: param.Description,
		FromDate:    param.FromDate,
		ToDate:      param.ToDate,
		MinCapacity: param.MinCapacity,
		MaxCapacity: param.MaxCapacity,
		VenueType:   param.VenueType,
		Location:    param.Location,
		Categories:  param.Categories,
	}

	eventController.eventService.UpdateEvent(eventDetails)

	if param.Banner != nil {
		log.Printf("Banner file received: %+v\n", param.Banner.Filename)
		objectPath := fmt.Sprintf("Events/Banners/%d", int(param.EventID))
		eventController.awsService.DeleteObject(objectPath)
		eventController.awsService.UploadObject(param.Banner, "Events/Banners", int(param.EventID))
	}

	trans := controller.GetTranslator(c, eventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.updateEvent")
	log.Println("Sending success response")
	controller.Response(c, 200, message, nil)
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
