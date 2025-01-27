package routes_http_v1

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"first-project/src/application"
	application_aws "first-project/src/application/aws"
	application_communication "first-project/src/application/communication/emailService"
	application_jwt "first-project/src/application/jwt"
	"first-project/src/bootstrap"
	controller_v1_category "first-project/src/controller/v1/category"
	controller_v1_comment "first-project/src/controller/v1/comment"
	controller_v1_event "first-project/src/controller/v1/event"
	controller_v1_journal "first-project/src/controller/v1/journal"
	controller_v1_news "first-project/src/controller/v1/news"
	controller_v1_podcast "first-project/src/controller/v1/podcast"

	controller_v1_user "first-project/src/controller/v1/user"
	repository_database "first-project/src/repository/database"
	repository_cache "first-project/src/repository/redis"
)

func SetupGeneralRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) {
	userRepository := repository_database.NewUserRepository()
	categoryRepository := repository_database.NewCategoryRepository()
	eventRepository := repository_database.NewEventRepository()
	commentRepository := repository_database.NewCommentRepository(db)
	podcastRepository := repository_database.NewPodcastRepository()
	newsRepository := repository_database.NewNewsRepository()
	journalRepository := repository_database.NewJournalRepository()
	purchaseRepository := repository_database.NewPurchaseRepository()
	userCache := repository_cache.NewUserCache(db, di.Constants, rdb, userRepository)

	jwtService := application_jwt.NewJWTToken()
	emailService := application_communication.NewEmailService(&di.Env.Email)
	otpService := application.NewOTPService()
	awsService := application_aws.NewS3Service(
		di.Constants, &di.Env.EventsBucket,
		&di.Env.PodcastsBucket, &di.Env.NewsBucket,
		&di.Env.JournalsBucket, &di.Env.ProfilesBucket,
	)
	categoryService := application.NewCategoryService(di.Constants, categoryRepository, db)
	eventService := application.NewEventService(di.Constants, awsService, categoryService, eventRepository, commentRepository, purchaseRepository, db)
	userService := application.NewUserService(di.Constants, userRepository, otpService, awsService, db)
	commentService := application.NewCommentService(di.Constants, commentRepository, userService, db)
	podcastService := application.NewPodcastService(di.Constants, awsService, categoryService, podcastRepository, commentRepository, userService, db)
	newsService := application.NewNewsService(di.Constants, awsService, categoryService, commentRepository, newsRepository, userService, db)
	journalService := application.NewJournalService(di.Constants, awsService, userService, journalRepository, db)

	generalCategoryController := controller_v1_category.NewGeneralCategoryController(categoryService)
	generalEventController := controller_v1_event.NewGeneralEventController(di.Constants, eventService, emailService)
	generalCommentController := controller_v1_comment.NewGeneralCommentController(di.Constants, commentService)
	generalPodcastController := controller_v1_podcast.NewGeneralPodcastController(di.Constants, podcastService)
	generalUserController := controller_v1_user.NewGeneralUserController(di.Constants, userService, emailService, userCache, otpService, jwtService)
	generalNewsController := controller_v1_news.NewGeneralNewsController(di.Constants, newsService)
	generalJournalController := controller_v1_journal.NewGeneralJournalController(di.Constants, journalService)

	const (
		searchEndpoint = "/search"
		filterEndpoint = "/filter"
	)

	public := routerGroup.Group("/public")
	{
		categories := public.Group("/categories")
		{
			categories.GET("", generalCategoryController.GetListCategoryNames)
		}

		events := public.Group("/events")
		{
			events.GET("/published", generalEventController.ListEvents)
			events.GET(searchEndpoint, generalEventController.SearchEvents)
			events.GET(filterEndpoint, generalEventController.FilterEvents)
			events.GET("/:eventID", generalEventController.GetEventDetails)
			events.GET("/:eventID/organizers", generalEventController.GetEventOrganizers)
		}

		podcasts := public.Group("/podcasts")
		{
			podcasts.GET("", generalPodcastController.GetPodcastsList)
			podcasts.GET(searchEndpoint, generalPodcastController.SearchPodcast)
			podcasts.GET(filterEndpoint, generalPodcastController.FilterPodcastByCategory)
			podcasts.GET("/:podcastID/episodes", generalPodcastController.GetEpisodesList)
			podcasts.GET("/:podcastID", generalPodcastController.GetPodcastDetails)
		}

		episodes := public.Group("/episodes")
		{
			episodes.GET("/:episodeID", generalPodcastController.GetEpisodeDetails)
		}

		comments := public.Group("/comments/:postID")
		{
			comments.GET("", generalCommentController.GetComments)
		}

		news := public.Group("/news")
		{
			news.GET("", generalNewsController.GetNewsList)
			news.GET("/:newsID", generalNewsController.GetNewsDetails)
			news.GET(searchEndpoint, generalNewsController.SearchNews)
			news.GET(filterEndpoint, generalNewsController.FilterNewsByCategory)
		}

		journals := public.Group("/journals")
		{
			journals.GET("", generalJournalController.GetJournalsList)
			journals.GET(searchEndpoint, generalJournalController.SearchJournals)
		}

		councilors := public.Group("/councilors")
		{
			councilors.GET("", generalUserController.GetCouncilors)
		}
	}

	auth := routerGroup.Group("/auth")
	{
		auth.POST("/register", generalUserController.Register)
		auth.POST("/register/verify", generalUserController.VerifyEmail)
		auth.POST("/login", generalUserController.Login)
		auth.POST("/forgot-password", generalUserController.ForgotPassword)
		auth.POST("/confirm-otp", generalUserController.ConfirmOTP)
		auth.POST("/refresh-token", generalUserController.RefreshToken)
	}
}
