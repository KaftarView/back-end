package application

import (
	"first-project/src/bootstrap"
	"first-project/src/entities"
	repository_database "first-project/src/repository/database"
)

type categoryService struct {
	constants          *bootstrap.Constants
	categoryRepository *repository_database.CategoryRepository
}

func NewCategoryService(
	constants *bootstrap.Constants,
	categoryRepository *repository_database.CategoryRepository,
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
