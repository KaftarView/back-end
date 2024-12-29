package application

import (
	"first-project/src/bootstrap"
	"first-project/src/dto"
	"first-project/src/exceptions"
	repository_database "first-project/src/repository/database"
)

type commentService struct {
	constants         *bootstrap.Constants
	commentRepository *repository_database.CommentRepository
	userRepository    *repository_database.UserRepository
}

func NewCommentService(
	constants *bootstrap.Constants,
	commentRepository *repository_database.CommentRepository,
	userRepository *repository_database.UserRepository,
) *commentService {
	return &commentService{
		constants:         constants,
		commentRepository: commentRepository,
		userRepository:    userRepository,
	}
}

func (commentService *commentService) GetPostComments(commentableID uint) []dto.CommentDetailsResponse {
	comments := commentService.commentRepository.GetCommentsByEventID(commentableID)
	var commentsDetails []dto.CommentDetailsResponse
	for _, comment := range comments {
		commentsDetails = append(commentsDetails, dto.CommentDetailsResponse{
			ID:          comment.ID,
			Content:     comment.Content,
			IsModerated: comment.IsModerated,
			AuthorName:  comment.Author.Name,
		})
	}
	return commentsDetails
}

func (commentService *commentService) CreateComment(authorID, commentableID uint, content string) {
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

func (commentService *commentService) EditComment(authorID, commentID uint, newContent string) {
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

func (commentService *commentService) DeleteComment(authorID, commentID uint, canModerateComment bool) {
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
