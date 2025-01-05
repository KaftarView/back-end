package routes_http_v1

import (
	"first-project/src/application"
	application_aws "first-project/src/application/aws"
	application_communication "first-project/src/application/communication/emailService"
	"first-project/src/bootstrap"
	controller_v1_comment "first-project/src/controller/v1/comment"
	controller_v1_event "first-project/src/controller/v1/event"
	controller_v1_podcast "first-project/src/controller/v1/podcast"
	controller_v1_user "first-project/src/controller/v1/user"

	repository_database "first-project/src/repository/database"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupCustomerRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) {
	userRepository := repository_database.NewUserRepository(db)
	categoryRepository := repository_database.NewCategoryRepository(db)
	eventRepository := repository_database.NewEventRepository(db)
	commentRepository := repository_database.NewCommentRepository(db)
	podcastRepository := repository_database.NewPodcastRepository(db)
	purchaseRepository := repository_database.NewPurchaseRepository(db)

	emailService := application_communication.NewEmailService(&di.Env.Email)
	awsService := application_aws.NewS3Service(
		di.Constants, &di.Env.EventsBucket,
		&di.Env.PodcastsBucket, &di.Env.NewsBucket,
		&di.Env.JournalsBucket, &di.Env.ProfilesBucket,
	)
	categoryService := application.NewCategoryService(di.Constants, categoryRepository, db)
	eventService := application.NewEventService(di.Constants, awsService, categoryService, eventRepository, commentRepository, purchaseRepository, db)
	commentService := application.NewCommentService(di.Constants, commentRepository, userRepository, db)
	podcastService := application.NewPodcastService(di.Constants, awsService, categoryService, podcastRepository, commentRepository, userRepository, db)
	otpService := application.NewOTPService()
	userService := application.NewUserService(di.Constants, userRepository, otpService, awsService, db)

	customerEventController := controller_v1_event.NewCustomerEventController(di.Constants, eventService, emailService)
	customerCommentController := controller_v1_comment.NewCustomerCommentController(di.Constants, commentService)
	customerPodcastController := controller_v1_podcast.NewCustomerPodcastController(di.Constants, podcastService)
	customerUserController := controller_v1_user.NewCustomerUserController(di.Constants, userService)

	event := routerGroup.Group("/events/:eventID")
	{
		event.GET("/tickets", customerEventController.GetAvailableEventTicketsList)
		event.GET("/media", customerEventController.GetEventMedia)
		event.POST("/reserve", customerEventController.ReserveTickets)
		event.POST("/purchase/:reservationID", customerEventController.PurchaseTickets)
	}

	comments := routerGroup.Group("/comments")
	{
		comments.POST("/post/:postID", customerCommentController.CreateComment)

		commentSubGroup := comments.Group("/:commentID")
		{
			commentSubGroup.PUT("", customerCommentController.EditComment)
			commentSubGroup.DELETE("", customerCommentController.DeleteComment)
		}
	}

	podcast := routerGroup.Group("/podcasts/:podcastID/subscribe")
	{
		podcast.POST("", customerPodcastController.SubscribePodcast)
		podcast.DELETE("", customerPodcastController.UnSubscribePodcast)
		podcast.GET("/status", customerPodcastController.SubscribeStatus)
	}

	profile := routerGroup.Group("/profile")
	{
		profile.PUT("/username", customerUserController.ChangeUsername)
		profile.PUT("/reset-password", customerUserController.ResetPassword)
		profile.GET("/events", customerEventController.GetAllUserJoinedEvents)
	}
}
