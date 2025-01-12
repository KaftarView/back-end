package application_interfaces

import (
	"first-project/src/dto"
)

type ChatService interface {
	CreateOrGetRoom(userID uint) []dto.RoomDetailsResponse
	SaveMessage(roomID, senderID uint, content string)
	GetRoomMessages(roomID uint) []dto.MessageDetailsResponse
}
