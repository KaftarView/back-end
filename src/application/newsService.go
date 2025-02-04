package application

import (
	application_aws "first-project/src/application/aws"
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"
	"first-project/src/dto"
	"first-project/src/entities"
	"first-project/src/enums"
	"first-project/src/exceptions"
	repository_database "first-project/src/repository/database"
	repository_database_interfaces "first-project/src/repository/database/interfaces"
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

type NewsService struct {
	constants         *bootstrap.Constants
	awsS3Service      *application_aws.S3Service
	categoryService   application_interfaces.CategoryService
	commentRepository repository_database_interfaces.CommentRepository
	newsRepository    repository_database_interfaces.NewsRepository
	userService       application_interfaces.UserService
	db                *gorm.DB
}

func NewNewsService(
	constants *bootstrap.Constants,
	awsS3Service *application_aws.S3Service,
	categoryService application_interfaces.CategoryService,
	commentRepository repository_database_interfaces.CommentRepository,
	newsRepository repository_database_interfaces.NewsRepository,
	userService application_interfaces.UserService,
	db *gorm.DB,
) *NewsService {
	return &NewsService{
		constants:         constants,
		awsS3Service:      awsS3Service,
		categoryService:   categoryService,
		commentRepository: commentRepository,
		newsRepository:    newsRepository,
		userService:       userService,
		db:                db,
	}
}

func (newsService *NewsService) fetchNewsByID(newsID uint) *entities.News {
	var notFoundError exceptions.NotFoundError
	news, newsExist := newsService.newsRepository.FindNewsByID(newsService.db, newsID)
	if !newsExist {
		notFoundError.ErrorField = newsService.constants.ErrorField.News
		panic(notFoundError)
	}
	return news
}

func (newsService *NewsService) validateUniqueNewsTittle(tittle string) {
	var conflictError exceptions.ConflictError
	_, newsExist := newsService.newsRepository.FindNewsByTitle(newsService.db, tittle)
	if newsExist {
		conflictError.AppendError(
			newsService.constants.ErrorField.Tittle,
			newsService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}
}

func (newsService *NewsService) getNewsBannerPath(news *entities.News, banner *multipart.FileHeader) string {
	bannerPath := newsService.constants.S3Service.GetNewsBannerKey(news.ID, banner.Filename)
	newsService.awsS3Service.UploadObject(enums.NewsBucket, bannerPath, banner)
	return bannerPath
}

func (newsService *NewsService) updateBannerPath(news *entities.News, currentPath *string, newBanner *multipart.FileHeader) {
	if *currentPath != "" {
		newsService.awsS3Service.DeleteObject(enums.NewsBucket, *currentPath)
	}
	*currentPath = newsService.getNewsBannerPath(news, newBanner)
}

func (newsService *NewsService) CreateNews(newsDetails dto.CreateNewsRequest) *entities.News {
	newsService.validateUniqueNewsTittle(newsDetails.Title)

	categories := newsService.categoryService.GetCategoriesByName(newsDetails.Categories)

	var news *entities.News
	err := repository_database.ExecuteInTransaction(newsService.db, func(tx *gorm.DB) error {
		commentable := newsService.commentRepository.CreateNewCommentable(tx)

		bannerPath := newsService.constants.S3Service.GetNewsBannerKey(commentable.CID, newsDetails.Banner.Filename)
		newsService.awsS3Service.UploadObject(enums.NewsBucket, bannerPath, newsDetails.Banner)

		banner2Path := ""
		if newsDetails.Banner2 != nil {
			banner2Path = newsService.constants.S3Service.GetNewsBannerKey(commentable.CID, newsDetails.Banner2.Filename)
			newsService.awsS3Service.UploadObject(enums.NewsBucket, banner2Path, newsDetails.Banner2)
		}

		news = &entities.News{
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
		if err := newsService.newsRepository.CreateNews(tx, news); err != nil {
			panic(err)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return news
}

func (newsService *NewsService) UpdateNews(newsDetails dto.UpdateNewsRequest) {
	news := newsService.fetchNewsByID(newsDetails.ID)

	err := repository_database.ExecuteInTransaction(newsService.db, func(tx *gorm.DB) error {
		if newsDetails.Title != nil {
			newsService.validateUniqueNewsTittle(*newsDetails.Title)
			news.Title = *newsDetails.Title
		}
		updateField(newsDetails.Description, &news.Description)
		updateField(newsDetails.Content, &news.Content)
		updateField(newsDetails.Content2, &news.Content2)

		if newsDetails.Banner != nil {
			newsService.updateBannerPath(news, &news.BannerPath, newsDetails.Banner)
		}
		if newsDetails.Banner2 != nil {
			newsService.updateBannerPath(news, &news.Banner2Path, newsDetails.Banner2)
		}
		if newsDetails.Categories != nil {
			categories := newsService.categoryService.GetCategoriesByName(*newsDetails.Categories)
			if err := newsService.newsRepository.UpdateNewsCategories(tx, newsDetails.ID, categories); err != nil {
				panic(err)
			}
		}
		if err := newsService.newsRepository.UpdateNews(tx, news); err != nil {
			panic(err)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (newsService *NewsService) DeleteNews(newsID uint) {
	news := newsService.fetchNewsByID(newsID)

	err := repository_database.ExecuteInTransaction(newsService.db, func(tx *gorm.DB) error {
		if err := newsService.newsRepository.DeleteNews(tx, newsID); err != nil {
			panic(err)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	if news.BannerPath != "" {
		newsService.awsS3Service.DeleteObject(enums.NewsBucket, news.BannerPath)
	}
	if news.Banner2Path != "" {
		newsService.awsS3Service.DeleteObject(enums.NewsBucket, news.Banner2Path)
	}
}

func (newsService *NewsService) GetNewsDetails(newsID uint) dto.NewsDetailsResponse {
	news := newsService.fetchNewsByID(newsID)

	banner1URL := ""
	banner2URL := ""
	if news.BannerPath != "" {
		banner1URL = newsService.awsS3Service.GetPresignedURL(enums.NewsBucket, news.BannerPath, 8*time.Hour)
	}
	if news.Banner2Path != "" {
		banner2URL = newsService.awsS3Service.GetPresignedURL(enums.NewsBucket, news.Banner2Path, 8*time.Hour)
	}

	author, _ := newsService.userService.FindByUserID(news.AuthorID)

	categories := newsService.newsRepository.FindNewsCategoriesByNews(newsService.db, news)
	categoryNames := make([]string, len(categories))
	for i, category := range categories {
		categoryNames[i] = category.Name
	}

	newsDetails := dto.NewsDetailsResponse{
		ID:          newsID,
		Title:       news.Title,
		CreatedAt:   news.CreatedAt,
		Description: news.Description,
		Content:     news.Content,
		Content2:    news.Content2,
		Banner:      banner1URL,
		Banner2:     banner2URL,
		Categories:  categoryNames,
		Author:      author.Name,
	}

	return newsDetails
}

func (newsService *NewsService) GetNewsList(page, pageSize int) []dto.NewsDetailsResponse {
	offset := (page - 1) * pageSize
	newsList, _ := newsService.newsRepository.FindAllNews(newsService.db, offset, pageSize)

	newsDetails := make([]dto.NewsDetailsResponse, len(newsList))
	for i, news := range newsList {
		banner := ""
		if news.BannerPath != "" {
			banner = newsService.awsS3Service.GetPresignedURL(enums.NewsBucket, news.BannerPath, 8*time.Hour)
		}

		author, _ := newsService.userService.FindByUserID(news.AuthorID)

		categories := newsService.newsRepository.FindNewsCategoriesByNews(newsService.db, news)
		categoryNames := make([]string, len(categories))
		for i, category := range categories {
			categoryNames[i] = category.Name
		}

		newsDetails[i] = dto.NewsDetailsResponse{
			ID:          news.ID,
			Title:       news.Title,
			CreatedAt:   news.CreatedAt,
			Description: news.Description,
			Banner:      banner,
			Categories:  categoryNames,
			Author:      author.Name,
		}
	}

	return newsDetails
}

func (newsService *NewsService) SearchNews(query string, page, pageSize int) []dto.NewsDetailsResponse {
	var newsList []*entities.News
	offset := (page - 1) * pageSize
	if query != "" {
		newsList = newsService.newsRepository.FullTextSearch(newsService.db, query, offset, pageSize)
	} else {
		newsList, _ = newsService.newsRepository.FindAllNews(newsService.db, offset, pageSize)
	}

	newsDetails := make([]dto.NewsDetailsResponse, len(newsList))
	for i, news := range newsList {
		banner := ""
		if news.BannerPath != "" {
			banner = newsService.awsS3Service.GetPresignedURL(enums.NewsBucket, news.BannerPath, 8*time.Hour)
		}

		author, _ := newsService.userService.FindByUserID(news.AuthorID)

		categories := newsService.newsRepository.FindNewsCategoriesByNews(newsService.db, news)
		categoryNames := make([]string, len(categories))
		for i, category := range categories {
			categoryNames[i] = category.Name
		}

		newsDetails[i] = dto.NewsDetailsResponse{
			ID:          news.ID,
			Title:       news.Title,
			CreatedAt:   news.CreatedAt,
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
		newsList, _ = newsService.newsRepository.FindAllNews(newsService.db, offset, pageSize)
	} else {
		newsList = newsService.newsRepository.FindNewsByCategoryName(newsService.db, categories, offset, pageSize)
	}

	newsDetails := make([]dto.NewsDetailsResponse, len(newsList))
	for i, news := range newsList {
		banner := ""
		if news.BannerPath != "" {
			banner = newsService.awsS3Service.GetPresignedURL(enums.NewsBucket, news.BannerPath, 30*time.Minute)
		}
		author, _ := newsService.userService.FindByUserID(news.AuthorID)
		categories := make([]string, len(news.Categories))
		for i, category := range news.Categories {
			categories[i] = category.Name
		}
		newsDetails[i] = dto.NewsDetailsResponse{
			ID:          news.ID,
			Title:       news.Title,
			CreatedAt:   news.CreatedAt,
			Description: news.Description,
			Banner:      banner,
			Categories:  categories,
			Author:      author.Name,
		}
	}

	return newsDetails
}
