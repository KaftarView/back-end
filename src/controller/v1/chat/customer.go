package controller_v1_chat

import (
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/websocket"

	"github.com/gin-gonic/gin"
)

type CustomerChatController struct {
	constants   *bootstrap.Constants
	chatService application_interfaces.ChatService
	hub         *websocket.Hub
}

func NewCustomerChatController(
	constants *bootstrap.Constants,
	chatService application_interfaces.ChatService,
	hub *websocket.Hub,
) *CustomerChatController {
	return &CustomerChatController{
		constants:   constants,
		chatService: chatService,
		hub:         hub,
	}
}

func (customerChatController *CustomerChatController) CreateOrGetRoom(c *gin.Context) {
	userID, _ := c.Get(customerChatController.constants.Context.UserID)
	roomsDetails := customerChatController.chatService.CreateOrGetRoom(userID.(uint))

	controller.Response(c, 200, "", roomsDetails)
}

func (customerChatController *CustomerChatController) GetMessages(c *gin.Context) {
	type getMessagesParams struct {
		RoomID uint `uri:"roomID" validate:"required"`
	}
	param := controller.Validated[getMessagesParams](c, &customerChatController.constants.Context)
	messages := customerChatController.chatService.GetRoomMessages(param.RoomID)

	controller.Response(c, 200, "", messages)
}

func (customerChatController *CustomerChatController) HandleWebsocket(c *gin.Context) {
	type roomConnectionParams struct {
		RoomID uint `uri:"roomID" validate:"required"`
	}
	param := controller.Validated[roomConnectionParams](c, &customerChatController.constants.Context)
	userID, _ := c.Get(customerChatController.constants.Context.UserID)
	conn, _ := c.Get(customerChatController.constants.Context.WebsocketConnection)

	client := websocket.NewClient(customerChatController.hub, conn, param.RoomID, userID.(uint), customerChatController.chatService)

	client.Hub.Register <- client

	go client.ReadPump()
	go client.WritePump()
}
