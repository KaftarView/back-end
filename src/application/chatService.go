package application

import (
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"
	"first-project/src/dto"
	"first-project/src/entities"
	"first-project/src/enums"
	"first-project/src/exceptions"
	repository_database_interfaces "first-project/src/repository/database/interfaces"

	"gorm.io/gorm"
)

type ChatService struct {
	constants      *bootstrap.Constants
	userService    application_interfaces.UserService
	chatRepository repository_database_interfaces.ChatRepository
	db             *gorm.DB
}

func NewChatService(
	constants *bootstrap.Constants,
	userService application_interfaces.UserService,
	chatRepository repository_database_interfaces.ChatRepository,
	db *gorm.DB,
) *ChatService {
	return &ChatService{
		constants:      constants,
		userService:    userService,
		chatRepository: chatRepository,
		db:             db,
	}
}

func getRoomDetails(room *entities.ChatRoom) dto.RoomDetailsResponse {
	admins := make([]dto.UserDetailsResponse, len(room.Admins))
	for i, admin := range room.Admins {
		admins[i] = dto.UserDetailsResponse{
			Name:  admin.Name,
			Email: admin.Email,
		}
	}
	return dto.RoomDetailsResponse{
		ID:     room.ID,
		Tag:    room.Tag.String(),
		Admins: admins,
		Member: dto.UserDetailsResponse{
			Name:  room.Member.Name,
			Email: room.Member.Email,
		},
	}
}

func (chatService *ChatService) CreateOrGetRoom(userID uint) []dto.RoomDetailsResponse {
	var roomsDetails []dto.RoomDetailsResponse
	rooms := chatService.chatRepository.GetRoomByUserID(chatService.db, userID)

	if len(rooms) == 0 {
		supportAdmins := chatService.userService.GetUsersByPermissions([]enums.PermissionType{enums.CustomerSupport, enums.All})
		room := &entities.ChatRoom{
			Tag:      enums.Support,
			MemberID: userID,
			Admins:   supportAdmins,
		}
		chatService.chatRepository.CreateRoom(chatService.db, room)
		roomsDetails = append(roomsDetails, getRoomDetails(room))
		return roomsDetails
	}
	for _, room := range rooms {
		roomsDetails = append(roomsDetails, getRoomDetails(room))
	}
	return roomsDetails
}

func (chatService *ChatService) GetRoomMessages(roomID uint) []dto.MessageDetailsResponse {
	var notFoundError exceptions.NotFoundError
	_, roomExist := chatService.chatRepository.GetRoomByID(chatService.db, roomID)
	if !roomExist {
		notFoundError.ErrorField = chatService.constants.ErrorField.Room
		panic(notFoundError)
	}
	messages := chatService.chatRepository.GetMessagesByRoomID(chatService.db, roomID)
	messagesDetails := make([]dto.MessageDetailsResponse, len(messages))
	for i, message := range messages {
		sender := dto.UserDetailsResponse{
			Name:  message.Sender.Name,
			Email: message.Sender.Email,
		}
		messagesDetails[i] = dto.MessageDetailsResponse{
			Sender:  sender,
			Content: message.Content,
		}
	}
	return messagesDetails
}

func (chatService *ChatService) SaveMessage(roomID, senderID uint, content string) {
	var notFoundError exceptions.NotFoundError
	room, roomExist := chatService.chatRepository.GetRoomByID(chatService.db, roomID)
	if !roomExist {
		notFoundError.ErrorField = chatService.constants.ErrorField.Room
		panic(notFoundError)
	}

	isValidSender := false
	if senderID == room.MemberID {
		isValidSender = true
	} else {
		for _, admin := range room.Admins {
			if senderID == admin.ID {
				isValidSender = true
				break
			}
		}
	}

	if !isValidSender {
		panic(exceptions.NewForbiddenError())
	}

	message := &entities.ChatMessage{
		RoomID:   roomID,
		SenderID: senderID,
		Content:  content,
	}
	if err := chatService.chatRepository.SaveMessage(chatService.db, message); err != nil {
		panic(err)
	}
}
