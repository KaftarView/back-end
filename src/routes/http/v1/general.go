package routes_http_v1

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"first-project/src/application"
	application_communication "first-project/src/application/communication/emailService"
	application_jwt "first-project/src/application/jwt"
	"first-project/src/bootstrap"
	controller_v1_general "first-project/src/controller/v1/general"
	repository_database "first-project/src/repository/database"
	repository_cache "first-project/src/repository/redis"
)

func SetupGeneralRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) *gin.RouterGroup {
	userRepository := repository_database.NewUserRepository(db)

	otpService := application.NewOTPService()
	userService := application.NewUserService(di.Constants, userRepository, otpService)
	emailService := application_communication.NewEmailService(&di.Env.Email)
	userCache := repository_cache.NewUserCache(di.Constants, rdb, userRepository)
	jwtService := application_jwt.NewJWTToken()
	userController := controller_v1_general.NewUserController(
		di.Constants, userService, emailService, userCache, otpService, jwtService)

	authController := controller_v1_general.NewAuthController(di.Constants, jwtService)

	routerGroup.POST("/register", userController.Register)
	routerGroup.POST("/register/verify", userController.VerifyEmail)
	routerGroup.POST("/login", userController.Login)
	routerGroup.POST("/forgot-password", userController.ForgotPassword)
	routerGroup.POST("/confirm-otp", userController.ConfirmOTP)
	routerGroup.PUT("/reset-password", userController.ResetPassword)
	routerGroup.POST("/refresh-token", authController.RefreshToken)

	return routerGroup
}
