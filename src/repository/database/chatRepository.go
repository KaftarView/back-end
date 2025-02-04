package repository_database

import (
	"first-project/src/entities"

	"gorm.io/gorm"
)

type ChatRepository struct{}

func NewChatRepository() *ChatRepository {
	return &ChatRepository{}
}

func (repo *ChatRepository) GetRoomByID(db *gorm.DB, roomID uint) (*entities.ChatRoom, bool) {
	var room entities.ChatRoom
	result := db.Preload("Admins").First(&room, "id = ?", roomID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &room, true
}

func (repo *ChatRepository) GetRoomByUserID(db *gorm.DB, userID uint) []*entities.ChatRoom {
	var rooms []*entities.ChatRoom
	result := db.Joins("LEFT JOIN chat_room_admins ON chat_room_admins.chat_room_id = chat_rooms.id").
		Where("chat_rooms.member_id = ? OR chat_room_admins.user_id = ?", userID, userID).
		Preload("Admins").
		Preload("Member").
		Distinct().
		Find(&rooms)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return rooms
}

func (repo *ChatRepository) CreateRoom(db *gorm.DB, room *entities.ChatRoom) error {
	return db.Create(room).Error
}

func (repo *ChatRepository) GetMessagesByRoomID(db *gorm.DB, roomID uint) []*entities.ChatMessage {
	var messages []*entities.ChatMessage

	result := db.Where("room_id = ?", roomID).
		Order("created_at").
		Preload("Sender").
		Find(&messages)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return messages
}

func (repo *ChatRepository) SaveMessage(db *gorm.DB, message *entities.ChatMessage) error {
	return db.Create(message).Error
}
