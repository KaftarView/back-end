package application

import (
	"first-project/src/bootstrap"
	"first-project/src/entities"
	repository_database_interfaces "first-project/src/repository/database/interfaces"

	"gorm.io/gorm"
)

type CategoryService struct {
	constants          *bootstrap.Constants
	categoryRepository repository_database_interfaces.CategoryRepository
	db                 *gorm.DB
}

func NewCategoryService(
	constants *bootstrap.Constants,
	categoryRepository repository_database_interfaces.CategoryRepository,
	db *gorm.DB,
) *CategoryService {
	return &CategoryService{
		constants:          constants,
		categoryRepository: categoryRepository,
		db:                 db,
	}
}

func (categoryService *CategoryService) GetListCategoryNames() []string {
	return categoryService.categoryRepository.FindAllCategories(categoryService.db)
}

func (categoryService *CategoryService) GetCategoriesByName(categoryNames []string) []entities.Category {
	var categories []entities.Category
	for _, name := range categoryNames {
		category := categoryService.categoryRepository.CreateOrGetCategoryByName(categoryService.db, name)
		categories = append(categories, category)
	}
	return categories
}
