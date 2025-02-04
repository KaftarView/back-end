package routes_http_v1

import (
	"first-project/src/application"
	application_communication "first-project/src/application/communication/emailService"
	"first-project/src/bootstrap"
	controller_v1_chat "first-project/src/controller/v1/chat"
	controller_v1_comment "first-project/src/controller/v1/comment"
	controller_v1_event "first-project/src/controller/v1/event"
	controller_v1_podcast "first-project/src/controller/v1/podcast"
	controller_v1_user "first-project/src/controller/v1/user"
	"first-project/src/websocket"

	repository_database "first-project/src/repository/database"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupCustomerRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client, hub *websocket.Hub) {
	userRepository := repository_database.NewUserRepository()
	categoryRepository := repository_database.NewCategoryRepository()
	eventRepository := repository_database.NewEventRepository()
	commentRepository := repository_database.NewCommentRepository()
	podcastRepository := repository_database.NewPodcastRepository()
	purchaseRepository := repository_database.NewPurchaseRepository()
	chatRepository := repository_database.NewChatRepository()

	emailService := application_communication.NewEmailService(&di.Env.Email)
	awsService := application.NewS3Service(di.Constants, &di.Env.Storage)
	categoryService := application.NewCategoryService(di.Constants, categoryRepository, db)
	otpService := application.NewOTPService()
	jwtService := application.NewJWTToken()
	userService := application.NewUserService(di.Constants, userRepository, otpService, awsService, db)
	eventService := application.NewEventService(di.Constants, awsService, categoryService, eventRepository, commentRepository, purchaseRepository, db)
	commentService := application.NewCommentService(di.Constants, commentRepository, userService, db)
	podcastService := application.NewPodcastService(di.Constants, awsService, categoryService, podcastRepository, commentRepository, userService, db)
	chatService := application.NewChatService(di.Constants, userService, chatRepository, db)

	customerEventController := controller_v1_event.NewCustomerEventController(di.Constants, eventService, emailService)
	customerCommentController := controller_v1_comment.NewCustomerCommentController(di.Constants, commentService)
	customerPodcastController := controller_v1_podcast.NewCustomerPodcastController(di.Constants, podcastService)
	customerUserController := controller_v1_user.NewCustomerUserController(di.Constants, userService)
	customerChatController := controller_v1_chat.NewCustomerChatController(di.Constants, chatService, jwtService, hub)

	event := routerGroup.Group("/events/:eventID")
	{
		event.GET("/attendance", customerEventController.IsUserAttended)
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

	chat := routerGroup.Group("/chat")
	{
		chat.POST("/room", customerChatController.CreateOrGetRoom)
		chat.GET("/room/:roomID/messages", customerChatController.GetMessages)
	}
}
