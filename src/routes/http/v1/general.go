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
	controller_v1_event "first-project/src/controller/v1/event"
	controller_v1_general "first-project/src/controller/v1/general"
	repository_database "first-project/src/repository/database"
	repository_cache "first-project/src/repository/redis"
)

func SetupGeneralRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) {
	userRepository := repository_database.NewUserRepository(db)

	otpService := application.NewOTPService()
	userService := application.NewUserService(di.Constants, userRepository, otpService)
	emailService := application_communication.NewEmailService(&di.Env.Email)
	userCache := repository_cache.NewUserCache(di.Constants, rdb, userRepository)
	jwtService := application_jwt.NewJWTToken()
	userController := controller_v1_general.NewUserController(
		di.Constants, userService, emailService, userCache, otpService, jwtService)
	authController := controller_v1_general.NewAuthController(di.Constants, jwtService)

	eventRepository := repository_database.NewEventRepository(db)
	eventService := application.NewEventService(di.Constants, eventRepository)
	awsService := application_aws.NewAWSS3(di.Constants, &di.Env.PrimaryBucket)
	eventController := controller_v1_event.NewEventController(di.Constants, eventService, awsService)

	public := routerGroup.Group("/public")
	{
		public.GET("/categories", eventController.ListCategories)

		events := public.Group("/events")
		{
			events.GET("/published", eventController.ListPublicEvents)
			events.GET("/:eventID", eventController.GetPublicEvent)
			events.GET("/search", eventController.SearchPublicEvents)
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
