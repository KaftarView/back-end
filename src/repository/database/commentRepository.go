package repository_database

import (
	"first-project/src/entities"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (repo *CommentRepository) CreateNewCommentable() entities.Commentable {
	commentable := entities.Commentable{}
	result := repo.db.Create(&commentable)
	if result.Error != nil {
		panic(result.Error)
	}
	return commentable
}

func (repo *CommentRepository) FindCommentableByID(commentableID uint) (entities.Commentable, bool) {
	var commentable entities.Commentable
	result := repo.db.First(&commentable, "c_id = ?", commentableID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return commentable, false
		}
		panic(result.Error)
	}
	return commentable, true
}

func (repo *CommentRepository) FindCommentByID(commentID uint) (entities.Comment, bool) {
	var comment entities.Comment
	result := repo.db.First(&comment, "id = ?", commentID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return comment, false
		}
		panic(result.Error)
	}
	return comment, true

}

func (repo *CommentRepository) UpdateCommentContent(comment entities.Comment, newContent string) {
	comment.Content = newContent
	comment.IsModerated = true
	repo.db.Save(&comment)
}

func (repo *CommentRepository) CreateNewComment(authorID, commentableID uint, content string) entities.Comment {
	comment := entities.Comment{
		AuthorID:      authorID,
		IsModerated:   false,
		Content:       content,
		CommentableID: commentableID,
	}
	result := repo.db.Create(&comment)
	if result.Error != nil {
		panic(result.Error)
	}
	return comment
}
