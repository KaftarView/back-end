package routes_http_v1

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"first-project/src/application"
	application_communication "first-project/src/application/communication/emailService"
	application_math "first-project/src/application/math"
	"first-project/src/bootstrap"
	controller_v1_general "first-project/src/controller/v1/general"
	"first-project/src/enums"
	middleware_authentication "first-project/src/middleware/Authentication"
	"first-project/src/repository"
)

func SetupGeneralRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB) *gin.RouterGroup {
	userRepository := repository.NewUserRepository(db)
	addService := application_math.NewAddService(userRepository)
	sampleController := controller_v1_general.NewSampleController(di.Constants, addService)

	userService := application.NewUserService(di.Constants, userRepository)
	emailService := application_communication.NewEmailService(&di.Env.Email)
	userController := controller_v1_general.NewUserController(
		di.Constants, userService, emailService)

	authMiddleware := middleware_authentication.NewAuthMiddleware(di.Constants, userRepository)
	authController := controller_v1_general.NewAuthController(di.Constants)

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

	return routerGroup
}
