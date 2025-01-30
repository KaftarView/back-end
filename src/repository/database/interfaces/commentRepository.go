package repository_database_interfaces

import (
	"first-project/src/entities"

	"gorm.io/gorm"
)

type CommentRepository interface {
	CreateNewComment(db *gorm.DB, authorID uint, commentableID uint, content string) *entities.Comment
	CreateNewCommentable(db *gorm.DB) *entities.Commentable
	DeleteCommentContent(db *gorm.DB, comment *entities.Comment)
	FindCommentByID(db *gorm.DB, commentID uint) (*entities.Comment, bool)
	FindCommentableByID(db *gorm.DB, commentableID uint) (*entities.Commentable, bool)
	GetCommentsByEventID(db *gorm.DB, eventID uint) []*entities.Comment
	UpdateCommentContent(db *gorm.DB, comment *entities.Comment, newContent string)
}
