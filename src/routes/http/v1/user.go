package routes_http_v1

import (
	"first-project/src/application"
	application_communication "first-project/src/application/communication/emailService"
	application_jwt "first-project/src/application/jwt"
	"first-project/src/bootstrap"
	controller_v1_general "first-project/src/controller/v1/general"
	"first-project/src/enums"
	middleware_authentication "first-project/src/middleware/Authentication"
	repository_database "first-project/src/repository/database"
	repository_cache "first-project/src/repository/redis"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupUserRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) {
	userRepository := repository_database.NewUserRepository(db)
	otpService := application.NewOTPService()
	userService := application.NewUserService(di.Constants, userRepository, otpService)
	emailService := application_communication.NewEmailService(&di.Env.Email)
	userCache := repository_cache.NewUserCache(di.Constants, rdb, userRepository)
	jwtService := application_jwt.NewJWTToken()
	userController := controller_v1_general.NewUserController(
		di.Constants, userService, emailService, userCache, otpService, jwtService)

	authMiddleware := middleware_authentication.NewAuthMiddleware(di.Constants, userRepository, jwtService)

	profile := routerGroup.Group("/profile")
	{
		profile.GET("") // some sample
	}

	users := routerGroup.Group("/users")
	users.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ManageUsers}))
	{
		users.PUT("/add-roles", userController.UpdateUserRoles)
	}

	routerGroup.GET(
		"/admin/hello",
		authMiddleware.RequirePermission([]enums.PermissionType{enums.ManageUsers}),
		userController.AdminSayHello,
	)
}
