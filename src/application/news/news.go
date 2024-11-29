package application_news

import (
	"first-project/src/entities"
	"first-project/src/enums"
	repository_database "first-project/src/repository/database"
)

type NewsService struct {
	newsRepo *repository_database.NewsRepository
}

func NewNewsService(newsRepo *repository_database.NewsRepository) *NewsService {
	return &NewsService{newsRepo: newsRepo}
}

func (ns *NewsService) CreateNews(title, description, content, category, author string) {
	categoryType, err := enums.StringToCategoryType(category)
	if err != nil {
		panic(err)
	}

	ns.newsRepo.CreateNews(title, description, content, categoryType, author)
}

func (ns *NewsService) GetNewsByID(newsID uint) (*entities.News, bool) {
	news, found := ns.newsRepo.GetNewsByID(newsID)
	return &news, found
}

func (ns *NewsService) UpdateNews(newsID uint, title, description, content, category, author string) (*entities.News, bool) {
	categoryType, err := enums.StringToCategoryType(category)
	if err != nil {
		panic(err)
	}

	updated, err := ns.newsRepo.UpdateNews(newsID, title, description, content, categoryType, author)
	if err != nil {
		panic(err)
	}
	return updated, true
}

func (ns *NewsService) DeleteNews(newsID uint) bool {
	_, found := ns.newsRepo.GetNewsByID(newsID)
	if !found {
		return false
	}
	ns.newsRepo.DeleteNews(newsID)
	return true
}

func (ns *NewsService) GetAllNews(categories []enums.CategoryType, limit int, offset int) []entities.News {
	news, err := ns.newsRepo.GetAllNews(categories, limit, offset)
	if err != nil {
		panic(err)
	}
	return news
}

func (ns *NewsService) GetTopKNews(limit int, categories []enums.CategoryType) ([]entities.News, error) {
	news, err := ns.newsRepo.GetTopKNews(limit, categories)
	if err != nil {
		return nil, err
	}
	return news, nil
}
