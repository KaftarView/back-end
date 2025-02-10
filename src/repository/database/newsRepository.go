package repository_database

import (
	"first-project/src/entities"
	"strings"

	"gorm.io/gorm"
)

type NewsRepository struct{}

func NewNewsRepository() *NewsRepository {
	return &NewsRepository{}
}

func (repo *NewsRepository) FindNewsByTitle(db *gorm.DB, name string) (*entities.News, bool) {
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

func (repo *NewsRepository) FindNewsByID(db *gorm.DB, newsID uint) (*entities.News, bool) {
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

func (repo *NewsRepository) CreateNews(db *gorm.DB, news *entities.News) error {
	return db.Create(news).Error
}

func (repo *NewsRepository) GetNewsByID(db *gorm.DB, id uint) (*entities.News, bool) {
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

func (repo *NewsRepository) UpdateNewsCategories(db *gorm.DB, newsID uint, categories []entities.Category) error {
	return db.Model(&entities.News{ID: newsID}).Association("Categories").Replace(categories)
}

func (repo *NewsRepository) UpdateNews(db *gorm.DB, news *entities.News) error {
	return db.Save(news).Error
}

func (repo *NewsRepository) DeleteNews(db *gorm.DB, newsID uint) error {
	return db.Unscoped().Delete(&entities.News{}, newsID).Error
}

func (repo *NewsRepository) FindAllNews(db *gorm.DB, offset, pageSize int) ([]*entities.News, bool) {
	var news []*entities.News
	result := OrderByCreatedAtDesc(db).Offset(offset).Limit(pageSize).Find(&news)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return news, true
}

func (repo *NewsRepository) FindNewsByCategoryName(db *gorm.DB, categories []string, offset, pageSize int) []*entities.News {
	var news []*entities.News

	result := OrderByCreatedAtDesc(db).
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

func (repo *NewsRepository) FindNewsCategoriesByNews(db *gorm.DB, news *entities.News) []entities.Category {
	if err := OrderByCreatedAtDesc(db).Model(news).Association("Categories").Find(&news.Categories); err != nil {
		panic(err)
	}
	return news.Categories
}

func (repo *NewsRepository) FullTextSearch(db *gorm.DB, query string, offset, pageSize int) []*entities.News {
	var news []*entities.News

	db.Exec(`ALTER TABLE news ADD FULLTEXT INDEX idx_title_description_content_content2 (title, description, content, content2)`)
	searchQuery := "+" + strings.Join(strings.Fields(query), "* +") + "*"

	result := OrderByCreatedAtDesc(db).
		Model(&entities.News{}).
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
