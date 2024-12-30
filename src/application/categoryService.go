package application

import (
	"first-project/src/bootstrap"
	"first-project/src/entities"
	repository_database_interfaces "first-project/src/repository/database/interfaces"
)

type categoryService struct {
	constants          *bootstrap.Constants
	categoryRepository repository_database_interfaces.CategoryRepository
}

func NewCategoryService(
	constants *bootstrap.Constants,
	categoryRepository repository_database_interfaces.CategoryRepository,
) *categoryService {
	return &categoryService{
		constants:          constants,
		categoryRepository: categoryRepository,
	}
}

func (categoryService *categoryService) GetListCategoryNames() []string {
	return categoryService.categoryRepository.FindAllCategories()
}

func (categoryService *categoryService) GetCategoriesByName(categoryNames []string) []entities.Category {
	var categories []entities.Category
	for _, name := range categoryNames {
		category := categoryService.categoryRepository.CreateOrGetCategoryByName(name)
		categories = append(categories, category)
	}
	return categories
}
