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
		Category    string                `form:"category"` // you can put some validate for this one
		Description string                `form:"description"`
		FromDate    time.Time             `form:"from-date" validate:"required"`
		ToDate      time.Time             `form:"to-date" validate:"required,gtfield=FromDate"`
		MinCapacity uint                  `form:"min-capacity" validate:"required,min=1"`
		MaxCapacity uint                  `form:"max-capacity" validate:"required,gtfield=MinCapacity"`
		VenueType   string                `form:"venue-type" validate:"required"`
		Location    string                `form:"location"`
		Banner      *multipart.FileHeader `form:"banner"`
	}
	param := controller.Validated[createEventParams](c, &eventController.constants.Context)
	eventController.eventService.ValidateEventCreationDetails(
		param.Name, param.VenueType, param.Location, param.FromDate, param.ToDate,
	)

	eventDetails := dto.CreateEventDetails{
		Name:        param.Name,
		Status:      param.Status,
		Category:    param.Category,
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
