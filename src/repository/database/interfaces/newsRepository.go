package repository_database_interfaces

import "first-project/src/entities"

type NewsRepository interface {
	CreateNews(news *entities.News) *entities.News
	DeleteNews(newsID uint)
	FindAllNews(offset int, pageSize int) ([]*entities.News, bool)
	FindNewsByCategoryName(categories []string, offset int, pageSize int) []*entities.News
	FindNewsByID(newsID uint) (*entities.News, bool)
	FindNewsByTitle(name string) (*entities.News, bool)
	FindNewsCategoriesByNews(news *entities.News) []entities.Category
	FullTextSearch(query string, offset int, pageSize int) []*entities.News
	GetNewsByID(id uint) (*entities.News, bool)
	UpdateNews(news *entities.News)
	UpdateNewsCategories(newsID uint, categories []entities.Category)
}
