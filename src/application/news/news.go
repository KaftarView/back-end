package application_news

import (
	"errors"
	"first-project/src/entities"
	"first-project/src/enums"
	repository_database "first-project/src/repository/database"

	"gorm.io/gorm"
)

type NewsService struct {
	newsRepo *repository_database.NewsRepository
}

func NewNewsService(newsRepo *repository_database.NewsRepository) *NewsService {
	return &NewsService{newsRepo: newsRepo}
}

func (ns *NewsService) CreateNews(news entities.News) (*entities.News, error) {
	createdNews, err := ns.newsRepo.CreateNews(news)
	if err != nil {
		return nil, err
	}
	return &createdNews, nil
}

func (ns *NewsService) GetNewsByID(newsID uint) (*entities.News, error) {
	news, err := ns.newsRepo.GetNewsByID(newsID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("news not found")
		}
		return nil, err
	}
	return &news, nil
}

func (ns *NewsService) UpdateNews(newsID uint, updatedNews entities.News) (*entities.News, error) {
	existingNews, err := ns.newsRepo.GetNewsByID(newsID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("news not found")
		}
		return nil, err
	}

	existingNews.Title = updatedNews.Title
	existingNews.Description = updatedNews.Description
	existingNews.Content = updatedNews.Content
	existingNews.ImageURL = updatedNews.ImageURL
	existingNews.Category = updatedNews.Category
	existingNews.PublishedAt = updatedNews.PublishedAt

	updated, err := ns.newsRepo.UpdateNews(existingNews)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (ns *NewsService) DeleteNews(newsID uint) error {
	_, err := ns.newsRepo.GetNewsByID(newsID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("news not found")
		}
		return err
	}
	return ns.newsRepo.DeleteNews(newsID)
}

func (ns *NewsService) GetAllNews(categories []enums.CategoryType, limit int, offset int) ([]entities.News, error) {
	news, err := ns.newsRepo.GetAllNews(categories, limit, offset)
	if err != nil {
		return nil, err
	}
	return news, nil
}

func (ns *NewsService) GetTopKNews(limit int, categories []enums.CategoryType) ([]entities.News, error) {
	news, err := ns.newsRepo.GetTopKNews(limit, categories)
	if err != nil {
		return nil, err
	}
	return news, nil
}
