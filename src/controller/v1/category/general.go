package controller_v1_category

import (
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/controller"

	"github.com/gin-gonic/gin"
)

type GeneralCategoryController struct {
	categoryService application_interfaces.CategoryService
}

func NewGeneralCategoryController(categoryService application_interfaces.CategoryService) *GeneralCategoryController {
	return &GeneralCategoryController{categoryService: categoryService}
}

func (generalCategoryController *GeneralCategoryController) GetListCategoryNames(c *gin.Context) {
	categoryNames := generalCategoryController.categoryService.GetListCategoryNames()

	controller.Response(c, 200, "", categoryNames)
}
