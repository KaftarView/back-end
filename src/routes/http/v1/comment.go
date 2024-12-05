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

func SetupCommentRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) {
	userRepository := repository_database.NewUserRepository(db)
	commentRepository := repository_database.NewCommentRepository(db)
	jwtService := application_jwt.NewJWTToken()
	authMiddleware := middleware_authentication.NewAuthMiddleware(di.Constants, userRepository, jwtService)
	commentService := application.NewCommentService(di.Constants, commentRepository, userRepository)
	commentController := controller_v1_private.NewCommentController(di.Constants, commentService)
	comments := routerGroup.Group("/comments")
	{
		crudUser := comments.Group("")
		{
			crudUser.POST("/post/:postID", commentController.CreateComment)
			crudUser.PUT("/:commentID", commentController.EditComment)
			crudUser.DELETE("/:commentID", commentController.DeleteComment)
		}

		crudAdmin := comments.Group("")
		crudAdmin.Use(authMiddleware.RequirePermission([]enums.PermissionType{enums.ModerateComments}))
		{
			crudUser.DELETE("/admin/:id")
		}
	}
}
