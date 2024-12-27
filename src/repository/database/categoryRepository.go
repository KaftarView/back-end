package repository_database

import (
	"first-project/src/entities"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (repo *CategoryRepository) FindAllCategories() []string {
	var categoryNames []string
	result := repo.db.Model(&entities.Category{}).Pluck("name", &categoryNames)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return []string{}
		}
		panic(result.Error)
	}
	return categoryNames
}

func (repo *CategoryRepository) CreateOrGetCategoryByName(name string) entities.Category {
	var category entities.Category
	if err := repo.db.FirstOrCreate(&category, entities.Category{Name: name}).Error; err != nil {
		panic(err)
	}

	return category
}
