package repository_database_interfaces

import (
	"first-project/src/entities"

	"gorm.io/gorm"
)

type ChatRepository interface {
	GetMessagesByRoomID(db *gorm.DB, roomID uint) []*entities.ChatMessage
	GetRoomByID(db *gorm.DB, roomID uint) (*entities.ChatRoom, bool)
	GetRoomByUserID(db *gorm.DB, userID uint) []*entities.ChatRoom
	CreateRoom(db *gorm.DB, room *entities.ChatRoom) error
	SaveMessage(db *gorm.DB, message *entities.ChatMessage) error
}
