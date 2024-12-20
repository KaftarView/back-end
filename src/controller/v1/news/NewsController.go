package controller_v1_news

import (
	application_news "first-project/src/application/news"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/dto"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type NewsController struct {
	constants   *bootstrap.Constants
	newsService *application_news.NewsService
}

func NewNewsController(
	constants *bootstrap.Constants,
	newsService *application_news.NewsService,
) *NewsController {
	return &NewsController{
		constants:   constants,
		newsService: newsService,
	}
}

func (newsController *NewsController) CreateNews(c *gin.Context) {
	type createNewsParams struct {
		Title       string                `form:"title" validate:"required"`
		Description string                `form:"description"`
		Content     string                `form:"content"`
		Content2    string                `form:"content2"`
		Banner      *multipart.FileHeader `form:"banner"`
		Banner2     *multipart.FileHeader `form:"banner2"`
		Categories  []string              `form:"categories"`
	}
	param := controller.Validated[createNewsParams](c, &newsController.constants.Context)
	userID, _ := c.Get(newsController.constants.Context.UserID)
	newsDetails := dto.RequestNewsDetails{
		Title:       param.Title,
		Description: param.Description,
		Content:     param.Content,
		Content2:    param.Content2,
		Banner:      param.Banner,
		Banner2:     param.Banner2,
		Categories:  param.Categories,
		AuthorID:    userID.(uint),
	}
	news := newsController.newsService.CreateNews(newsDetails)

	trans := controller.GetTranslator(c, newsController.constants.Context.Translator)
	message, _ := trans.T("successMessage.NewsCreation")
	controller.Response(c, 200, message, news.ID)
}

func (newsController *NewsController) UpdateNews(c *gin.Context) {
	type editNewsParams struct {
		Title       *string               `form:"title"`
		Description *string               `form:"description"`
		Content     *string               `form:"content"`
		Content2    *string               `form:"content2"`
		Banner      *multipart.FileHeader `form:"banner"`
		Banner2     *multipart.FileHeader `form:"banner2"`
		Categories  *[]string             `form:"categories"`
		NewsID      uint                  `uri:"newsID" validate:"required"`
	}
	param := controller.Validated[editNewsParams](c, &newsController.constants.Context)
	newsUpdatingDetails := dto.RequestUpdateNewsDetails{
		ID:          param.NewsID,
		Title:       param.Title,
		Description: param.Description,
		Content:     param.Content,
		Content2:    param.Content2,
		Banner:      param.Banner,
		Banner2:     param.Banner2,
		Categories:  param.Categories,
	}
	newsController.newsService.UpdateNews(newsUpdatingDetails)

	trans := controller.GetTranslator(c, newsController.constants.Context.Translator)
	message, _ := trans.T("successMessage.NewsUpdated")
	controller.Response(c, 200, message, nil)
}

func (newsController *NewsController) DeleteNews(c *gin.Context) {
	type deleteNewsParams struct {
		NewsID uint `uri:"newsID" validate:"required"`
	}
	param := controller.Validated[deleteNewsParams](c, &newsController.constants.Context)
	newsController.newsService.DeleteNews(param.NewsID)

	trans := controller.GetTranslator(c, newsController.constants.Context.Translator)
	message, _ := trans.T("successMessage.NewsDeleted")
	controller.Response(c, 200, message, nil)
}

func (newsController *NewsController) GetNewsDetails(c *gin.Context) {
	type getNewsDetailsParams struct {
		NewsID uint `uri:"newsID" validate:"required"`
	}
	param := controller.Validated[getNewsDetailsParams](c, &newsController.constants.Context)
	news := newsController.newsService.GetNewsDetails(param.NewsID)

	controller.Response(c, 200, "", news)
}

func (newsController *NewsController) GetNewsList(c *gin.Context) {
	type getNewsListParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	param := controller.Validated[getNewsListParams](c, &newsController.constants.Context)
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	news := newsController.newsService.GetNewsList(param.Page, param.PageSize)

	controller.Response(c, 200, "", news)
}

func (newsController *NewsController) SearchNews(c *gin.Context) {
	type searchNewsForAdminParams struct {
		Query    string `form:"query"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
	}
	param := controller.Validated[searchNewsForAdminParams](c, &newsController.constants.Context)
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	news := newsController.newsService.SearchNews(param.Query, param.Page, param.PageSize)

	controller.Response(c, 200, "", news)
}

func (newsController *NewsController) FilterNewsByCategory(c *gin.Context) {
	type filterNewsByCategoryParams struct {
		Categories []string `form:"categories"`
		Page       int      `form:"page"`
		PageSize   int      `form:"pageSize"`
	}
	param := controller.Validated[filterNewsByCategoryParams](c, &newsController.constants.Context)
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	news := newsController.newsService.FilterNewsByCategory(param.Categories, param.Page, param.PageSize)

	controller.Response(c, 200, "", news)
}
