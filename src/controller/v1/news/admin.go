package controller_v1_news

import (
	"first-project/src/application"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/dto"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type AdminNewsController struct {
	constants   *bootstrap.Constants
	newsService *application.NewsService
}

func NewAdminNewsController(
	constants *bootstrap.Constants,
	newsService *application.NewsService,
) *AdminNewsController {
	return &AdminNewsController{
		constants:   constants,
		newsService: newsService,
	}
}

func (adminNewsController *AdminNewsController) CreateNews(c *gin.Context) {
	type createNewsParams struct {
		Title       string                `form:"title" validate:"required"`
		Description string                `form:"description"`
		Content     string                `form:"content"`
		Content2    string                `form:"content2"`
		Banner      *multipart.FileHeader `form:"banner"`
		Banner2     *multipart.FileHeader `form:"banner2"`
		Categories  []string              `form:"categories"`
	}
	param := controller.Validated[createNewsParams](c, &adminNewsController.constants.Context)
	userID, _ := c.Get(adminNewsController.constants.Context.UserID)
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
	news := adminNewsController.newsService.CreateNews(newsDetails)

	trans := controller.GetTranslator(c, adminNewsController.constants.Context.Translator)
	message, _ := trans.T("successMessage.NewsCreation")
	controller.Response(c, 200, message, news.ID)
}

func (adminNewsController *AdminNewsController) UpdateNews(c *gin.Context) {
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
	param := controller.Validated[editNewsParams](c, &adminNewsController.constants.Context)
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
	adminNewsController.newsService.UpdateNews(newsUpdatingDetails)

	trans := controller.GetTranslator(c, adminNewsController.constants.Context.Translator)
	message, _ := trans.T("successMessage.NewsUpdated")
	controller.Response(c, 200, message, nil)
}

func (adminNewsController *AdminNewsController) DeleteNews(c *gin.Context) {
	type deleteNewsParams struct {
		NewsID uint `uri:"newsID" validate:"required"`
	}
	param := controller.Validated[deleteNewsParams](c, &adminNewsController.constants.Context)
	adminNewsController.newsService.DeleteNews(param.NewsID)

	trans := controller.GetTranslator(c, adminNewsController.constants.Context.Translator)
	message, _ := trans.T("successMessage.NewsDeleted")
	controller.Response(c, 200, message, nil)
}
