package controller_v1_event

import (
	"encoding/base64"
	"first-project/src/application"
	application_aws "first-project/src/application/aws"
	application_communication "first-project/src/application/communication/emailService"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/dto"
	"first-project/src/enums"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"
)

type EventController struct {
	constants    *bootstrap.Constants
	eventService *application.EventService
	awsService   *application_aws.S3service
	emailService *application_communication.EmailService
}

func NewEventController(
	constants *bootstrap.Constants,
	eventService *application.EventService,
	awsService *application_aws.S3service,
	emailService *application_communication.EmailService,
) *EventController {
	return &EventController{
		constants:    constants,
		eventService: eventService,
		awsService:   awsService,
		emailService: emailService,
	}
}

func getTemplatePath(c *gin.Context, transKey string) string {
	trans := controller.GetTranslator(c, transKey)
	if trans.Locale() == "fa_IR" {
		return "fa.html"
	}
	return "en.html"
}

func (eventController *EventController) GetEventsListForAdmin(c *gin.Context) {
	allowedStatus := []enums.EventStatus{enums.Published, enums.Draft, enums.Completed, enums.Cancelled}
	events := eventController.eventService.GetEventsList(allowedStatus)
	for i := range events {
		events[i].Banner = eventController.awsService.GetPresignedURL(enums.BannersBucket, events[i].Banner, 8*time.Hour)
	}
	controller.Response(c, 200, "", events)
}

