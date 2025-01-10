package controller_v1_chat

import (
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"

	"github.com/gin-gonic/gin"
)

type CustomerChatController struct {
	constants   *bootstrap.Constants
	chatService application_interfaces.ChatService
}

func NewCustomerCommentController(
	constants *bootstrap.Constants,
	chatService application_interfaces.ChatService,
) *CustomerChatController {
	return &CustomerChatController{
		constants:   constants,
		chatService: chatService,
	}
}

func (customerChatController *CustomerChatController) CreateOrGetRoom(c *gin.Context) {
	// some codes here ...
}

func (customerChatController *CustomerChatController) HandleWebsocket(c *gin.Context) {
	// some codes here ...
}
