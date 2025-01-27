package controller_v1_event

import (
	application_communication "first-project/src/application/communication/emailService"
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/enums"

	"github.com/gin-gonic/gin"
)

type GeneralEventController struct {
	constants    *bootstrap.Constants
	eventService application_interfaces.EventService
	emailService *application_communication.EmailService
}

func NewGeneralEventController(
	constants *bootstrap.Constants,
	eventService application_interfaces.EventService,
	emailService *application_communication.EmailService,
) *GeneralEventController {
	return &GeneralEventController{
		constants:    constants,
		eventService: eventService,
		emailService: emailService,
	}
}

func (generalEventController *GeneralEventController) ListEvents(c *gin.Context) {
	pagination := controller.GetPagination(c, &generalEventController.constants.Context)
	allowedStatus := []enums.EventStatus{enums.Published}
	events := generalEventController.eventService.GetEventsList(allowedStatus, pagination.Page, pagination.PageSize)

	controller.Response(c, 200, "", events)
}

func (generalEventController *GeneralEventController) SearchEvents(c *gin.Context) {
	type searchEventParams struct {
		Query string `form:"query"`
	}
	param := controller.Validated[searchEventParams](c, &generalEventController.constants.Context)
	pagination := controller.GetPagination(c, &generalEventController.constants.Context)
	allowedStatus := []enums.EventStatus{enums.Published, enums.Completed}
	events := generalEventController.eventService.SearchEvents(param.Query, pagination.Page, pagination.PageSize, allowedStatus)

	controller.Response(c, 200, "", events)
}

func (generalEventController *GeneralEventController) FilterEvents(c *gin.Context) {
	type filterEventParams struct {
		Categories []string `form:"categories"`
	}
	param := controller.Validated[filterEventParams](c, &generalEventController.constants.Context)
	pagination := controller.GetPagination(c, &generalEventController.constants.Context)
	allowedStatus := []enums.EventStatus{enums.Published, enums.Completed}
	events := generalEventController.eventService.FilterEventsByCategories(param.Categories, pagination.Page, pagination.PageSize, allowedStatus)

	controller.Response(c, 200, "", events)
}

func (generalEventController *GeneralEventController) GetEventDetails(c *gin.Context) {
	type getEventDetailsParams struct {
		EventID uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[getEventDetailsParams](c, &generalEventController.constants.Context)
	allowedStatus := []enums.EventStatus{enums.Published, enums.Completed}
	eventDetails := generalEventController.eventService.GetEventDetails(allowedStatus, param.EventID)

	controller.Response(c, 200, "", eventDetails)
}

func (generalEventController *GeneralEventController) GetEventOrganizers(c *gin.Context) {
	type getEventOrganizersParams struct {
		EventID uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[getEventOrganizersParams](c, &generalEventController.constants.Context)
	organizersDetails := generalEventController.eventService.GetEventOrganizers(param.EventID)

	controller.Response(c, 200, "", organizersDetails)
}
