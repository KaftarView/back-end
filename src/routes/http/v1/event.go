package routes_http_v1

import (
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
	jwtService := application_jwt.NewJWTToken()
	authMiddleware := middleware_authentication.NewAuthMiddleware(di.Constants, userRepository, jwtService)
	eventController := controller_v1_event.NewEventController()
	events := routerGroup.Group("/events")
	{
		read := events.Group("")
		read.Use(func(c *gin.Context) {
			authMiddleware.RequirePermission(c, []enums.PermissionType{enums.CreateEvent})
		})
		{
			events.GET("", eventController.ListEvents)
			events.GET("/:id", eventController.GetEvent)
		}

		events.POST("", func(c *gin.Context) {
			authMiddleware.RequirePermission(c, []enums.PermissionType{enums.CreateEvent})
		}, eventController.CreateEvent)

		createOrEdit := events.Group("")
		createOrEdit.Use(func(c *gin.Context) {
			authMiddleware.RequirePermission(c, []enums.PermissionType{enums.CreateEvent, enums.EditEvent})
		})
		{
			createOrEdit.PUT("/:id", eventController.UpdateEvent)
			createOrEdit.DELETE("/:id", eventController.DeleteEvent)
			createOrEdit.POST("/:id/media", eventController.UploadEventMedia)
			createOrEdit.DELETE("/:id/media/:mediaId", eventController.DeleteEventMedia)
		}

		publish := events.Group("")
		publish.Use(func(c *gin.Context) {
			authMiddleware.RequirePermission(c, []enums.PermissionType{enums.PublishEvent})
		})
		{
			publish.POST("/:id/publish", eventController.PublishEvent)
			publish.POST("/:id/unpublish", eventController.UnpublishEvent)
		}
	}
}
