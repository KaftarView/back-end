package repository_database

import (
	"first-project/src/entities"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *categoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (repo *categoryRepository) FindAllCategories() []string {
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

func (repo *categoryRepository) CreateOrGetCategoryByName(name string) entities.Category {
	var category entities.Category
	if err := repo.db.FirstOrCreate(&category, entities.Category{Name: name}).Error; err != nil {
		panic(err)
	}

	return category
}
