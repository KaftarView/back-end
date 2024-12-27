package controller_v1_category

import (
	"first-project/src/application"
	"first-project/src/controller"

	"github.com/gin-gonic/gin"
)

type GeneralCategoryController struct {
	categoryService *application.CategoryService
}

func NewGeneralCategoryController(categoryService *application.CategoryService) *GeneralCategoryController {
	return &GeneralCategoryController{categoryService: categoryService}
}

func (generalCategoryController *GeneralCategoryController) GetListCategoryNames(c *gin.Context) {
	categoryNames := generalCategoryController.categoryService.GetListCategoryNames()

	controller.Response(c, 200, "", categoryNames)
}
