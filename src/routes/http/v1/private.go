package routes_http_v1

import (
	"first-project/src/application"
	application_aws "first-project/src/application/aws"
	application_communication "first-project/src/application/communication/emailService"
	application_jwt "first-project/src/application/jwt"
	application_news "first-project/src/application/news"
	"first-project/src/bootstrap"
	controller_v1_event "first-project/src/controller/v1/event"
	controller_v1_journal "first-project/src/controller/v1/journal"
	controller_v1_news "first-project/src/controller/v1/news"
	controller_v1_private "first-project/src/controller/v1/private"
	"first-project/src/enums"
	middleware_authentication "first-project/src/middleware/Authentication"
	repository_database "first-project/src/repository/database"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupPrivateRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) {
	userRepository := repository_database.NewUserRepository(db)
	eventRepository := repository_database.NewEventRepository(db)
	commentRepository := repository_database.NewCommentRepository(db)
	podcastRepository := repository_database.NewPodcastRepository(db)
	newsRepository := repository_database.NewNewsRepository(db)
	journalRepository := repository_database.NewJournalRepository(db)

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

	authMiddleware := middleware_authentication.NewAuthMiddleware(di.Constants, userRepository, jwtService)

	eventController := controller_v1_event.NewEventController(di.Constants, eventService, emailService)
	commentController := controller_v1_private.NewCommentController(di.Constants, commentService)
	podcastController := controller_v1_private.NewPodcastController(di.Constants, podcastService)
	roleController := controller_v1_private.NewRoleController(di.Constants, userService)
	newsController := controller_v1_news.NewNewsController(di.Constants, newsService)
	journalController := controller_v1_journal.NewJournalController(di.Constants, journalService)

	events := routerGroup.Group("/events")
	{
		readGroup := events.Group("")
		readGroup.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ManageEvent}))
		{
			readGroup.GET("", eventController.GetEventsListForAdmin)
			readGroup.GET("/event-details/:eventID", eventController.GetEventDetailsForAdmin)
			readGroup.GET("/ticket-details/:eventID", eventController.GetAllTicketDetails)
			readGroup.GET("/discount-details/:eventID", eventController.GetAllDiscountDetails)
			readGroup.GET("/search", eventController.SearchEventsForAdmin)
			readGroup.GET("/filter", eventController.FilterEventsForAdmin)

			readGroup.GET("ticket/:ticketID", eventController.GetTicketDetails)
			readGroup.GET("discount/:discountID", eventController.GetDiscountDetails)
		}

		createGroup := events.Group("")
		createGroup.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.CreateEvent}))
		{
			createGroup.POST("/create", eventController.CreateEvent)
			createGroup.POST("/add-ticket/:eventID", eventController.AddEventTicket)
			createGroup.POST("/add-discount/:eventID", eventController.AddEventDiscount)
			createGroup.POST("/add-organizer/:eventID", eventController.AddEventOrganizer)
		}

		manageEventsGroup := events.Group("")
		manageEventsGroup.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.CreateEvent, enums.EditEvent}))
		{
			eventSubGroup := events.Group("/:eventID")
			{
				eventSubGroup.PUT("", eventController.UpdateEvent)
				eventSubGroup.DELETE("", eventController.DeleteEvent)
				eventSubGroup.POST("/media", eventController.UploadEventMedia)
			}

			ticketSubGroup := manageEventsGroup.Group("/ticket/:ticketID")
			{
				ticketSubGroup.PUT("", eventController.UpdateEventTicket)
				ticketSubGroup.DELETE("", eventController.DeleteTicket)
			}

			discountSubGroup := manageEventsGroup.Group("/discount/:discountID")
			{
				discountSubGroup.PUT("", eventController.UpdateEventDiscount)
				discountSubGroup.DELETE("", eventController.DeleteDiscount)
			}

			manageEventsGroup.DELETE("/organizer/:organizerID", eventController.DeleteOrganizer)
			manageEventsGroup.DELETE("/media/:mediaId", eventController.DeleteEventMedia)
		}

		publishGroup := events.Group("/:eventID")
		publishGroup.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.PublishEvent}))
		{
			publishGroup.POST("/publish", eventController.PublishEvent)
			publishGroup.POST("/unpublish", eventController.UnpublishEvent)
		}
	}

	comments := routerGroup.Group("/comments")
	{
		comments.POST("/post/:postID", commentController.CreateComment)

		commentSubGroup := comments.Group("/:commentID")
		{
			commentSubGroup.PUT("", commentController.EditComment)
			commentSubGroup.DELETE("", commentController.DeleteCommentByUser)
		}

		moderateCommentsGroup := comments.Group("/admin")
		moderateCommentsGroup.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ModerateComments}))
		{
			moderateCommentsGroup.DELETE("/:commentID", commentController.DeleteCommentByAdmin)
		}
	}

	podcasts := routerGroup.Group("/podcasts")
	podcasts.POST("/:podcastID/subscribe", podcastController.SubscribePodcast)
	podcasts.DELETE("/:podcastID/subscribe", podcastController.UnSubscribePodcast)
	podcasts.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ManagePodcasts}))
	{
		podcasts.POST("", podcastController.CreatePodcast)
		podcastSubGroup := podcasts.Group("/:podcastID")
		{
			podcastSubGroup.PUT("", podcastController.UpdatePodcast)
			podcastSubGroup.DELETE("", podcastController.DeletePodcast)

			podcastEpisodesSubRouter := podcastSubGroup.Group("/episodes")
			{
				podcastEpisodesSubRouter.POST("", podcastController.CreateEpisode)
			}
		}
	}

	podcastEpisodes := routerGroup.Group("/episodes")
	podcastEpisodes.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ManagePodcasts}))
	{
		podcastEpisodes.PUT("/:episodeID", podcastController.UpdateEpisode)
		podcastEpisodes.DELETE("/:episodeID", podcastController.DeleteEpisode)
	}

	profile := routerGroup.Group("/profile")
	{
		profile.GET("") // some sample api here ...
	}

	users := routerGroup.Group("")
	users.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ManageUsers, enums.ManageRoles}))
	{
		roles := users.Group("/roles")
		{
			roles.GET("", roleController.GetRolesList)
			roles.POST("", roleController.CreateRole)

			roleSubGroup := roles.Group("/:roleID")
			{
				roleSubGroup.GET("", roleController.GetRoleOwners)
				roleSubGroup.DELETE("", roleController.DeleteRole)

				rolePermissions := roleSubGroup.Group("/permissions")
				{
					rolePermissions.POST("", roleController.UpdateRole)
					rolePermissions.DELETE("/:permissionID", roleController.DeleteRolePermission)
				}
			}
		}

		users.GET("/permissions", roleController.GetPermissionsList)

		userRoles := users.Group("/users/roles")
		{
			userRoles.POST("", roleController.UpdateUserRoles)
			userRoles.DELETE("/:roleID", roleController.DeleteUserRole)
		}
	}

	news := routerGroup.Group("/news")
	news.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ManageNews}))
	{
		news.POST("", newsController.CreateNews)
		newsSubGroup := news.Group("/:newsID")
		{
			newsSubGroup.PUT("", newsController.UpdateNews)
			newsSubGroup.DELETE("", newsController.DeleteNews)
		}
	}

	journals := routerGroup.Group("/journal")
	journals.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ManageJournal}))
	{
		journals.POST("", journalController.CreateJournal)
		journalSubGroup := journals.Group("/:journalID")
		{
			journalSubGroup.PUT("", journalController.UpdateJournal)
			journalSubGroup.DELETE("", journalController.DeleteJournal)
		}
	}
}
