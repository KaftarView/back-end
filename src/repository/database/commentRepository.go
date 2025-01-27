package repository_database

import (
	"first-project/src/entities"

	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *commentRepository {
	return &commentRepository{
		db: db,
	}
}

func (repo *commentRepository) GetCommentsByEventID(eventID uint) []*entities.Comment {
	var comments []*entities.Comment

	result := repo.db.Where("commentable_id = ?", eventID).Preload("Author").Order("created_at DESC").Find(&comments)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return comments
}

func (repo *commentRepository) CreateNewCommentable(db *gorm.DB) *entities.Commentable {
	commentable := entities.Commentable{}
	result := db.Create(&commentable)
	if result.Error != nil {
		panic(result.Error)
	}
	return &commentable
}

func (repo *commentRepository) FindCommentableByID(commentableID uint) (*entities.Commentable, bool) {
	var commentable entities.Commentable
	result := repo.db.First(&commentable, "c_id = ?", commentableID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &commentable, true
}

func (repo *commentRepository) FindCommentByID(commentID uint) (*entities.Comment, bool) {
	var comment entities.Comment
	result := repo.db.First(&comment, "id = ?", commentID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &comment, true

}

func (repo *commentRepository) UpdateCommentContent(comment *entities.Comment, newContent string) {
	comment.Content = newContent
	comment.IsModerated = true
	repo.db.Save(comment)
}

func (repo *commentRepository) DeleteCommentContent(comment *entities.Comment) {
	err := repo.db.Unscoped().Delete(comment).Error
	if err != nil {
		panic(err)
	}

}

func (repo *commentRepository) CreateNewComment(authorID, commentableID uint, content string) *entities.Comment {
	comment := &entities.Comment{
		AuthorID:      authorID,
		IsModerated:   false,
		Content:       content,
		CommentableID: commentableID,
	}
	result := repo.db.Create(comment)
	if result.Error != nil {
		panic(result.Error)
	}
	return comment
}
