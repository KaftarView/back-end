package routes_http_v1

import (
	"first-project/src/application"
	application_jwt "first-project/src/application/jwt"
	"first-project/src/bootstrap"
	controller_v1_private "first-project/src/controller/v1/private"
	"first-project/src/enums"
	middleware_authentication "first-project/src/middleware/Authentication"
	repository_database "first-project/src/repository/database"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupUserRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) {
	userRepository := repository_database.NewUserRepository(db)
	otpService := application.NewOTPService()
	userService := application.NewUserService(di.Constants, userRepository, otpService)
	jwtService := application_jwt.NewJWTToken()
	roleController := controller_v1_private.NewRoleController(di.Constants, userService)

	authMiddleware := middleware_authentication.NewAuthMiddleware(di.Constants, userRepository, jwtService)

	profile := routerGroup.Group("/profile")
	{
		profile.GET("") // some sample
	}

	users := routerGroup.Group("/users")
	users.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ManageUsers, enums.ManageRoles}))
	{
		users.POST("/add-role", roleController.CreateRole)
		users.POST("/roles/:roleID/permissions", roleController.UpdateRole)
		users.POST("/:userID/roles", roleController.UpdateUserRoles)
	}
}
