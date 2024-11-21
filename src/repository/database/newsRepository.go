package repository_database

import (
	"first-project/src/entities"
	"first-project/src/enums"
	"log"

	"gorm.io/gorm"
)

type NewsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) *NewsRepository {
	return &NewsRepository{
		db: db,
	}
}

func (repo *NewsRepository) CreateNews(title, description, content, imageURL string, category enums.CategoryType, author string) entities.News {
	news := entities.News{
		Title:       title,
		Description: description,
		Content:     content,
		ImageURL:    imageURL,
		Category:    category,
		Author:      author,
	}

	err := repo.db.Create(&news).Error
	if err != nil {
		panic(err)
	}
	return news
}

func (repo *NewsRepository) GetNewsByID(id uint) (entities.News, bool) {
	var news entities.News
	query := repo.db.Where("id = ?", id)
	err := query.First(&news).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return news, false
		}
		panic(err)
	}
	return news, true
}

func (repo *NewsRepository) UpdateNews(newsID uint, title, description, content, imageURL string, category enums.CategoryType, author string) (*entities.News, error) {
	var news entities.News
	query := repo.db.Where("id = ?", newsID).First(&news)
	if query.Error != nil {
		if query.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, query.Error
	}

	news.Title = title
	news.Description = description
	news.Content = content
	news.ImageURL = imageURL
	news.Category = category
	news.Author = author

	err := repo.db.Save(&news).Error
	if err != nil {
		return nil, err
	}
	return &news, nil
}

func (repo *NewsRepository) DeleteNews(id uint) {
	query := repo.db.Where("id = ?", id)
	err := query.Delete(&entities.News{}).Error
	if err != nil {
		panic(err)
	}
}

func (repo *NewsRepository) GetAllNews(categories []enums.CategoryType, limit int, offset int) ([]entities.News, error) {
	var news []entities.News
	query := repo.db

	if len(categories) > 0 {
		query = query.Where("category IN ?", categories)
		log.Printf("Applied category filter: %v", categories)
	}

	query = query.Debug() // Enable debugging to log SQL and params

	err := query.Limit(limit).Offset(offset).Find(&news).Error
	if err != nil {
		panic(err) // should be handled
	}

	log.Printf("Generated SQL Query: %s", query.Statement.SQL.String())
	log.Printf("SQL Vars: %v", query.Statement.Vars)
	log.Printf("Query executed successfully. Retrieved %d news items.", len(news))

	return news, nil
}

func (repo *NewsRepository) GetTopKNews(limit int, categories []enums.CategoryType) ([]entities.News, error) {
	var news []entities.News
	query := repo.db.Order("published_at DESC").Limit(limit)
	if len(categories) > 0 {
		query = query.Where("category IN ?", categories)
	}
	err := query.Find(&news).Error
	if err != nil {
		return nil, err
	}
	return news, nil
}
