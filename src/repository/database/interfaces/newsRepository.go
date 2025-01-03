package repository_database_interfaces

import (
	"first-project/src/entities"

	"gorm.io/gorm"
)

type NewsRepository interface {
	CreateNews(db *gorm.DB, news *entities.News) error
	DeleteNews(db *gorm.DB, newsID uint) error
	FindAllNews(db *gorm.DB, offset int, pageSize int) ([]*entities.News, bool)
	FindNewsByCategoryName(db *gorm.DB, categories []string, offset int, pageSize int) []*entities.News
	FindNewsByID(db *gorm.DB, newsID uint) (*entities.News, bool)
	FindNewsByTitle(db *gorm.DB, name string) (*entities.News, bool)
	FindNewsCategoriesByNews(db *gorm.DB, news *entities.News) []entities.Category
	FullTextSearch(db *gorm.DB, query string, offset int, pageSize int) []*entities.News
	GetNewsByID(db *gorm.DB, id uint) (*entities.News, bool)
	UpdateNews(db *gorm.DB, news *entities.News) error
	UpdateNewsCategories(db *gorm.DB, newsID uint, categories []entities.Category) error
}
