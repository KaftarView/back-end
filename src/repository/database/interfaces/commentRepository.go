package repository_database_interfaces

import (
	"first-project/src/entities"

	"gorm.io/gorm"
)

type CommentRepository interface {
	CreateNewComment(authorID uint, commentableID uint, content string) *entities.Comment
	CreateNewCommentable(db *gorm.DB) *entities.Commentable
	DeleteCommentContent(comment *entities.Comment)
	FindCommentByID(commentID uint) (*entities.Comment, bool)
	FindCommentableByID(commentableID uint) (*entities.Commentable, bool)
	GetCommentsByEventID(eventID uint) []*entities.Comment
	UpdateCommentContent(comment *entities.Comment, newContent string)
}
