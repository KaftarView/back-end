package routes_websocket_v1

import (
	"first-project/src/application"
	application_aws "first-project/src/application/aws"
	application_jwt "first-project/src/application/jwt"
	"first-project/src/bootstrap"
	controller_v1_chat "first-project/src/controller/v1/chat"
	repository_database "first-project/src/repository/database"
	"first-project/src/websocket"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupCustomerRoutes(ws *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client, hub *websocket.Hub) {
	chatRepository := repository_database.NewChatRepository()
	userRepository := repository_database.NewUserRepository()

	awsService := application_aws.NewS3Service(
		di.Constants, &di.Env.EventsBucket,
		&di.Env.PodcastsBucket, &di.Env.NewsBucket,
		&di.Env.JournalsBucket, &di.Env.ProfilesBucket,
	)
	otpService := application.NewOTPService()
	jwtService := application_jwt.NewJWTToken()
	userService := application.NewUserService(di.Constants, userRepository, otpService, awsService, db)
	chatService := application.NewChatService(di.Constants, userService, chatRepository, db)

	customerChatController := controller_v1_chat.NewCustomerChatController(di.Constants, chatService, jwtService, hub)

	chat := ws.Group("/chat")
	{
		chat.GET("/room/:roomID/token/:token", customerChatController.HandleWebsocket)
	}
}
