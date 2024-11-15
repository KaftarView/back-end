package repository_database

import (
	"first-project/src/entities"
	"first-project/src/enums"
	"time"

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

func (repo *NewsRepository) CreateNews(news entities.News) (entities.News, error) {
	err := repo.db.Create(&news).Error
	if err != nil {
		return entities.News{}, err
	}
	return news, nil
}

func (repo *NewsRepository) GetNewsByID(id uint) (entities.News, error) {
	var news entities.News
	query := repo.db.Where("id = ?", id)
	err := query.First(&news).Error
	if err != nil {
		return entities.News{}, err
	}
	return news, nil
}

func (repo *NewsRepository) UpdateNews(news entities.News) (entities.News, error) {
	err := repo.db.Save(&news).Error
	if err != nil {
		return entities.News{}, err
	}
	return news, nil
}

func (repo *NewsRepository) DeleteNews(id uint) error {
	query := repo.db.Where("id = ?", id)
	err := query.Delete(&entities.News{}).Error
	return err
}

func (repo *NewsRepository) GetAllNews(categories []enums.CategoryType, limit int, offset int) ([]entities.News, error) {
	var news []entities.News
	query := repo.db
	if len(categories) > 0 {
		query = query.Where("category IN ?", categories)
	}
	err := query.Limit(limit).Offset(offset).Find(&news).Error
	if err != nil {
		return nil, err
	}
	return news, nil
}

func (repo *NewsRepository) GetNewsByDateRange(start, end time.Time, categories []enums.CategoryType) ([]entities.News, error) {
	var news []entities.News
	query := repo.db.Where("published_at BETWEEN ? AND ?", start, end)
	if len(categories) > 0 {
		query = query.Where("category IN ?", categories)
	}
	err := query.Find(&news).Error
	if err != nil {
		return nil, err
	}
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
