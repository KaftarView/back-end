package routes_http_v1

import (
	"first-project/src/application"
	application_aws "first-project/src/application/aws"
	application_jwt "first-project/src/application/jwt"
	"first-project/src/bootstrap"
	controller_v1_private "first-project/src/controller/v1/private"
	"first-project/src/enums"
	middleware_authentication "first-project/src/middleware/Authentication"
	repository_database "first-project/src/repository/database"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupPodcastRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) {
	userRepository := repository_database.NewUserRepository(db)
	podcastRepository := repository_database.NewPodcastRepository(db)
	commentRepository := repository_database.NewCommentRepository(db)
	jwtService := application_jwt.NewJWTToken()
	authMiddleware := middleware_authentication.NewAuthMiddleware(di.Constants, userRepository, jwtService)
	podcastService := application.NewPodcastService(di.Constants, podcastRepository, commentRepository)
	awsService := application_aws.NewS3Service(di.Constants, &di.Env.BannersBucket, &di.Env.SessionsBucket, &di.Env.PodcastsBucket, &di.Env.ProfileBucket)
	podcastController := controller_v1_private.NewPodcastController(di.Constants, podcastService, awsService)

	podcasts := routerGroup.Group("/podcasts")
	podcasts.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ManagePodcasts}))
	{
		podcasts.GET("", podcastController.GetPodcastsList)
		podcasts.POST("", podcastController.CreatePodcast)

		podcastSubRouter := podcasts.Group("/:podcastID")
		{
			podcasts.GET("", podcastController.GetPodcastDetails)
			podcasts.PUT("", podcastController.UpdatePodcast)
			podcasts.DELETE("", podcastController.DeletePodcast)
			podcasts.POST("/subscribe", podcastController.SubscribePodcast)
			podcasts.DELETE("/subscribe", podcastController.UnSubscribePodcast)

			podcastEpisodesSubRouter := podcastSubRouter.Group("/episodes")
			{
				podcastEpisodesSubRouter.GET("", podcastController.GetEpisodesList)
				podcastEpisodesSubRouter.POST("", podcastController.CreateEpisode)
			}
		}

	}

	episodes := routerGroup.Group("/episodes")
	episodes.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ManagePodcasts}))
	{
		episodes.PUT("/:episodeID", podcastController.UpdateEpisode)
		episodes.DELETE("/:episodeID", podcastController.DeleteEpisode)
	}
}
