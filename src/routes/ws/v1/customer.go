package routes_websocket_v1

import (
	"first-project/src/wire"

	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(ws *gin.RouterGroup, app *wire.Application) {
	chat := ws.Group("/chat")
	{
		chat.GET("/room/:roomID/token/:token", app.CustomerControllers.ChatController.HandleWebsocket)
	}
}
