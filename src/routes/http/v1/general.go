package routes_http_v1

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"first-project/src/application"
	application_aws "first-project/src/application/aws"
	application_communication "first-project/src/application/communication/emailService"
	application_jwt "first-project/src/application/jwt"
	application_news "first-project/src/application/news"
	"first-project/src/bootstrap"
	controller_v1_event "first-project/src/controller/v1/event"
	controller_v1_general "first-project/src/controller/v1/general"
	controller_v1_journal "first-project/src/controller/v1/journal"
	controller_v1_news "first-project/src/controller/v1/news"
	controller_v1_private "first-project/src/controller/v1/private"
	repository_database "first-project/src/repository/database"
	repository_cache "first-project/src/repository/redis"
)

func SetupGeneralRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) {
	userRepository := repository_database.NewUserRepository(db)
	eventRepository := repository_database.NewEventRepository(db)
	commentRepository := repository_database.NewCommentRepository(db)
	podcastRepository := repository_database.NewPodcastRepository(db)
	newsRepository := repository_database.NewNewsRepository(db)
	journalRepository := repository_database.NewJournalRepository(db)
	userCache := repository_cache.NewUserCache(di.Constants, rdb, userRepository)

	jwtService := application_jwt.NewJWTToken()
	emailService := application_communication.NewEmailService(&di.Env.Email)
	otpService := application.NewOTPService()
	awsService := application_aws.NewS3Service(di.Constants, &di.Env.BannersBucket, &di.Env.SessionsBucket, &di.Env.PodcastsBucket, &di.Env.ProfileBucket)
	eventService := application.NewEventService(di.Constants, awsService, eventRepository, commentRepository)
	commentService := application.NewCommentService(di.Constants, commentRepository, userRepository)
	podcastService := application.NewPodcastService(di.Constants, awsService, podcastRepository, commentRepository, userRepository)
	userService := application.NewUserService(di.Constants, userRepository, otpService)
	newsService := application_news.NewNewsService(di.Constants, awsService, commentRepository, newsRepository, userRepository)
	journalService := application.NewJournalService(di.Constants, awsService, userRepository, journalRepository)

	eventController := controller_v1_event.NewEventController(di.Constants, eventService, emailService)
	commentController := controller_v1_private.NewCommentController(di.Constants, commentService)
	authController := controller_v1_general.NewAuthController(di.Constants, jwtService)
	podcastController := controller_v1_private.NewPodcastController(di.Constants, podcastService)
	userController := controller_v1_general.NewUserController(di.Constants, userService, emailService, userCache, otpService, jwtService)
	newsController := controller_v1_news.NewNewsController(di.Constants, newsService)
	journalController := controller_v1_journal.NewJournalController(di.Constants, journalService)

	const (
		searchEndpoint = "/search"
		filterEndpoint = "/filter"
	)

	public := routerGroup.Group("/public")
	{
		public.GET("/categories", eventController.ListCategories)

		events := public.Group("/events")
		{
			events.GET("/published", eventController.ListPublicEvents)
			events.GET(searchEndpoint, eventController.SearchPublicEvents)
			events.GET(filterEndpoint, eventController.FilterPublicEvents)

			eventSubGroup := events.Group("/:eventID")
			{
				eventSubGroup.GET("", eventController.GetPublicEventDetails)
				eventSubGroup.GET("/tickets", eventController.GetAvailableTicketDetails)
			}
		}

		podcasts := public.Group("/podcasts")
		{
			podcasts.GET("", podcastController.GetPodcastsList)
			podcasts.GET("/:podcastID", podcastController.GetPodcastDetails)
			podcasts.GET("/:podcastID/episodes", podcastController.GetEpisodesList)
			podcasts.GET(searchEndpoint, podcastController.SearchPodcast)
			podcasts.GET(filterEndpoint, podcastController.FilterPodcastByCategory)
		}
		episodes := public.Group("/episodes")
		{
			episodes.GET("/:episodeID", podcastController.GetEpisodeDetails)
		}
		public.GET("comments/:postID", commentController.GetComments)

		news := public.Group("/news")
		{
			news.GET("", newsController.GetNewsList)
			news.GET("/:newsID", newsController.GetNewsDetails)
			news.GET(searchEndpoint, newsController.SearchNews)
			news.GET(filterEndpoint, newsController.FilterNewsByCategory)
		}

		journals := public.Group("/journals")
		{
			journals.GET("", journalController.GetJournalsList)
			journals.GET(searchEndpoint, journalController.SearchJournals)
		}
	}

	auth := routerGroup.Group("/auth")
	{
		auth.POST("/register", userController.Register)
		auth.POST("/register/verify", userController.VerifyEmail)
		auth.POST("/login", userController.Login)
		auth.POST("/forgot-password", userController.ForgotPassword)
		auth.POST("/confirm-otp", userController.ConfirmOTP)
		auth.PUT("/reset-password", userController.ResetPassword)
		auth.POST("/refresh-token", authController.RefreshToken)
	}
}
