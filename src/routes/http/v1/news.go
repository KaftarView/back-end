package routes_http_v1

import (
	application_aws "first-project/src/application/aws"
	application_news "first-project/src/application/news"
	"first-project/src/bootstrap"
	controller_v1_general "first-project/src/controller/v1/general"
	database "first-project/src/repository/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupNewsRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB) {
	newsRepository := database.NewNewsRepository(db)
	newsService := application_news.NewNewsService(newsRepository)
	awsService := application_aws.NewS3Service(di.Constants, &di.Env.BannersBucket, &di.Env.SessionsBucket, &di.Env.PodcastsBucket)
	newsController := controller_v1_general.NewNewsController(di.Constants, newsService, awsService)

	news := routerGroup.Group("/news")
	{
		news.POST("/create", newsController.CreateNews)
		news.GET("/:id", newsController.GetNewsByID)
		news.PUT("/:id", newsController.UpdateNews)
		news.DELETE("/:id", newsController.DeleteNews)
		news.GET("", newsController.GetNewsList)
		news.GET("/topk", newsController.GetTopKNews)
		news.POST("/filtered", newsController.GetNewsByCategory)
	}
}
