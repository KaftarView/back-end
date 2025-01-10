package application

import (
	"first-project/src/bootstrap"
	repository_database_interfaces "first-project/src/repository/database/interfaces"

	"gorm.io/gorm"
)

type chatService struct {
	constants      *bootstrap.Constants
	chatRepository repository_database_interfaces.ChatRepository
	db             *gorm.DB
}

func NewChatService(
	constants *bootstrap.Constants,
	chatRepository repository_database_interfaces.ChatRepository,
	db *gorm.DB,
) *chatService {
	return &chatService{
		constants:      constants,
		chatRepository: chatRepository,
		db:             db,
	}
}
