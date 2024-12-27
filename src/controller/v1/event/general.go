package controller_v1_event

import (
	"first-project/src/application"
	application_communication "first-project/src/application/communication/emailService"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/enums"

	"github.com/gin-gonic/gin"
)

type GeneralEventController struct {
	constants    *bootstrap.Constants
	eventService *application.EventService
	emailService *application_communication.EmailService
}

func NewGeneralEventController(
	constants *bootstrap.Constants,
	eventService *application.EventService,
	emailService *application_communication.EmailService,
) *GeneralEventController {
	return &GeneralEventController{
		constants:    constants,
		eventService: eventService,
		emailService: emailService,
	}
}

func (generalEventController *GeneralEventController) ListPublicEvents(c *gin.Context) {
	type publicEventsListParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	param := controller.Validated[publicEventsListParams](c, &generalEventController.constants.Context)
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

func (generalEventController *GeneralEventController) SearchPublicEvents(c *gin.Context) {
	type searchPublicEventParams struct {
		Query    string `form:"query"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
	}
	param := controller.Validated[searchPublicEventParams](c, &generalEventController.constants.Context)
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

func (generalEventController *GeneralEventController) FilterPublicEvents(c *gin.Context) {
	type filterPublicEventParams struct {
		Categories []string `form:"categories"`
		Page       int      `form:"page"`
		PageSize   int      `form:"pageSize"`
	}
	param := controller.Validated[filterPublicEventParams](c, &generalEventController.constants.Context)
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

func (generalEventController *GeneralEventController) GetPublicEventDetails(c *gin.Context) {
	type getPublicEventParams struct {
		EventID uint `uri:"eventID" validate:"required"`
	}
	param := controller.Validated[getPublicEventParams](c, &generalEventController.constants.Context)
	allowedStatus := []enums.EventStatus{enums.Published, enums.Completed}
	eventDetails := generalEventController.eventService.GetEventDetails(allowedStatus, param.EventID)

	controller.Response(c, 200, "", eventDetails)
}
