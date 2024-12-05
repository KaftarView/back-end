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

func (commentService *CommentService) EditComment(authorID, commentID uint, newContent string) {
	var notFoundError exceptions.NotFoundError
	_, authorExist := commentService.userRepository.FindByUserID(authorID)
	if !authorExist {
		notFoundError.ErrorField = commentService.constants.ErrorField.User
		panic(notFoundError)
	}
	comment, commentExist := commentService.commentRepository.FindCommentByID(commentID)
	if !commentExist {
		notFoundError.ErrorField = commentService.constants.ErrorField.Comment
		panic(notFoundError)
	}
	if comment.AuthorID != authorID {
		authError := exceptions.NewForbiddenError()
		panic(authError)
	}
	commentService.commentRepository.UpdateCommentContent(comment, newContent)
}

func (commentService *CommentService) DeleteComment(authorID, commentID uint, canModerateComment bool) {
	var notFoundError exceptions.NotFoundError
	_, authorExist := commentService.userRepository.FindByUserID(authorID)
	if !authorExist {
		notFoundError.ErrorField = commentService.constants.ErrorField.User
		panic(notFoundError)
	}
	comment, commentExist := commentService.commentRepository.FindCommentByID(commentID)
	if !commentExist {
		notFoundError.ErrorField = commentService.constants.ErrorField.Comment
		panic(notFoundError)
	}
	if !canModerateComment && comment.AuthorID != authorID {
		authError := exceptions.NewForbiddenError()
		panic(authError)
	}
	commentService.commentRepository.DeleteCommentContent(comment)
}