func (eventController *EventController) GetEventDetailsForAdmin(c *gin.Context) {
	type getEventParams struct {
		EventID uint `uri:"id" validate:"required"`
	}
	param := controller.Validated[getEventParams](c, &eventController.constants.Context)
	allowedStatus := []enums.EventStatus{enums.Published, enums.Draft, enums.Completed, enums.Cancelled}
	eventDetails := eventController.eventService.GetEventDetails(allowedStatus, param.EventID)
	eventDetails.Banner = eventController.awsService.GetPresignedURL(enums.BannersBucket, eventDetails.Banner, 8*time.Hour)
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
	objectPath := fmt.Sprintf("banners/events/%d/images/%s", event.ID, param.Banner.Filename)
	eventController.awsService.UploadObject(enums.BannersBucket, objectPath, param.Banner)
	eventController.eventService.SetBannerPathForEvent(objectPath, event.ID)

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

func (eventController *EventController) EditEventTicket(c *gin.Context) {
	type Params struct {
		EventID  uint `uri:"eventID" validate:"required"`
		TicketID uint `uri:"ticketID" validate:"required"`
	}
	param := controller.Validated[Params](c, &eventController.constants.Context)
	ticketDetails := eventController.eventService.GetTicketDetails(param.TicketID)
	controller.Response(c, 200, "", ticketDetails)
}

func (eventController *EventController) UpdateEventTicket(c *gin.Context) {
	type EditEventTicketParams struct {
		Name           *string    `json:"name"`
		Description    *string    `json:"description"`
		Price          *float64   `json:"price"`
		Quantity       *uint      `json:"quantity" `
		SoldCount      *uint      `json:"soldCount"`
		IsAvailable    *bool      `json:"isAvailable" `
		AvailableFrom  *time.Time `json:"availableFrom" `
		AvailableUntil *time.Time `json:"availableUntil" `
		EventID        uint       `uri:"eventID" validate:"required"`
		TicketID       uint       `uri:"ticketID" validate:"required"`
	}
	param := controller.Validated[EditEventTicketParams](c, &eventController.constants.Context)
	ticketDetails := dto.EditTicketDetaitls{
		Name:           param.Name,
		Description:    param.Description,
		Price:          param.Price,
		Quantity:       param.Quantity,
		SoldCount:      param.SoldCount,
		IsAvailable:    param.IsAvailable,
		AvailableFrom:  param.AvailableFrom,
		AvailableUntil: param.AvailableUntil,
		EventID:        param.EventID,
		TicketID:       param.TicketID,
	}
	trans := controller.GetTranslator(c, eventController.constants.Context.Translator)
	eventController.eventService.UpdateEventTicket(ticketDetails)
	message, _ := trans.T("successMessage.addDiscount")
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
func (eventController *EventController) EditEventDiscount(c *gin.Context) {
	type EditDiscountParams struct {
		EventID    uint `uri:"eventID" validate:"required"`
		DiscountID uint `uri:"discountID" validate:"required"`
	}
	param := controller.Validated[EditDiscountParams](c, &eventController.constants.Context)
	discountDetails := eventController.eventService.GetDiscountDetails(param.DiscountID)
	controller.Response(c, 200, "", discountDetails)
}
func (eventController *EventController) UpdateEventDiscount(c *gin.Context) {
	type updateEventDiscountParams struct {
		Code       *string    `json:"code"`
		Type       *string    `json:"type"`
		Value      *float64   `json:"value"`
		ValidFrom  *time.Time `json:"validFrom"`
		ValidUntil *time.Time `json:"validUntil"`
		Quantity   *uint      `json:"quantity"`
		UsedCount  *uint      `json:"usedCount"`
		MinTickets *uint      `json:"minTickets"`
		EventID    uint       `uri:"eventID" validate:"required"`
		DiscountID uint       `uri:"discountID" validate:"required"`
	}
	param := controller.Validated[updateEventDiscountParams](c, &eventController.constants.Context)

	discountDetails := dto.EditDiscountDetails{
		Code:           param.Code,
		Type:           param.Type,
		Value:          param.Value,
		AvailableFrom:  param.ValidFrom,
		AvailableUntil: param.ValidUntil,
		Quantity:       param.Quantity,
		UsedCount:      param.UsedCount,
		MinTickets:     param.MinTickets,
		EventID:        param.EventID,
		DiscountID:     param.DiscountID,
	}
	eventController.eventService.UpdateEventDiscount(discountDetails)

	trans := controller.GetTranslator(c, eventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.updateDiscount")
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
func (eventController *EventController) AddEventOrganizer(c *gin.Context) {
	type addEventOrganizerParams struct {
		Name        string                `form:"name" validate:"required,max=50"`
		Email       string                `form:"email" validate:"required,email"`
		Description string                `form:"description"`
		Profile     *multipart.FileHeader `form:"profile"`
		EventID     uint                  `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[addEventOrganizerParams](c, &eventController.constants.Context)
	token := application.GenerateSecureToken(32)
	objectPath := fmt.Sprintf("profiles/organizers/%d/images/%s", param.EventID, param.Profile.Filename)
	eventController.awsService.UploadObject(enums.ProfileBucket, objectPath, param.Profile)
	organizerID := eventController.eventService.UpdateOrCreateEventOrganizer(param.EventID, param.Name, param.Email, param.Description, token)
	eventName := eventController.eventService.GetEventByID(param.EventID).Name
	encodedOrganizerID := base64.StdEncoding.EncodeToString([]byte(string(organizerID)))
	encodedEventID := base64.StdEncoding.EncodeToString([]byte(string(param.EventID)))
	emailTemplateData := struct {
		Name      string
		EventName string
		Link      string
	}{
		Name:      param.Name,
		EventName: eventName,
		Link:      encodedOrganizerID + "/" + encodedEventID + "/" + token,
	}
	eventController.eventService.SetProfilePathForOrganizer(objectPath, organizerID)
	templatePath := getTemplatePath(c, eventController.constants.Context.Translator)
	eventController.emailService.SendEmail(
		param.Email, "Accept invitation", "acceptInvitation/"+templatePath, emailTemplateData)

	trans := controller.GetTranslator(c, eventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.organizerRegistration")
	controller.Response(c, 200, message, nil)
}

func (eventController *EventController) VerifyEmail(c *gin.Context) {
	type verifyEmailParams struct {
		EncodedOrganizerID string `json:"encodedOrganizerID" validate:"required"`
		EncodedEventID     string `json:"encodedEventID" validate:"required"`
		Token              string `json:"token" validate:"required"`
	}
	param := controller.Validated[verifyEmailParams](c, &eventController.constants.Context)
	eventController.eventService.ActivateUser(param.EncodedOrganizerID, param.EncodedEventID, param.Token)

	trans := controller.GetTranslator(c, eventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.organizerActivated")
	controller.Response(c, 200, message, nil)
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
		eventController.awsService.DeleteObject(enums.BannersBucket, objectPath)
		eventController.awsService.UploadObject(enums.BannersBucket, objectPath, param.Banner)
	}

	trans := controller.GetTranslator(c, eventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.updateEvent")
	log.Println("Sending success response")
	controller.Response(c, 200, message, nil)
}

func (eventController *EventController) DeleteEvent(c *gin.Context) {
	type deleteEventParams struct {
		EventID uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[deleteEventParams](c, &eventController.constants.Context)
	event := eventController.eventService.GetEventByID(param.EventID)
	eventController.eventService.DeleteEvent(param.EventID)
	eventController.awsService.DeleteObject(enums.BannersBucket, event.BannerPath)
	listEventMedia := eventController.eventService.GetListEventMedia(param.EventID)
	for _, media := range listEventMedia {
		eventController.awsService.DeleteObject(enums.SessionsBucket, media.Path)
	}

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

func (eventController *EventController) DeleteOrganizer(c *gin.Context) {
	type deleteOrganizerParams struct {
		OrganizerID uint `uri:"organizerID" validate:"required"`
		EventID     uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[deleteOrganizerParams](c, &eventController.constants.Context)
	profilePath := eventController.eventService.GetOrganizerProfilePath(param.OrganizerID)
	eventController.eventService.DeleteOrganizer(param.EventID, param.OrganizerID)
	eventController.awsService.DeleteObject(enums.ProfileBucket, profilePath)

	trans := controller.GetTranslator(c, eventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteOrganizer")
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
	mediaPath := fmt.Sprintf("media/events/%d/resources/%s", param.EventID, param.Media.Filename)
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
	media := eventController.eventService.GetEventMediaDetails(param.MediaID, param.EventID)
	eventController.awsService.DeleteObject(enums.SessionsBucket, media.Path)
	eventController.eventService.DeleteEventMedia(param.MediaID)
	eventController.awsService.DeleteObject(enums.SessionsBucket, media.Path)

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
	for i := range events {
		events[i].Banner = eventController.awsService.GetPresignedURL(enums.BannersBucket, events[i].Banner, 8*time.Hour)
	}
	controller.Response(c, 200, "", events)
}

func (eventController *EventController) GetPublicEvent(c *gin.Context) {
	type getPublicEventParams struct {
		EventID uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[getPublicEventParams](c, &eventController.constants.Context)
	allowedStatus := []enums.EventStatus{enums.Published, enums.Completed}
	eventDetails := eventController.eventService.GetEventDetails(allowedStatus, param.EventID)
	eventDetails.Banner = eventController.awsService.GetPresignedURL(enums.BannersBucket, eventDetails.Banner, 8*time.Hour)
	controller.Response(c, 200, "", eventDetails)
}

func (ec *EventController) ListCategories(c *gin.Context) {
	categoryList := ec.eventService.GetListOfCategories()
	controller.Response(c, 200, "", categoryList)
}

func (eventController *EventController) SearchPublicEvents(c *gin.Context) {
	// some code here ...
}
