package routes_http_v1

import (
	"first-project/src/application"
	application_aws "first-project/src/application/aws"
	application_jwt "first-project/src/application/jwt"
	"first-project/src/bootstrap"
	controller_v1_event "first-project/src/controller/v1/event"
	"first-project/src/enums"
	middleware_authentication "first-project/src/middleware/Authentication"
	repository_database "first-project/src/repository/database"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupEventRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) {
	userRepository := repository_database.NewUserRepository(db)
	eventRepository := repository_database.NewEventRepository(db)
	jwtService := application_jwt.NewJWTToken()
	authMiddleware := middleware_authentication.NewAuthMiddleware(di.Constants, userRepository, jwtService)
	eventService := application.NewEventService(di.Constants, eventRepository)
	awsService := application_aws.NewS3Service(di.Constants, &di.Env.BannersBucket, &di.Env.SessionsBucket, &di.Env.PodcastsBucket)
	eventController := controller_v1_event.NewEventController(di.Constants, eventService, awsService)

	events := routerGroup.Group("/events")
	{
		read := events.Group("")
		read.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ViewReports}))
		{
			read.GET("", eventController.GetEventsListForAdmin)
			read.GET("/event-details/:id", eventController.GetEventDetailsForAdmin)
			read.GET("/ticket-details/:id", eventController.GetTicketDetails)
			read.GET("/discount-details/:id", eventController.GetDiscountDetails)
		}

		create := events.Group("")
		create.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.CreateEvent}))
		{
			create.POST("/create", eventController.CreateEvent)
			create.POST("/add-ticket/:eventID", eventController.AddEventTicket)
			create.POST("/add-discount/:eventID", eventController.AddEventDiscount)
		}

		updateOrDelete := events.Group("")
		updateOrDelete.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.CreateEvent, enums.EditEvent}))
		{
			updateOrDelete.PUT("/:eventID", eventController.UpdateEvent)
			updateOrDelete.DELETE("/:eventID", eventController.DeleteEvent)
			updateOrDelete.DELETE("/:eventID/ticket/:ticketID", eventController.DeleteTicket)
			updateOrDelete.DELETE("/:eventID/discount/:discountID", eventController.DeleteDiscount)
			updateOrDelete.POST("/:eventID/media", eventController.UploadEventMedia)
			updateOrDelete.DELETE("/:eventID/media/:mediaId", eventController.DeleteEventMedia)
		}

		publish := events.Group("")
		publish.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.PublishEvent}))
		{
			publish.POST("/:eventID/publish", eventController.PublishEvent)
			publish.POST("/:eventID/unpublish", eventController.UnpublishEvent)
		}
	}
}
