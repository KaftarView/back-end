package controller_v1_general

import (
	application_news "first-project/src/application/news"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/enums"
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
	type createParams struct {
		Title       string `json:"title" validate:"required"`
		Description string `json:"description"`
		Content     string `json:"content"`
		ImageURL    string `json:"image_url"`
		Category    string `json:"category" validate:"required"`
		Author      string `json:"author" validate:"required"`
	}
	param := controller.Validated[createParams](c, &nc.constants.Context)

	nc.newsService.CreateNews(
		param.Title,
		param.Description,
		param.Content,
		param.ImageURL,
		param.Category,
		param.Author,
	)

	trans := controller.GetTranslator(c, nc.constants.Context.Translator)
	message, _ := trans.T("successMessage.NewsCreation")
	controller.Response(c, 201, message, nil)
}

func (nc *NewsController) UpdateNews(c *gin.Context) {
	newsID := c.Param("id")
	id, err := strconv.Atoi(newsID)
	if err != nil {
		controller.Response(c, 400, "Invalid news ID", nil)
		return
	}

	type editParams struct {
		Title       string `json:"title" validate:"required"`
		Description string `json:"description"`
		Content     string `json:"content"`
		ImageURL    string `json:"image_url"`
		Category    string `json:"category" validate:"required"`
		Author      string `json:"author" validate:"required"`
	}

	param := controller.Validated[editParams](c, &nc.constants.Context)

	updatedNewsPointer, found := nc.newsService.UpdateNews(
		uint(id),
		param.Title,
		param.Description,
		param.Content,
		param.ImageURL,
		param.Category,
		param.Author,
	)

	if !found {
		controller.Response(c, 400, "No news with this id", nil)
		return
	}

	trans := controller.GetTranslator(c, nc.constants.Context.Translator)
	message, _ := trans.T("successMessage.NewsUpdated")
	controller.Response(c, 201, message, updatedNewsPointer)
}

func (nc *NewsController) DeleteNews(c *gin.Context) {
	newsIDStr := c.Param("id")
	id, err := strconv.Atoi(newsIDStr)
	if err != nil {
		controller.Response(c, 400, "Invalid news ID", nil)
		return
	}

	found := nc.newsService.DeleteNews(uint(id))
	if !found {
		controller.Response(c, 400, "Invalid news ID", nil)
		return
	}

	controller.Response(c, 200, "News deleted successfully", nil)
}

func (nc *NewsController) GetNewsByID(c *gin.Context) {
	newsID := c.Param("id")
	id, err := strconv.Atoi(newsID)
	if err != nil {
		controller.Response(c, 400, "Invalid news ID", nil)
		return
	}

	news, found := nc.newsService.GetNewsByID(uint(id))
	if !found {
		controller.Response(c, 400, "No news with this id", nil)
		return
	}

	controller.Response(c, 200, "Success", news)
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

	newsList := nc.newsService.GetAllNews(parsedCategories, limit, offset)

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

	newsList := nc.newsService.GetAllNews(categories, 10, 0)
	controller.Response(c, 200, "Success", newsList)
}
