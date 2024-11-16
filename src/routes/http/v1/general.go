package routes_http_v1

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"first-project/src/application"
	application_math "first-project/src/application/math"
	"first-project/src/bootstrap"
	controller_v1_general "first-project/src/controller/v1/general"
	"first-project/src/repository"
)

func SetupGeneralRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB) *gin.RouterGroup {
	userRepository := repository.NewUserRepository(db)
	addService := application_math.NewAddService(userRepository)
	sampleController := controller_v1_general.NewSampleController(di.Constants, addService)

	userService := application.NewUserService(di.Constants, userRepository)
	otpService := application.NewOTPService(di.Constants, userRepository)
	emailService := application.NewEmailService(&di.Env.Email)
	userController := controller_v1_general.NewUserController(
		di.Constants, userService, otpService, emailService,
	)

	routerGroup.GET("/ping", controller_v1_general.Pong)
	routerGroup.GET("/add/:num1/:num2", sampleController.Add)
	routerGroup.POST("/register", userController.Register)
	routerGroup.POST("/register/activate", userController.VerifyEmail)
	routerGroup.POST("/login", userController.Login)

	return routerGroup
}
