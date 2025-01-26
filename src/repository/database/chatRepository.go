package repository_database

import (
	"first-project/src/entities"

	"gorm.io/gorm"
)

type chatRepository struct{}

func NewChatRepository() *chatRepository {
	return &chatRepository{}
}

func (repo *chatRepository) GetRoomByID(db *gorm.DB, roomID uint) (*entities.ChatRoom, bool) {
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

func (repo *chatRepository) GetRoomByUserID(db *gorm.DB, userID uint) []*entities.ChatRoom {
	var rooms []*entities.ChatRoom
	result := db.Joins("LEFT JOIN chat_room_admins ON chat_room_admins.chat_room_id = chat_rooms.id").
		Where("chat_rooms.member_id = ? OR chat_room_admins.user_id = ?", userID, userID).
		Preload("Admins").
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

func (repo *chatRepository) CreateRoom(db *gorm.DB, room *entities.ChatRoom) error {
	return db.Create(room).Error
}

func (repo *chatRepository) GetMessagesByRoomID(db *gorm.DB, roomID uint) []*entities.ChatMessage {
	var messages []*entities.ChatMessage

	result := db.Where("room_id = ?", roomID).
		Order("created_at desc").
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

func (repo *chatRepository) SaveMessage(db *gorm.DB, message *entities.ChatMessage) error {
	return db.Create(message).Error
}
