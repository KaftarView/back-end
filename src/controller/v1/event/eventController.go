package controller_v1_event

import "github.com/gin-gonic/gin"

type EventController struct {
}

func NewEventController() *EventController {
	return &EventController{}
}

func (eventController *EventController) ListEvents(c *gin.Context) {
	// some code here ...
}

func (eventController *EventController) GetEvent(c *gin.Context) {
	// some code here ...
}

func (eventController *EventController) CreateEvent(c *gin.Context) {
	// some code here ...
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
