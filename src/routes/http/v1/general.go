package routes_http_v1

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"first-project/src/application"
	application_aws "first-project/src/application/aws"
	application_communication "first-project/src/application/communication/emailService"
	application_math "first-project/src/application/math"
	"first-project/src/bootstrap"
	controller_v1_general "first-project/src/controller/v1/general"
	"first-project/src/enums"
	middleware_authentication "first-project/src/middleware/Authentication"
	cache "first-project/src/redis"
	"first-project/src/repository"
)

func SetupGeneralRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) *gin.RouterGroup {
	userRepository := repository.NewUserRepository(db)
	addService := application_math.NewAddService(userRepository)
	sampleController := controller_v1_general.NewSampleController(di.Constants, addService)

	otpService := application.NewOTPService()
	userService := application.NewUserService(di.Constants, userRepository, otpService)
	emailService := application_communication.NewEmailService(&di.Env.Email)
	userCache := cache.NewUserCache(rdb, userRepository)
	userController := controller_v1_general.NewUserController(
		di.Constants, userService, emailService, userCache, otpService)

	authMiddleware := middleware_authentication.NewAuthMiddleware(di.Constants, userRepository)
	authController := controller_v1_general.NewAuthController(di.Constants)

	awsService := application_aws.NewAWSS3(&di.Env.PrimaryBucket)
	awsController := controller_v1_general.NewAWSController(di.Constants, awsService)

	routerGroup.GET("/ping", controller_v1_general.Pong)
	routerGroup.GET("/add/:num1/:num2", sampleController.Add)
	routerGroup.POST("/register", userController.Register)
	routerGroup.POST("/register/activate", userController.VerifyEmail)
	routerGroup.POST("/login", userController.Login)
	routerGroup.POST("/forgotPassword", userController.ForgotPassword)
	routerGroup.PUT("/resetPassword", userController.ResetPassword)
	routerGroup.GET("/admin", func(c *gin.Context) {
		authMiddleware.AuthenticateMiddleware(c, []enums.RoleType{enums.Admin})
	}, userController.AdminSayHello)
	routerGroup.POST("/refreshToken", authController.RefreshToken)
	routerGroup.POST("/bucket/create", awsController.CreateBucketController)
	routerGroup.POST("/bucket/upload", awsController.UploadObjectController)
	routerGroup.POST("/bucket/delete", awsController.DeleteObjectController)

	return routerGroup
}
