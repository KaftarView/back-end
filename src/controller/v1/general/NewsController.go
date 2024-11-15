package controller_v1_general

import (
	application_news "first-project/src/application/news"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/entities"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NewsController struct {
	constants   *bootstrap.Constants
	newsService *application_news.NewsService
}

func NewNewsController(constants *bootstrap.Constants, newsService *application_news.NewsService) *NewsController {
	return &NewsController{
		constants:   constants,
		newsService: newsService,
	}
}

func (nc *NewsController) CreateNews(c *gin.Context) {
	var newNews entities.News
	if err := c.ShouldBindJSON(&newNews); err != nil {
		controller.Response(c, 400, "Invalid input", nil)
		return
	}

	createdNews, err := nc.newsService.CreateNews(newNews)
	if err != nil {
		controller.Response(c, 500, "Error creating news", nil)
		return
	}

	controller.Response(c, 201, "News created successfully", createdNews)
}

func (nc *NewsController) GetNewsByID(c *gin.Context) {
	newsID := c.Param("id")
	id, err := strconv.Atoi(newsID)
	if err != nil {
		controller.Response(c, 400, "Invalid news ID", nil)
		return
	}

	news, err := nc.newsService.GetNewsByID(uint(id))
	if err != nil {
		//notFoundError := exceptions.NewNotFoundError("News not found") // should be handled in exceptions
		panic(err)
	}

	controller.Response(c, 200, "Success", news)
}

func (nc *NewsController) UpdateNews(c *gin.Context) {
	newsID := c.Param("id")
	id, err := strconv.Atoi(newsID)
	if err != nil {
		controller.Response(c, 400, "Invalid news ID", nil)
		return
	}

	var updatedNews entities.News
	if err := c.ShouldBindJSON(&updatedNews); err != nil {
		controller.Response(c, 400, "Invalid input", nil)
		return
	}

	updatedNewsPointer, err := nc.newsService.UpdateNews(uint(id), updatedNews)

	if err != nil {
		controller.Response(c, 500, "Error updating news", nil)
		return
	}

	updatedNews = *updatedNewsPointer
	controller.Response(c, 200, "News updated successfully", updatedNews)
}

func (nc *NewsController) DeleteNews(c *gin.Context) {
	newsIDStr := c.Param("id")
	id, err := strconv.Atoi(newsIDStr)
	if err != nil {
		controller.Response(c, 400, "Invalid news ID", nil)
		return
	}

	newsID := uint(id)
	err = nc.newsService.DeleteNews(newsID)
	if err != nil {
		panic(err)
	}

	controller.Response(c, 200, "News deleted successfully", nil)
}
