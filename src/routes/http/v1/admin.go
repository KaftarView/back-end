package routes_http_v1

import (
	"first-project/src/application"
	application_aws "first-project/src/application/aws"
	application_communication "first-project/src/application/communication/emailService"
	application_jwt "first-project/src/application/jwt"
	"first-project/src/bootstrap"
	controller_v1_comment "first-project/src/controller/v1/comment"
	controller_v1_event "first-project/src/controller/v1/event"
	controller_v1_journal "first-project/src/controller/v1/journal"
	controller_v1_news "first-project/src/controller/v1/news"
	controller_v1_podcast "first-project/src/controller/v1/podcast"
	controller_v1_user "first-project/src/controller/v1/user"
	"first-project/src/enums"
	middleware_authentication "first-project/src/middleware/Authentication"
	repository_database "first-project/src/repository/database"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupAdminRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) {
	userRepository := repository_database.NewUserRepository()
	categoryRepository := repository_database.NewCategoryRepository()
	eventRepository := repository_database.NewEventRepository()
	commentRepository := repository_database.NewCommentRepository(db)
	podcastRepository := repository_database.NewPodcastRepository()
	newsRepository := repository_database.NewNewsRepository()
	journalRepository := repository_database.NewJournalRepository()
	purchaseRepository := repository_database.NewPurchaseRepository()

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
	commentService := application.NewCommentService(di.Constants, commentRepository, userRepository, db)
	podcastService := application.NewPodcastService(di.Constants, awsService, categoryService, podcastRepository, commentRepository, userRepository, db)
	userService := application.NewUserService(di.Constants, userRepository, otpService, awsService, db)
	newsService := application.NewNewsService(di.Constants, awsService, categoryService, commentRepository, newsRepository, userRepository, db)
	journalService := application.NewJournalService(di.Constants, awsService, userRepository, journalRepository, db)

	authMiddleware := middleware_authentication.NewAuthMiddleware(di.Constants, userRepository, jwtService, db)

	adminEventController := controller_v1_event.NewAdminEventController(di.Constants, eventService, emailService)
	adminCommentController := controller_v1_comment.NewAdminCommentController(di.Constants, commentService)
	adminPodcastController := controller_v1_podcast.NewAdminPodcastController(di.Constants, podcastService)
	adminUserController := controller_v1_user.NewAdminUserController(di.Constants, userService)
	adminNewsController := controller_v1_news.NewAdminNewsController(di.Constants, newsService)
	adminJournalController := controller_v1_journal.NewAdminJournalController(di.Constants, journalService)

	events := routerGroup.Group("/events")
	{
		readGroup := events.Group("")
		readGroup.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.CreateEvent, enums.ManageEvent, enums.ReviewEvent}))
		{
			readGroup.GET("", adminEventController.GetEventsList)
			readGroup.GET("/search", adminEventController.SearchEvents)
			readGroup.GET("/filter", adminEventController.FilterEvents)

			readGroup.GET("ticket/:ticketID", adminEventController.GetTicketDetails)
			readGroup.GET("discount/:discountID", adminEventController.GetDiscountDetails)
			readGroup.GET("media/:mediaID", adminEventController.GetMediaDetails)

			readSingleEventGroup := readGroup.Group("/:eventID")
			{
				readSingleEventGroup.GET("", adminEventController.GetEventDetails)
				readSingleEventGroup.GET("/tickets", adminEventController.GetAllTicketDetails)
				readSingleEventGroup.GET("/discounts", adminEventController.GetAllDiscountDetails)
				readSingleEventGroup.GET("/media", adminEventController.GetEventMedia)
				readSingleEventGroup.GET("/attendees", adminEventController.GetEventAttendees)
			}
		}

		createGroup := events.Group("")
		createGroup.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.CreateEvent}))
		{
			createGroup.POST("/create", adminEventController.CreateEvent)
			createGroup.POST("/add-ticket/:eventID", adminEventController.AddEventTicket)
			createGroup.POST("/add-discount/:eventID", adminEventController.AddEventDiscount)
			createGroup.POST("/add-organizer/:eventID", adminEventController.AddEventOrganizer)
		}

		manageEventsGroup := events.Group("")
		manageEventsGroup.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.CreateEvent, enums.ManageEvent}))
		{
			eventSubGroup := events.Group("/:eventID")
			{
				eventSubGroup.PUT("", adminEventController.UpdateEvent)
				eventSubGroup.POST("/publish", adminEventController.PublishEvent)
				eventSubGroup.POST("/unpublish", adminEventController.UnpublishEvent)
				eventSubGroup.DELETE("", adminEventController.DeleteEvent)
				eventSubGroup.POST("/media", adminEventController.UploadEventMedia)
			}

			ticketSubGroup := manageEventsGroup.Group("/ticket/:ticketID")
			{
				ticketSubGroup.PUT("", adminEventController.UpdateEventTicket)
				ticketSubGroup.DELETE("", adminEventController.DeleteTicket)
			}

			discountSubGroup := manageEventsGroup.Group("/discount/:discountID")
			{
				discountSubGroup.PUT("", adminEventController.UpdateEventDiscount)
				discountSubGroup.DELETE("", adminEventController.DeleteDiscount)
			}

			mediaSubGroup := manageEventsGroup.Group("/media/:mediaID")
			{
				mediaSubGroup.PUT("", adminEventController.UpdateEventMedia)
				mediaSubGroup.DELETE("", adminEventController.DeleteEventMedia)
			}

			manageEventsGroup.DELETE("/organizer/:organizerID", adminEventController.DeleteOrganizer)
		}
	}

	comments := routerGroup.Group("/comments")
	comments.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ModerateComments}))
	{
		comments.DELETE("/:commentID", adminCommentController.DeleteComment)
	}

	podcasts := routerGroup.Group("/podcasts")
	podcasts.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ManagePodcasts}))
	{
		podcasts.POST("", adminPodcastController.CreatePodcast)
		podcastSubGroup := podcasts.Group("/:podcastID")
		{
			podcastSubGroup.PUT("", adminPodcastController.UpdatePodcast)
			podcastSubGroup.DELETE("", adminPodcastController.DeletePodcast)

			podcastEpisodesSubRouter := podcastSubGroup.Group("/episodes")
			{
				podcastEpisodesSubRouter.POST("", adminPodcastController.CreateEpisode)
			}
		}
	}

	podcastEpisodes := routerGroup.Group("/episodes")
	podcastEpisodes.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ManagePodcasts}))
	{
		podcastEpisodes.PUT("/:episodeID", adminPodcastController.UpdateEpisode)
		podcastEpisodes.DELETE("/:episodeID", adminPodcastController.DeleteEpisode)
	}

	accessManagement := routerGroup.Group("")
	accessManagement.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ManageUsers, enums.ManageRoles}))
	{
		roles := accessManagement.Group("/roles")
		{
			roles.GET("", adminUserController.GetRolesList)
			roles.POST("", adminUserController.CreateRole)

			roleSubGroup := roles.Group("/:roleID")
			{
				roleSubGroup.GET("/owners", adminUserController.GetRoleOwners)
				roleSubGroup.DELETE("", adminUserController.DeleteRole)

				rolePermissions := roleSubGroup.Group("/permissions")
				{
					rolePermissions.PUT("", adminUserController.UpdateRole)
					rolePermissions.DELETE("/:permissionID", adminUserController.DeleteRolePermission)
				}
			}
		}

		accessManagement.GET("/permissions", adminUserController.GetPermissionsList)

		userRoles := accessManagement.Group("/users/roles")
		{
			userRoles.PUT("", adminUserController.UpdateUserRoles)
			userRoles.DELETE("/:roleID", adminUserController.DeleteUserRole)
		}

		councilors := accessManagement.Group("/councilors")
		{
			councilors.POST("", adminUserController.CreateCouncilor)
			councilors.DELETE("/:councilorID", adminUserController.DeleteCouncilor)
		}
	}

	news := routerGroup.Group("/news")
	news.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ManageNews}))
	{
		news.POST("", adminNewsController.CreateNews)
		newsSubGroup := news.Group("/:newsID")
		{
			newsSubGroup.PUT("", adminNewsController.UpdateNews)
			newsSubGroup.DELETE("", adminNewsController.DeleteNews)
		}
	}

	journals := routerGroup.Group("/journal")
	journals.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ManageJournal}))
	{
		journals.POST("", adminJournalController.CreateJournal)
		journalSubGroup := journals.Group("/:journalID")
		{
			journalSubGroup.PUT("", adminJournalController.UpdateJournal)
			journalSubGroup.DELETE("", adminJournalController.DeleteJournal)
		}
	}
}
