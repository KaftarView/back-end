package controller_v1_general

import (
	application_news "first-project/src/application/news"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/entities"
	"first-project/src/enums"
	"first-project/src/exceptions"
	"log"
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
		notFoundError := exceptions.NewNotFoundInDatabaseError() // should be handled in exceptions?
		panic(notFoundError)
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

func (nc *NewsController) GetNewsList(c *gin.Context) {
	categories := c.QueryArray("categories")
	var parsedCategories []enums.CategoryType
	for _, category := range categories {
		parsedCategory, err := strconv.Atoi(category)
		if err != nil {
			controller.Response(c, 400, "Invalid category", nil)
			return
		}
		parsedCategories = append(parsedCategories, enums.CategoryType(parsedCategory))
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		controller.Response(c, 400, "Invalid limit", nil)
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		controller.Response(c, 400, "Invalid offset", nil)
		return
	}

	newsList, err := nc.newsService.GetAllNews(parsedCategories, limit, offset)
	if err != nil {
		controller.Response(c, 500, "Error fetching news list", nil)
		return
	}

	controller.Response(c, 200, "Success", newsList)
}

func (nc *NewsController) GetTopKNews(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil || limit <= 0 {
		controller.Response(c, 400, "Invalid limit", nil)
		return
	}

	categories := c.QueryArray("categories")
	var parsedCategories []enums.CategoryType
	for _, category := range categories {
		parsedCategory, err := strconv.Atoi(category)
		if err != nil {
			controller.Response(c, 400, "Invalid category", nil)
			return
		}
		parsedCategories = append(parsedCategories, enums.CategoryType(parsedCategory))
	}

	topKNews, err := nc.newsService.GetTopKNews(limit, parsedCategories)
	if err != nil {
		controller.Response(c, 500, "Error fetching top K news", nil)
		return
	}

	controller.Response(c, 200, "Success", topKNews)
}
func (nc *NewsController) GetNewsByCategory(c *gin.Context) {
	var requestBody struct {
		Categories []string `json:"categories"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		controller.Response(c, 400, "Invalid input", nil)
		return
	}

	var categories []enums.CategoryType
	for _, categoryName := range requestBody.Categories {
		category, err := enums.GetCategoryTypeByName(categoryName)
		if err != nil {
			controller.Response(c, 400, "Invalid category name: "+categoryName, nil)
			return
		}
		categories = append(categories, category)
	}
	var c2 enums.CategoryType = enums.Public
	log.Printf("C3 %v ", uint(c2))
	log.Printf("Categories %v ", categories)
	newsList, err := nc.newsService.GetAllNews(categories, 10, 0) // Limit and Offset can be dynamic
	if err != nil {
		controller.Response(c, 500, "Error fetching news by category", nil)
		return
	}

	controller.Response(c, 200, "Success", newsList)
}
