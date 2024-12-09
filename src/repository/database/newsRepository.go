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

func (repo *NewsRepository) CreateNews(news entities.News) entities.News {

	err := repo.db.Create(&news).Error

	if err != nil {
		log.Printf("Error while creating news: %v", err)
		panic(err)
	}
	return news
}
func (repo *NewsRepository) UpdateNewsBannerByNewsID(mediaPaths string, eventID uint) {
	var news entities.News
	if err := repo.db.Model(&news).Where("id = ?", eventID).Update("banner_paths", mediaPaths).Error; err != nil {
		panic(err)
	}
}
func (repo *NewsRepository) GetNewsByID(id uint) (*entities.News, bool) {
	var news entities.News
	query := repo.db.Where("id = ?", id)
	err := query.First(&news).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &news, false
		}
		panic(err)
	}
	log.Print(news.Categories)
	return &news, true
}

func (repo *NewsRepository) UpdateNews(newsID uint, title, description, content, content2 string, categories []entities.Category, author string) (*entities.News, error) {
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
	news.Content2 = content2
	news.Categories = categories
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

func (repo *NewsRepository) GetAllNews(categories []string, limit int, offset int) ([]entities.News, error) {
	var news []entities.News
	if len(categories) == 0 {
		err := repo.db.Limit(limit).
			Offset(offset).
			Find(&news).Error
		if err != nil {
			return nil, err
		}
		return news, nil
	}
	err := repo.db.Joins("JOIN news_categories ON news.id = news_categories.news_id").
		Joins("JOIN categories ON categories.id = news_categories.category_id").
		Where("categories.name IN ?", categories).
		Limit(limit).
		Offset(offset).
		Find(&news).Error
	if err != nil {
		return nil, err
	}
	log.Print(categories)
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

func (repo *NewsRepository) FindNewsCategories(news entities.News) entities.News {
	if err := repo.db.Model(&news).Association("Categories").Find(&news.Categories); err != nil {
		panic(err)
	}
	return news
}

func (repo *NewsRepository) FindCategoriesByNames(categoryNames []string) []entities.Category {
	var categories []entities.Category
	for _, categoryName := range categoryNames {
		var category entities.Category
		if err := repo.db.FirstOrCreate(&category, entities.Category{Name: categoryName}).Error; err != nil {
			panic(err)
		}
		categories = append(categories, category)
	}
	return categories
}
