package controller_v1_general

import (
	application_aws "first-project/src/application/aws"
	application_news "first-project/src/application/news"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/enums"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type NewsController struct {
	constants   *bootstrap.Constants
	newsService *application_news.NewsService
	awsService  *application_aws.S3service
}

func NewNewsController(constants *bootstrap.Constants, newsService *application_news.NewsService,
	awsService *application_aws.S3service) *NewsController {
	return &NewsController{
		constants:   constants,
		newsService: newsService,
		awsService:  awsService,
	}
}

func (nc *NewsController) CreateNews(c *gin.Context) {
	type createParams struct {
		Title       string                `json:"title" validate:"required"`
		Description string                `json:"description"`
		Content     string                `json:"content"`
		Content2    string                `json:"content2"`
		Banner      *multipart.FileHeader `json:"banner"`
		Banner2     *multipart.FileHeader `json:"banner2"`
		Category    []string              `json:"category" validate:"required"`
		Author      string                `json:"author" validate:"required"`
	}
	param := controller.Validated[createParams](c, &nc.constants.Context)

	news := nc.newsService.CreateNews(
		param.Title,
		param.Description,
		param.Content,
		param.Content2,
		param.Author,
		param.Category,
	)

	objectPath := fmt.Sprintf("news/%d/banners/%s", news.ID, param.Banner.Filename)
	nc.awsService.UploadObject(enums.BannersBucket, objectPath, param.Banner)
	BannerPaths := objectPath

	if param.Banner2 != nil {
		objectPath = fmt.Sprintf("news/%d/banners/%s", news.ID, param.Banner2.Filename)
		nc.awsService.UploadObject(enums.BannersBucket, objectPath, param.Banner)
		BannerPaths = BannerPaths + "," + objectPath
	}

	nc.newsService.SetBannerPath(BannerPaths, news.ID)

	trans := controller.GetTranslator(c, nc.constants.Context.Translator)
	message, _ := trans.T("successMessage.NewsCreation")
	controller.Response(c, 201, message, nil)
}

func (nc *NewsController) UpdateNews(c *gin.Context) {
	newsID := c.Param("id")
	id, err := strconv.Atoi(newsID)
	if err != nil {
		trans := controller.GetTranslator(c, nc.constants.Context.Translator)
		message, _ := trans.T("successMessage.NewsNotFound")
		controller.Response(c, 400, message, nil)
		return
	}

	type editParams struct {
		Title       string                `json:"title" validate:"required"`
		Description string                `json:"description"`
		Content     string                `json:"content"`
		Content2    string                `json:"content2"`
		Banner      *multipart.FileHeader `json:"banner"`
		Banner2     *multipart.FileHeader `json:"banner2"`
		Category    []string              `json:"category" validate:"required"`
		Author      string                `json:"author" validate:"required"`
	}

	param := controller.Validated[editParams](c, &nc.constants.Context)

	updatedNewsPointer, found := nc.newsService.UpdateNews(
		uint(id),
		param.Title,
		param.Description,
		param.Content,
		param.Content2,
		param.Author,
		param.Category,
	)

	if !found {
		trans := controller.GetTranslator(c, nc.constants.Context.Translator)
		message, _ := trans.T("successMessage.NewsNotFound")
		controller.Response(c, 400, message, nil)
		return
	}

	NewsBannerPaths := strings.Split(updatedNewsPointer.BannerPaths, ",")

	nc.awsService.DeleteObject(enums.BannersBucket, NewsBannerPaths[0])

	objectPath := fmt.Sprintf("news/%d/banners/%s", id, param.Banner.Filename)
	nc.awsService.UploadObject(enums.BannersBucket, objectPath, param.Banner)
	BannerPaths := objectPath

	if len(NewsBannerPaths) > 1 {
		nc.awsService.DeleteObject(enums.BannersBucket, NewsBannerPaths[1])
	}

	if param.Banner2 != nil {
		objectPath = fmt.Sprintf("news/%d/banners/%s", id, param.Banner2.Filename)
		BannerPaths = BannerPaths + "," + objectPath
	}

	nc.newsService.SetBannerPath(BannerPaths, uint(id))

	trans := controller.GetTranslator(c, nc.constants.Context.Translator)
	message, _ := trans.T("successMessage.NewsUpdated")
	controller.Response(c, 201, message, updatedNewsPointer)
}

