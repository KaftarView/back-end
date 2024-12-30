package repository_database_interfaces

import "first-project/src/entities"

type CategoryRepository interface {
	CreateOrGetCategoryByName(name string) entities.Category
	FindAllCategories() []string
}
