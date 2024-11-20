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
	awsService := application_aws.NewAWSS3(di.Constants, &di.Env.PrimaryBucket)
	eventController := controller_v1_event.NewEventController(di.Constants, eventService, awsService)

	events := routerGroup.Group("/events")
	{
		read := events.Group("")
		read.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ViewReports}))
		{
			read.GET("", eventController.ListEvents)
			read.GET("/:id", eventController.GetEvent)
		}

		create := events.Group("")
		create.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.CreateEvent}))
		{
			create.POST("/create", eventController.CreateEvent)
			create.POST("/add-ticket/:id", eventController.AddEventTicket)
		}

		updateOrDelete := events.Group("")
		updateOrDelete.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.CreateEvent, enums.EditEvent}))
		{
			updateOrDelete.PUT("/:id", eventController.UpdateEvent)
			updateOrDelete.DELETE("/:id", eventController.DeleteEvent)
			updateOrDelete.POST("/:id/media", eventController.UploadEventMedia)
			updateOrDelete.DELETE("/:id/media/:mediaId", eventController.DeleteEventMedia)
		}

		publish := events.Group("")
		publish.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.PublishEvent}))
		{
			publish.POST("/:id/publish", eventController.PublishEvent)
			publish.POST("/:id/unpublish", eventController.UnpublishEvent)
		}
	}
}
