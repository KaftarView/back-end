package repository_database_interfaces

import (
	"first-project/src/entities"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateOrGetCategoryByName(db *gorm.DB, name string) entities.Category
	FindAllCategories(db *gorm.DB) []string
}