func (nc *NewsController) DeleteNews(c *gin.Context) {
	newsIDStr := c.Param("id")
	id, err := strconv.Atoi(newsIDStr)

	if err != nil {
		trans := controller.GetTranslator(c, nc.constants.Context.Translator)
		message, _ := trans.T("successMessage.NewsNotFound")
		controller.Response(c, 400, message, nil)
		return
	}

	News, found := nc.newsService.DeleteNews(uint(id))
	if !found {
		trans := controller.GetTranslator(c, nc.constants.Context.Translator)
		message, _ := trans.T("successMessage.NewsNotFound")
		controller.Response(c, 400, message, nil)
		return
	}

	NewsBannerPaths := strings.Split(News.BannerPaths, ",")
	nc.awsService.DeleteObject(enums.BannersBucket, NewsBannerPaths[0])
	if len(NewsBannerPaths) > 1 {
		nc.awsService.DeleteObject(enums.BannersBucket, NewsBannerPaths[1])
	}

	trans := controller.GetTranslator(c, nc.constants.Context.Translator)
	message, _ := trans.T("successMessage.NewsDeleted")
	controller.Response(c, 201, message, nil)
}

func (nc *NewsController) GetNewsByID(c *gin.Context) {
	newsID := c.Param("id")
	id, err := strconv.Atoi(newsID)
	if err != nil {
		trans := controller.GetTranslator(c, nc.constants.Context.Translator)
		message, _ := trans.T("successMessage.NewsNotFound")
		controller.Response(c, 400, message, nil)
		return
	}

	news, found := nc.newsService.GetNewsByID(uint(id))

	if !found {
		trans := controller.GetTranslator(c, nc.constants.Context.Translator)
		message, _ := trans.T("successMessage.NewsNotFound")
		controller.Response(c, 400, message, nil)
		return
	}

	trans := controller.GetTranslator(c, nc.constants.Context.Translator)
	message, _ := trans.T("successMessage.NewsFound")
	controller.Response(c, 201, message, news)
}
func (nc *NewsController) GetNewsList(c *gin.Context) {
	categories := c.QueryArray("categories")

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		trans := controller.GetTranslator(c, nc.constants.Context.Translator)
		message, _ := trans.T("successMessage.BadRequest")
		controller.Response(c, 400, message, nil)
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		trans := controller.GetTranslator(c, nc.constants.Context.Translator)
		message, _ := trans.T("successMessage.BadRequest")
		controller.Response(c, 400, message, nil)
		return
	}

	newsList := nc.newsService.GetAllNews(categories, limit, offset)

	trans := controller.GetTranslator(c, nc.constants.Context.Translator)
	message, _ := trans.T("successMessage.NewsFound")
	controller.Response(c, 201, message, newsList)
}

func (nc *NewsController) GetTopKNews(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil || limit <= 0 {
		trans := controller.GetTranslator(c, nc.constants.Context.Translator)
		message, _ := trans.T("successMessage.BadRequest")
		controller.Response(c, 400, message, nil)
		return
	}

	categories := c.QueryArray("categories")
	var parsedCategories []enums.CategoryType
	for _, category := range categories {
		parsedCategory, err := strconv.Atoi(category)
		if err != nil {
			trans := controller.GetTranslator(c, nc.constants.Context.Translator)
			message, _ := trans.T("successMessage.BadRequest")
			controller.Response(c, 400, message, nil)
			return
		}
		parsedCategories = append(parsedCategories, enums.CategoryType(parsedCategory))
	}

	topKNews, err := nc.newsService.GetTopKNews(limit, parsedCategories)
	if err != nil {
		trans := controller.GetTranslator(c, nc.constants.Context.Translator)
		message, _ := trans.T("successMessage.BadRequest")
		controller.Response(c, 400, message, nil)
		return
	}
	trans := controller.GetTranslator(c, nc.constants.Context.Translator)
	message, _ := trans.T("successMessage.NewsFound")
	controller.Response(c, 200, message, topKNews)
}

func (nc *NewsController) GetNewsByCategory(c *gin.Context) {
	type requestBody struct {
		Categories []string `json:"categories"`
	}

	categgories := controller.Validated[requestBody](c, &nc.constants.Context)
	newsList := nc.newsService.GetAllNews(categgories.Categories, 10, 0)
	trans := controller.GetTranslator(c, nc.constants.Context.Translator)
	message, _ := trans.T("successMessage.NewsFound")
	controller.Response(c, 200, message, newsList)
}
