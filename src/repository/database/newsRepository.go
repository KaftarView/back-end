package repository_database

import (
	"first-project/src/entities"
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

func (repo *NewsRepository) FindNewsByTitle(name string) (entities.News, bool) {
	var news entities.News
	result := repo.db.First(&news, "title = ?", name)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return news, false
		}
		panic(result.Error)
	}
	return news, true
}

func (repo *NewsRepository) FindNewsByID(newsID uint) (entities.News, bool) {
	var news entities.News
	result := repo.db.First(&news, queryByID, newsID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return news, false
		}
		panic(result.Error)
	}
	return news, true
}

func (repo *NewsRepository) CreateNews(news entities.News) entities.News {
	err := repo.db.Create(&news).Error
	if err != nil {
		panic(err)
	}
	return news
}
func (repo *NewsRepository) UpdateNewsBannerByNewsID(mediaPaths string, eventID uint) {
	var news entities.News
	if err := repo.db.Model(&news).Where(queryByID, eventID).Update("banner_paths", mediaPaths).Error; err != nil {
		panic(err)
	}
}
func (repo *NewsRepository) GetNewsByID(id uint) (*entities.News, bool) {
	var news entities.News
	query := repo.db.Where(queryByID, id)
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

func (repo *NewsRepository) UpdateNews(news entities.News) {
	err := repo.db.Save(&news).Error
	if err != nil {
		panic(err)
	}
}

func (repo *NewsRepository) DeleteNews(newsID uint) {
	err := repo.db.Unscoped().Delete(&entities.News{}, newsID).Error
	if err != nil {
		panic(err)
	}
}

func (repo *NewsRepository) FindAllNews(offset, pageSize int) ([]entities.News, bool) {
	var news []entities.News
	result := repo.db.Offset(offset).Limit(pageSize).Find(&news)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return news, false
		}
		panic(result.Error)
	}
	return news, true
}

func (repo *NewsRepository) FindNewsByCategoryName(categories []string, offset, pageSize int) []entities.News {
	var news []entities.News

	result := repo.db.
		Joins("JOIN news_categories ON news.id = news_categories.news_id").
		Joins("JOIN categories ON categories.id = news_categories.category_id").
		Where("categories.name IN ?", categories).
		Limit(pageSize).
		Offset(offset).
		Find(&news)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return []entities.News{}
		}
		panic(result.Error)
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

func (repo *NewsRepository) FindNewsCategoriesByNews(news entities.News) []entities.Category {
	if err := repo.db.Model(&news).Association("Categories").Find(&news.Categories); err != nil {
		panic(err)
	}
	return news.Categories
}
