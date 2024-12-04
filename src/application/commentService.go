package application

import (
	"first-project/src/bootstrap"
	"first-project/src/exceptions"
	repository_database "first-project/src/repository/database"
)

type CommentService struct {
	constants         *bootstrap.Constants
	commentRepository *repository_database.CommentRepository
	userRepository    *repository_database.UserRepository
}

func NewCommentService(
	constants *bootstrap.Constants,
	commentRepository *repository_database.CommentRepository,
	userRepository *repository_database.UserRepository,
) *CommentService {
	return &CommentService{
		constants:         constants,
		commentRepository: commentRepository,
		userRepository:    userRepository,
	}
}

func (commentService *CommentService) CreateComment(authorID, commentableID uint, content string) {
	var notFoundError exceptions.NotFoundError
	_, authorExist := commentService.userRepository.FindByUserID(authorID)
	if !authorExist {
		notFoundError.ErrorField = commentService.constants.ErrorField.User
		panic(notFoundError)
	}
	_, postExist := commentService.commentRepository.FindCommentableByID(commentableID)
	if !postExist {
		notFoundError.ErrorField = commentService.constants.ErrorField.Post
		panic(notFoundError)
	}
	commentService.commentRepository.CreateNewComment(authorID, commentableID, content)
}
