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
	awsService := application_aws.NewS3Service(di.Constants, &di.Env.BannersBucket, &di.Env.SessionsBucket, &di.Env.PodcastsBucket, &di.Env.ProfileBucket)
	podcastService := application.NewPodcastService(di.Constants, awsService, podcastRepository, commentRepository, userRepository)
	podcastController := controller_v1_private.NewPodcastController(di.Constants, podcastService)

	podcasts := routerGroup.Group("/podcasts")

	podcasts.POST("/:podcastID/subscribe", podcastController.SubscribePodcast)
	podcasts.DELETE("/:podcastID/subscribe", podcastController.UnSubscribePodcast)

	podcasts.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ManagePodcasts}))
	{
		podcasts.GET("", podcastController.GetPodcastsList)
		podcasts.POST("", podcastController.CreatePodcast)

		podcastSubRouter := podcasts.Group("/:podcastID")
		{
			podcastSubRouter.GET("", podcastController.GetPodcastDetails)
			podcastSubRouter.PUT("", podcastController.UpdatePodcast)
			podcastSubRouter.DELETE("", podcastController.DeletePodcast)

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
