package controller_v1_news

import (
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"
	"first-project/src/controller"

	"github.com/gin-gonic/gin"
)

type GeneralNewsController struct {
	constants   *bootstrap.Constants
	newsService application_interfaces.NewsService
}

func NewGeneralNewsController(
	constants *bootstrap.Constants,
	newsService application_interfaces.NewsService,
) *GeneralNewsController {
	return &GeneralNewsController{
		constants:   constants,
		newsService: newsService,
	}
}

func (generalNewsController *GeneralNewsController) GetNewsList(c *gin.Context) {
	pagination := controller.GetPagination(c, &generalNewsController.constants.Context)
	news := generalNewsController.newsService.GetNewsList(pagination.Page, pagination.PageSize)

	controller.Response(c, 200, "", news)
}

func (generalNewsController *GeneralNewsController) GetNewsDetails(c *gin.Context) {
	type getNewsDetailsParams struct {
		NewsID uint `uri:"newsID" validate:"required"`
	}
	param := controller.Validated[getNewsDetailsParams](c, &generalNewsController.constants.Context)
	news := generalNewsController.newsService.GetNewsDetails(param.NewsID)

	controller.Response(c, 200, "", news)
}

func (generalNewsController *GeneralNewsController) SearchNews(c *gin.Context) {
	type searchNewsForAdminParams struct {
		Query string `form:"query"`
	}
	param := controller.Validated[searchNewsForAdminParams](c, &generalNewsController.constants.Context)
	pagination := controller.GetPagination(c, &generalNewsController.constants.Context)
	news := generalNewsController.newsService.SearchNews(param.Query, pagination.Page, pagination.PageSize)

	controller.Response(c, 200, "", news)
}

func (generalNewsController *GeneralNewsController) FilterNewsByCategory(c *gin.Context) {
	type filterNewsByCategoryParams struct {
		Categories []string `form:"categories"`
	}
	param := controller.Validated[filterNewsByCategoryParams](c, &generalNewsController.constants.Context)
	pagination := controller.GetPagination(c, &generalNewsController.constants.Context)
	news := generalNewsController.newsService.FilterNewsByCategory(param.Categories, pagination.Page, pagination.PageSize)

	controller.Response(c, 200, "", news)
}
