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
	type eventsListParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	param := controller.Validated[eventsListParams](c, &generalEventController.constants.Context)
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	allowedStatus := []enums.EventStatus{enums.Published}
	events := generalEventController.eventService.GetEventsList(allowedStatus, param.Page, param.PageSize)

	controller.Response(c, 200, "", events)
}

func (generalEventController *GeneralEventController) SearchEvents(c *gin.Context) {
	type searchEventParams struct {
		Query    string `form:"query"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
	}
	param := controller.Validated[searchEventParams](c, &generalEventController.constants.Context)
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	allowedStatus := []enums.EventStatus{enums.Published, enums.Completed}
	events := generalEventController.eventService.SearchEvents(param.Query, param.Page, param.PageSize, allowedStatus)

	controller.Response(c, 200, "", events)
}

func (generalEventController *GeneralEventController) FilterEvents(c *gin.Context) {
	type filterEventParams struct {
		Categories []string `form:"categories"`
		Page       int      `form:"page"`
		PageSize   int      `form:"pageSize"`
	}
	param := controller.Validated[filterEventParams](c, &generalEventController.constants.Context)
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	allowedStatus := []enums.EventStatus{enums.Published, enums.Completed}
	events := generalEventController.eventService.FilterEventsByCategories(param.Categories, param.Page, param.PageSize, allowedStatus)

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
