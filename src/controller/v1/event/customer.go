package controller_v1_event

import (
	application_communication "first-project/src/application/communication/emailService"
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/dto"

	"github.com/gin-gonic/gin"
)

type CustomerEventController struct {
	constants    *bootstrap.Constants
	eventService application_interfaces.EventService
	emailService *application_communication.EmailService
}

func NewCustomerEventController(
	constants *bootstrap.Constants,
	eventService application_interfaces.EventService,
	emailService *application_communication.EmailService,
) *CustomerEventController {
	return &CustomerEventController{
		constants:    constants,
		eventService: eventService,
		emailService: emailService,
	}
}

func (customerEventController *CustomerEventController) GetAllUserJoinedEvents(c *gin.Context) {
	userID, _ := c.Get(customerEventController.constants.Context.UserID)
	events := customerEventController.eventService.GetAllUserJoinedEvents(userID.(uint))

	controller.Response(c, 200, "", events)
}

func (customerEventController *CustomerEventController) GetAvailableEventTicketsList(c *gin.Context) {
	type getEventParams struct {
		EventID uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[getEventParams](c, &customerEventController.constants.Context)
	ticketDetails := customerEventController.eventService.GetAvailableEventTickets(param.EventID)
	controller.Response(c, 200, "", ticketDetails)
}

func (customerEventController *CustomerEventController) ReserveTickets(c *gin.Context) {
	type ticketParams struct {
		TicketID uint `json:"ticketID" validate:"required"`
		Quantity uint `json:"quantity" validate:"required"`
	}
	type reserveTicketsParams struct {
		Tickets      []ticketParams `json:"tickets" validate:"required"`
		DiscountCode *string        `json:"discountCode"`
		EventID      uint           `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[reserveTicketsParams](c, &customerEventController.constants.Context)

	userID, _ := c.Get(customerEventController.constants.Context.UserID)
	tickets := make([]dto.ReserveTicketRequest, len(param.Tickets))
	for i, ticket := range param.Tickets {
		tickets[i] = dto.ReserveTicketRequest{
			ID:       ticket.TicketID,
			Quantity: ticket.Quantity,
		}
	}
	totalPrice := customerEventController.eventService.ReserveEventTicket(userID.(uint), param.EventID, param.DiscountCode, tickets)

	trans := controller.GetTranslator(c, customerEventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.reserveTicket")
	controller.Response(c, 200, message, totalPrice)
}

func (customerEventController *CustomerEventController) PurchaseTickets(c *gin.Context) {
	type purchaseTicketsParams struct {
		ReservationID uint `uri:"reservationID" validate:"required"`
		EventID       uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[purchaseTicketsParams](c, &customerEventController.constants.Context)
	userID, _ := c.Get(customerEventController.constants.Context.UserID)
	customerEventController.eventService.PurchaseEventTicket(userID.(uint), param.EventID, param.ReservationID)

	trans := controller.GetTranslator(c, customerEventController.constants.Context.Translator)
	message, _ := trans.T("successMessage.purchaseTicket")
	controller.Response(c, 200, message, nil)
}
