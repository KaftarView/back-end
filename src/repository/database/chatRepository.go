package repository_database

type chatRepository struct{}

func NewChatRepository() *chatRepository {
	return &chatRepository{}
}
