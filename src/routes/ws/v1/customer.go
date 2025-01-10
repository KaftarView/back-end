package routes_websocket_v1

import (
	"first-project/src/application"
	"first-project/src/bootstrap"
	controller_v1_chat "first-project/src/controller/v1/chat"
	repository_database "first-project/src/repository/database"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupCustomerRoutes(ws *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) {
	chatRepository := repository_database.NewChatRepository()
	chatService := application.NewChatService(di.Constants, chatRepository, db)
	customerChatController := controller_v1_chat.NewCustomerCommentController(di.Constants, chatService)

	chat := ws.Group("/chat")
	{
		chat.GET("", customerChatController.HandleWebsocket)
	}
}
