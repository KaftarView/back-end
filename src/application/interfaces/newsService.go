package application_interfaces

import (
	"first-project/src/dto"
	"first-project/src/entities"
)

type NewsService interface {
	CreateNews(newsDetails dto.CreateNewsRequest) *entities.News
	DeleteNews(newsID uint)
	FilterNewsByCategory(categories []string, page int, pageSize int) []dto.NewsDetailsResponse
	GetNewsDetails(newsID uint) dto.NewsDetailsResponse
	GetNewsList(page int, pageSize int) []dto.NewsDetailsResponse
	SearchNews(query string, page int, pageSize int) []dto.NewsDetailsResponse
	UpdateNews(newsDetails dto.UpdateNewsRequest)
}
