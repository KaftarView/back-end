package repository_database

import (
	"first-project/src/entities"

	"gorm.io/gorm"
)

type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (repo *CategoryRepository) FindAllCategories(db *gorm.DB) []string {
	var categoryNames []string
	result := db.Model(&entities.Category{}).Order("created_at DESC").Pluck("name", &categoryNames)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return []string{}
		}
		panic(result.Error)
	}
	return categoryNames
}

func (repo *CategoryRepository) CreateOrGetCategoryByName(db *gorm.DB, name string) entities.Category {
	var category entities.Category
	if err := db.FirstOrCreate(&category, entities.Category{Name: name}).Error; err != nil {
		panic(err)
	}

	return category
}
