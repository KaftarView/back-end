package application_news

import (
	"first-project/src/entities"
	"first-project/src/enums"
	repository_database "first-project/src/repository/database"
	"log"
)

type NewsService struct {
	newsRepo *repository_database.NewsRepository
}

func NewNewsService(newsRepo *repository_database.NewsRepository) *NewsService {
	return &NewsService{newsRepo: newsRepo}
}

func (ns *NewsService) CreateNews(title, description, content, content2 string, author string, category []string) entities.News {
	//categoryType, err := enums.StringToCategoryType(category)
	categories := ns.newsRepo.FindCategoriesByNames(category)
	log.Printf("Categories: %v\n", categories)
	res := ns.newsRepo.CreateNews(title, description, content, content2, categories, author)
	return res
}
func (ns *NewsService) SetBannerPath(mediaPaths []string, eventID uint) {
	ns.newsRepo.UpdateNewsBannerByNewsID(mediaPaths, eventID)
}
func (ns *NewsService) GetNewsByID(newsID uint) (*entities.News, bool) {
	news, found := ns.newsRepo.GetNewsByID(newsID)
	cat := ns.newsRepo.FindNewsCategories(*news)
	news.Categories = cat.Categories
	return news, found
}

func (ns *NewsService) UpdateNews(newsID uint, title, description, content, content2, author string, category []string) (*entities.News, bool) {
	categories := ns.newsRepo.FindCategoriesByNames(category)

	updated, err := ns.newsRepo.UpdateNews(newsID, title, description, content, content2, categories, author)
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

func (ns *NewsService) GetAllNews(categories []string, limit int, offset int) []entities.News {
	news, err := ns.newsRepo.GetAllNews(categories, limit, offset)
	for i, job := range news {
		cat := ns.newsRepo.FindNewsCategories(job)
		news[i].Categories = cat.Categories
	}
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
