package routes_http_v1

import (
	"first-project/src/application"
	application_aws "first-project/src/application/aws"
	application_communication "first-project/src/application/communication/emailService"
	application_jwt "first-project/src/application/jwt"
	"first-project/src/bootstrap"
	controller_v1_event "first-project/src/controller/v1/event"
	controller_v1_private "first-project/src/controller/v1/private"
	middleware_authentication "first-project/src/middleware/Authentication"
	repository_database "first-project/src/repository/database"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupCustomerRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) {
	userRepository := repository_database.NewUserRepository(db)
	eventRepository := repository_database.NewEventRepository(db)
	commentRepository := repository_database.NewCommentRepository(db)
	podcastRepository := repository_database.NewPodcastRepository(db)

	jwtService := application_jwt.NewJWTToken()
	emailService := application_communication.NewEmailService(&di.Env.Email)
	awsService := application_aws.NewS3Service(di.Constants, &di.Env.BannersBucket, &di.Env.SessionsBucket, &di.Env.PodcastsBucket, &di.Env.ProfileBucket)
	eventService := application.NewEventService(di.Constants, awsService, eventRepository, commentRepository)
	commentService := application.NewCommentService(di.Constants, commentRepository, userRepository)
	podcastService := application.NewPodcastService(di.Constants, awsService, podcastRepository, commentRepository, userRepository)

	authMiddleware := middleware_authentication.NewAuthMiddleware(di.Constants, userRepository, jwtService)

	eventController := controller_v1_event.NewEventController(di.Constants, eventService, emailService)
	commentController := controller_v1_private.NewCommentController(di.Constants, commentService)
	podcastController := controller_v1_private.NewPodcastController(di.Constants, podcastService)

	event := routerGroup.Group("/events/:eventID")
	{
		event.GET("", authMiddleware.OptionalAuth, eventController.GetPublicEventDetails)
		event.GET("/tickets", authMiddleware.AuthRequired, eventController.GetAvailableEventTicketsList)
	}

	comments := routerGroup.Group("/comments")
	comments.Use(authMiddleware.AuthRequired)
	{
		comments.POST("/post/:postID", commentController.CreateComment)

		commentSubGroup := comments.Group("/:commentID")
		{
			commentSubGroup.PUT("", commentController.EditComment)
			commentSubGroup.DELETE("", commentController.DeleteCommentByUser)
		}
	}

	podcast := routerGroup.Group("/podcasts/:podcastID")
	{
		podcast.GET("", authMiddleware.OptionalAuth, podcastController.GetPodcastDetails)

		podcastSubscription := podcast.Group("/subscribe")
		podcastSubscription.Use(authMiddleware.AuthRequired)
		{
			podcastSubscription.POST("", podcastController.SubscribePodcast)
			podcastSubscription.DELETE("", podcastController.UnSubscribePodcast)
		}
	}

	profile := routerGroup.Group("/profile")
	profile.Use(authMiddleware.AuthRequired)
	{
		profile.GET("") // some sample api here ...
	}
}
