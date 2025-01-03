package repository_database

import (
	"first-project/src/entities"
	"strings"

	"gorm.io/gorm"
)

type newsRepository struct{}

func NewNewsRepository(db *gorm.DB) *newsRepository {
	return &newsRepository{}
}

func (repo *newsRepository) FindNewsByTitle(db *gorm.DB, name string) (*entities.News, bool) {
	var news entities.News
	result := db.First(&news, "title = ?", name)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &news, true
}

func (repo *newsRepository) FindNewsByID(db *gorm.DB, newsID uint) (*entities.News, bool) {
	var news entities.News
	result := db.First(&news, queryByID, newsID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &news, true
}

func (repo *newsRepository) CreateNews(db *gorm.DB, news *entities.News) error {
	return db.Create(news).Error
}

func (repo *newsRepository) GetNewsByID(db *gorm.DB, id uint) (*entities.News, bool) {
	var news entities.News
	query := db.Where(queryByID, id)
	err := query.First(&news).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(err)
	}
	return &news, true
}

func (repo *newsRepository) UpdateNewsCategories(db *gorm.DB, newsID uint, categories []entities.Category) error {
	return db.Model(&entities.News{ID: newsID}).Association("Categories").Replace(categories)
}

func (repo *newsRepository) UpdateNews(db *gorm.DB, news *entities.News) error {
	return db.Save(news).Error
}

func (repo *newsRepository) DeleteNews(db *gorm.DB, newsID uint) error {
	return db.Unscoped().Delete(&entities.News{}, newsID).Error
}

func (repo *newsRepository) FindAllNews(db *gorm.DB, offset, pageSize int) ([]*entities.News, bool) {
	var news []*entities.News
	result := db.Offset(offset).Limit(pageSize).Find(&news)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return news, true
}

func (repo *newsRepository) FindNewsByCategoryName(db *gorm.DB, categories []string, offset, pageSize int) []*entities.News {
	var news []*entities.News

	result := db.
		Distinct("news.*").
		Joins("JOIN news_categories ON news.id = news_categories.news_id").
		Joins("JOIN categories ON categories.id = news_categories.category_id").
		Where("categories.name IN ?", categories).
		Limit(pageSize).
		Offset(offset).
		Find(&news)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}

	return news
}

func (repo *newsRepository) FindNewsCategoriesByNews(db *gorm.DB, news *entities.News) []entities.Category {
	if err := db.Model(news).Association("Categories").Find(&news.Categories); err != nil {
		panic(err)
	}
	return news.Categories
}

func (repo *newsRepository) FullTextSearch(db *gorm.DB, query string, offset, pageSize int) []*entities.News {
	var news []*entities.News

	db.Exec(`ALTER TABLE news ADD FULLTEXT INDEX idx_title_description_content_content2 (title, description, content, content2)`)
	searchQuery := "+" + strings.Join(strings.Fields(query), "* +") + "*"

	result := db.Model(&entities.News{}).
		Where("MATCH(title, description, content, content2) AGAINST(? IN BOOLEAN MODE)", searchQuery).
		Offset(offset).
		Limit(pageSize).
		Find(&news)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return news
}
