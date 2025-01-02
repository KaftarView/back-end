package application_interfaces

import "first-project/src/entities"

type CategoryService interface {
	GetCategoriesByName(categoryNames []string) []entities.Category
	GetListCategoryNames() []string
}
