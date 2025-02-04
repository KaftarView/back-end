package repository_database

import (
	"first-project/src/entities"

	"gorm.io/gorm"
)

type CommentRepository struct{}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{}
}

func (repo *CommentRepository) GetCommentsByEventID(db *gorm.DB, eventID uint) []*entities.Comment {
	var comments []*entities.Comment

	result := db.Where("commentable_id = ?", eventID).Preload("Author").Order("created_at DESC").Find(&comments)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return comments
}

func (repo *CommentRepository) CreateNewCommentable(db *gorm.DB) *entities.Commentable {
	commentable := entities.Commentable{}
	result := db.Create(&commentable)
	if result.Error != nil {
		panic(result.Error)
	}
	return &commentable
}

func (repo *CommentRepository) FindCommentableByID(db *gorm.DB, commentableID uint) (*entities.Commentable, bool) {
	var commentable entities.Commentable
	result := db.First(&commentable, "c_id = ?", commentableID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &commentable, true
}

func (repo *CommentRepository) FindCommentByID(db *gorm.DB, commentID uint) (*entities.Comment, bool) {
	var comment entities.Comment
	result := db.First(&comment, "id = ?", commentID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &comment, true

}

func (repo *CommentRepository) UpdateCommentContent(db *gorm.DB, comment *entities.Comment, newContent string) {
	comment.Content = newContent
	comment.IsModerated = true
	db.Save(comment)
}

func (repo *CommentRepository) DeleteCommentContent(db *gorm.DB, comment *entities.Comment) {
	err := db.Unscoped().Delete(comment).Error
	if err != nil {
		panic(err)
	}

}

func (repo *CommentRepository) CreateNewComment(db *gorm.DB, authorID, commentableID uint, content string) *entities.Comment {
	comment := &entities.Comment{
		AuthorID:      authorID,
		IsModerated:   false,
		Content:       content,
		CommentableID: commentableID,
	}
	result := db.Create(comment)
	if result.Error != nil {
		panic(result.Error)
	}
	return comment
}
