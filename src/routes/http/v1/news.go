package routes_http_v1

import (
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
	newsController := controller_v1_general.NewNewsController(di.Constants, newsService)

	news := routerGroup.Group("/news")
	{
		news.POST("/create", newsController.CreateNews)
		news.GET("/:id", newsController.GetNewsByID)
		news.PUT("/:id", newsController.UpdateNews)
		news.DELETE("/:id", newsController.DeleteNews)
		news.GET("", newsController.GetNewsList)
		news.GET("/topk", newsController.GetTopKNews)
		news.GET("/category/:category", newsController.GetNewsByCategory)
	}
}
