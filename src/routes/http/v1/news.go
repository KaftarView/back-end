package routes_http_v1

import (
	application_aws "first-project/src/application/aws"
	application_jwt "first-project/src/application/jwt"
	application_news "first-project/src/application/news"
	"first-project/src/bootstrap"
	controller_v1_general "first-project/src/controller/v1/general"
	enums "first-project/src/enums"
	middleware_authentication "first-project/src/middleware/Authentication"
	database "first-project/src/repository/database"
	repository_database "first-project/src/repository/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupNewsRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB) {
	userRepository := repository_database.NewUserRepository(db)
	newsRepository := database.NewNewsRepository(db)
	newsService := application_news.NewNewsService(newsRepository)
	awsService := application_aws.NewS3Service(di.Constants, &di.Env.BannersBucket, &di.Env.SessionsBucket, &di.Env.PodcastsBucket, &di.Env.ProfileBucket)
	newsController := controller_v1_general.NewNewsController(di.Constants, newsService, awsService)
	jwtService := application_jwt.NewJWTToken()
	authMiddleware := middleware_authentication.NewAuthMiddleware(di.Constants, userRepository, jwtService)
	news := routerGroup.Group("/news").Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ManageNewsAndBlogs}))
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
