package routes

import (
	"first-project/src/bootstrap"
	routes_http_v1 "first-project/src/routes/http/v1"
	routes_websocket_v1 "first-project/src/routes/ws/v1"
	"first-project/src/websocket"
	"first-project/src/wire"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Run(ginEngine *gin.Engine, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client, hub *websocket.Hub) {
	app, err := wire.InitializeApplication(di, db, rdb, hub)
	if err != nil {
		panic(err)
	}

	ginEngine.Use(app.Middlewares.Localization.Localization)
	ginEngine.Use(app.Middlewares.Recovery.Recovery)
	ginEngine.Use(app.Middlewares.RateLimit.RateLimit)

	v1 := ginEngine.Group("/v1")

	registerGeneralRoutes(v1, app)
	registerCustomerRoutes(v1, app)
	registerAdminRoutes(v1, app)
}

func registerGeneralRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	routes_http_v1.SetupGeneralRoutes(v1, app)
}

func registerCustomerRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	customerGroup := v1.Group("")
	customerGroup.Use(app.Middlewares.Auth.AuthRequired)
	routes_http_v1.SetupCustomerRoutes(customerGroup, app)

	ws := v1.Group("/ws")
	ws.Use(app.Middlewares.Websocket.UpgradeToWebSocket)
	routes_websocket_v1.SetupCustomerRoutes(ws, app)
}

func registerAdminRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	adminGroup := v1.Group("/admin")
	adminGroup.Use(app.Middlewares.Auth.AuthRequired)
	routes_http_v1.SetupAdminRoutes(adminGroup, app)
}
