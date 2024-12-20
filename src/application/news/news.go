package application_news

import (
	application_aws "first-project/src/application/aws"
	"first-project/src/bootstrap"
	"first-project/src/dto"
	"first-project/src/entities"
	"first-project/src/enums"
	"first-project/src/exceptions"
	repository_database "first-project/src/repository/database"
	"fmt"
	"time"
)

type NewsService struct {
	constants         *bootstrap.Constants
	awsS3Service      *application_aws.S3service
	commentRepository *repository_database.CommentRepository
	newsRepository    *repository_database.NewsRepository
	userRepository    *repository_database.UserRepository
}

func NewNewsService(
	constants *bootstrap.Constants,
	awsS3Service *application_aws.S3service,
	commentRepository *repository_database.CommentRepository,
	newsRepository *repository_database.NewsRepository,
	userRepository *repository_database.UserRepository,
) *NewsService {
	return &NewsService{
		constants:         constants,
		awsS3Service:      awsS3Service,
		commentRepository: commentRepository,
		newsRepository:    newsRepository,
		userRepository:    userRepository,
	}
}

const bannerPathFormat = "banners/podcasts/%d/images/%s"

func (newsService *NewsService) CreateNews(newsDetails dto.RequestNewsDetails) *entities.News {
	var conflictError exceptions.ConflictError
	_, newsExist := newsService.newsRepository.FindNewsByTitle(newsDetails.Title)
	if newsExist {
		conflictError.AppendError(
			newsService.constants.ErrorField.Tittle,
			newsService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}

	categories := newsService.newsRepository.FindCategoriesByNames(newsDetails.Categories)
	commentable := newsService.commentRepository.CreateNewCommentable()

	bannerPath := fmt.Sprintf(bannerPathFormat, commentable.CID, newsDetails.Banner.Filename)
	newsService.awsS3Service.UploadObject(enums.BannersBucket, bannerPath, newsDetails.Banner)

	banner2Path := ""
	if newsDetails.Banner2 != nil {
		banner2Path = fmt.Sprintf(bannerPathFormat, commentable.CID, newsDetails.Banner2.Filename)
		newsService.awsS3Service.UploadObject(enums.BannersBucket, banner2Path, newsDetails.Banner2)
	}

	newsModel := &entities.News{
		ID:          commentable.CID,
		Title:       newsDetails.Title,
		Description: newsDetails.Description,
		Content:     newsDetails.Content,
		Content2:    newsDetails.Content2,
		AuthorID:    newsDetails.AuthorID,
		Categories:  categories,
		BannerPath:  bannerPath,
		Banner2Path: banner2Path,
	}

	newsService.newsRepository.CreateNews(newsModel)

	return newsModel
}

func (newsService *NewsService) UpdateNews(newsDetails dto.RequestUpdateNewsDetails) {
	var conflictError exceptions.ConflictError
	var notFoundError exceptions.NotFoundError
	news, newsExist := newsService.newsRepository.FindNewsByID(newsDetails.ID)
	if !newsExist {
		notFoundError.ErrorField = newsService.constants.ErrorField.News
		panic(notFoundError)
	}

	if newsDetails.Title != nil {
		_, newsExist := newsService.newsRepository.FindNewsByTitle(*newsDetails.Title)
		if newsExist {
			conflictError.AppendError(
				newsService.constants.ErrorField.Tittle,
				newsService.constants.ErrorTag.AlreadyExist)
			panic(conflictError)
		}
		news.Title = *newsDetails.Title
	}
	if newsDetails.Description != nil {
		news.Description = *newsDetails.Description
	}
	if newsDetails.Content != nil {
		news.Content = *newsDetails.Content
	}
	if newsDetails.Content2 != nil {
		news.Content2 = *newsDetails.Content2
	}
	if newsDetails.Banner != nil {
		if news.BannerPath != "" {
			newsService.awsS3Service.DeleteObject(enums.BannersBucket, news.BannerPath)
		}
		banner1Path := fmt.Sprintf(bannerPathFormat, news.ID, newsDetails.Banner.Filename)
		newsService.awsS3Service.UploadObject(enums.BannersBucket, banner1Path, newsDetails.Banner)
		news.BannerPath = banner1Path
	}
	if newsDetails.Banner2 != nil {
		if news.Banner2Path != "" {
			newsService.awsS3Service.DeleteObject(enums.BannersBucket, news.Banner2Path)
		}
		banner2Path := fmt.Sprintf(bannerPathFormat, news.ID, newsDetails.Banner2.Filename)
		newsService.awsS3Service.UploadObject(enums.BannersBucket, banner2Path, newsDetails.Banner2)
		news.Banner2Path = banner2Path
	}
	if newsDetails.Categories != nil {
		news.Categories = newsService.newsRepository.FindCategoriesByNames(*newsDetails.Categories)
	}

	newsService.newsRepository.UpdateNews(news)
}

func (newsService *NewsService) DeleteNews(newsID uint) {
	var notFoundError exceptions.NotFoundError
	news, newsExist := newsService.newsRepository.FindNewsByID(newsID)
	if !newsExist {
		notFoundError.ErrorField = newsService.constants.ErrorField.News
		panic(notFoundError)
	}
	if news.BannerPath != "" {
		newsService.awsS3Service.DeleteObject(enums.BannersBucket, news.BannerPath)
	}
	if news.Banner2Path != "" {
		newsService.awsS3Service.DeleteObject(enums.BannersBucket, news.Banner2Path)
	}
	newsService.newsRepository.DeleteNews(newsID)
}

func (newsService *NewsService) GetNewsDetails(newsID uint) dto.NewsDetailsResponse {
	var notFoundError exceptions.NotFoundError
	news, newsExist := newsService.newsRepository.FindNewsByID(newsID)
	if !newsExist {
		notFoundError.ErrorField = newsService.constants.ErrorField.News
		panic(notFoundError)
	}

	banner1URL := ""
	banner2URL := ""
	if news.BannerPath != "" {
		banner1URL = newsService.awsS3Service.GetPresignedURL(enums.BannersBucket, news.BannerPath, 8*time.Hour)
	}
	if news.Banner2Path != "" {
		banner2URL = newsService.awsS3Service.GetPresignedURL(enums.BannersBucket, news.Banner2Path, 8*time.Hour)
	}

	author, _ := newsService.userRepository.FindByUserID(news.AuthorID)

	categories := newsService.newsRepository.FindNewsCategoriesByNews(news)
	categoryNames := make([]string, len(categories))
	for i, category := range categories {
		categoryNames[i] = category.Name
	}

	newsDetails := dto.NewsDetailsResponse{
		ID:          newsID,
		Title:       news.Title,
		Description: news.Description,
		Content:     news.Content,
		Content2:    news.Content,
		Banner:      banner1URL,
		Banner2:     banner2URL,
		Categories:  categoryNames,
		Author:      author.Name,
	}

	return newsDetails
}

func (newsService *NewsService) GetNewsList(page, pageSize int) []dto.NewsDetailsResponse {
	offset := (page - 1) * pageSize
	newsList, _ := newsService.newsRepository.FindAllNews(offset, pageSize)

	newsDetails := make([]dto.NewsDetailsResponse, len(newsList))
	for i, news := range newsList {
		banner := ""
		if news.BannerPath != "" {
			banner = newsService.awsS3Service.GetPresignedURL(enums.BannersBucket, news.BannerPath, 8*time.Hour)
		}

		author, _ := newsService.userRepository.FindByUserID(news.AuthorID)

		categories := newsService.newsRepository.FindNewsCategoriesByNews(news)
		categoryNames := make([]string, len(categories))
		for i, category := range categories {
			categoryNames[i] = category.Name
		}

		newsDetails[i] = dto.NewsDetailsResponse{
			ID:          news.ID,
			Title:       news.Title,
			Description: news.Description,
			Banner:      banner,
			Categories:  categoryNames,
			Author:      author.Name,
		}
	}

	return newsDetails
}

func (newsService *NewsService) FilterNewsByCategory(categories []string, page, pageSize int) []dto.NewsDetailsResponse {
	var newsList []*entities.News
	offset := (page - 1) * pageSize
	if len(categories) == 0 {
		newsList, _ = newsService.newsRepository.FindAllNews(offset, pageSize)
	} else {
		newsList = newsService.newsRepository.FindNewsByCategoryName(categories, offset, pageSize)
	}

	newsDetails := make([]dto.NewsDetailsResponse, len(newsList))
	for i, news := range newsList {
		banner := ""
		if news.BannerPath != "" {
			banner = newsService.awsS3Service.GetPresignedURL(enums.BannersBucket, news.BannerPath, 8*time.Hour)
		}
		author, _ := newsService.userRepository.FindByUserID(news.AuthorID)
		categories := make([]string, len(news.Categories))
		for i, category := range news.Categories {
			categories[i] = category.Name
		}
		newsDetails[i] = dto.NewsDetailsResponse{
			ID:          news.ID,
			Title:       news.Title,
			Description: news.Description,
			Banner:      banner,
			Categories:  categories,
			Author:      author.Name,
		}
	}

	return newsDetails

}
