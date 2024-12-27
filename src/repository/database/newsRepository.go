package repository_database

import (
	"first-project/src/entities"
	"strings"

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

func (repo *NewsRepository) FindNewsByTitle(name string) (*entities.News, bool) {
	var news entities.News
	result := repo.db.First(&news, "title = ?", name)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &news, true
}

func (repo *NewsRepository) FindNewsByID(newsID uint) (*entities.News, bool) {
	var news entities.News
	result := repo.db.First(&news, queryByID, newsID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &news, true
}

func (repo *NewsRepository) CreateNews(news *entities.News) *entities.News {
	err := repo.db.Create(news).Error
	if err != nil {
		panic(err)
	}
	return news
}

func (repo *NewsRepository) GetNewsByID(id uint) (*entities.News, bool) {
	var news entities.News
	query := repo.db.Where(queryByID, id)
	err := query.First(&news).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(err)
	}
	return &news, true
}

func (repo *NewsRepository) UpdateNewsCategories(newsID uint, categories []entities.Category) {
	err := repo.db.Model(&entities.News{ID: newsID}).Association("Categories").Replace(categories)
	if err != nil {
		panic(err)
	}
}

func (repo *NewsRepository) UpdateNews(news *entities.News) {
	err := repo.db.Save(news).Error
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

func (repo *NewsRepository) FindAllNews(offset, pageSize int) ([]*entities.News, bool) {
	var news []*entities.News
	result := repo.db.Offset(offset).Limit(pageSize).Find(&news)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return news, true
}

func (repo *NewsRepository) FindNewsByCategoryName(categories []string, offset, pageSize int) []*entities.News {
	var news []*entities.News

	result := repo.db.
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

func (repo *NewsRepository) FindNewsCategoriesByNews(news *entities.News) []entities.Category {
	if err := repo.db.Model(news).Association("Categories").Find(&news.Categories); err != nil {
		panic(err)
	}
	return news.Categories
}

func (repo *NewsRepository) FullTextSearch(query string, offset, pageSize int) []*entities.News {
	var news []*entities.News

	repo.db.Exec(`ALTER TABLE news ADD FULLTEXT INDEX idx_title_description_content_content2 (title, description, content, content2)`)
	searchQuery := "+" + strings.Join(strings.Fields(query), "* +") + "*"

	result := repo.db.Model(&entities.News{}).
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
