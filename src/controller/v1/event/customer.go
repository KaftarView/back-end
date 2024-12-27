package controller_v1_event

import (
	"first-project/src/application"
	application_communication "first-project/src/application/communication/emailService"
	"first-project/src/bootstrap"
	"first-project/src/controller"

	"github.com/gin-gonic/gin"
)

type CustomerEventController struct {
	constants    *bootstrap.Constants
	eventService *application.EventService
	emailService *application_communication.EmailService
}

func NewCustomerEventController(
	constants *bootstrap.Constants,
	eventService *application.EventService,
	emailService *application_communication.EmailService,
) *CustomerEventController {
	return &CustomerEventController{
		constants:    constants,
		eventService: eventService,
		emailService: emailService,
	}
}

func (customerEventController *CustomerEventController) GetAvailableEventTicketsList(c *gin.Context) {
	type getEventParams struct {
		EventID uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[getEventParams](c, &customerEventController.constants.Context)
	ticketDetails := customerEventController.eventService.GetEventTickets(param.EventID, []bool{true})
	controller.Response(c, 200, "", ticketDetails)
}
