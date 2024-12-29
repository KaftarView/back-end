package repository_database_interfaces

import "first-project/src/entities"

type CommentRepository interface {
	CreateNewComment(authorID uint, commentableID uint, content string) *entities.Comment
	CreateNewCommentable() *entities.Commentable
	DeleteCommentContent(comment *entities.Comment)
	FindCommentByID(commentID uint) (*entities.Comment, bool)
	FindCommentableByID(commentableID uint) (*entities.Commentable, bool)
	GetCommentsByEventID(eventID uint) []*entities.Comment
	UpdateCommentContent(comment *entities.Comment, newContent string)
}
