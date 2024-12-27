package controller_v1_news

import (
	"first-project/src/application"
	"first-project/src/bootstrap"
	"first-project/src/controller"

	"github.com/gin-gonic/gin"
)

type GeneralNewsController struct {
	constants   *bootstrap.Constants
	newsService *application.NewsService
}

func NewGeneralNewsController(
	constants *bootstrap.Constants,
	newsService *application.NewsService,
) *GeneralNewsController {
	return &GeneralNewsController{
		constants:   constants,
		newsService: newsService,
	}
}

func (generalNewsController *GeneralNewsController) GetNewsList(c *gin.Context) {
	type getNewsListParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	param := controller.Validated[getNewsListParams](c, &generalNewsController.constants.Context)
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	news := generalNewsController.newsService.GetNewsList(param.Page, param.PageSize)

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
		Query    string `form:"query"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
	}
	param := controller.Validated[searchNewsForAdminParams](c, &generalNewsController.constants.Context)
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	news := generalNewsController.newsService.SearchNews(param.Query, param.Page, param.PageSize)

	controller.Response(c, 200, "", news)
}

func (generalNewsController *GeneralNewsController) FilterNewsByCategory(c *gin.Context) {
	type filterNewsByCategoryParams struct {
		Categories []string `form:"categories"`
		Page       int      `form:"page"`
		PageSize   int      `form:"pageSize"`
	}
	param := controller.Validated[filterNewsByCategoryParams](c, &generalNewsController.constants.Context)
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	news := generalNewsController.newsService.FilterNewsByCategory(param.Categories, param.Page, param.PageSize)

	controller.Response(c, 200, "", news)
}
