package routes

import (
	routes_http_v1 "first-project/src/routes/http/v1"
	routes_websocket_v1 "first-project/src/routes/ws/v1"
	"first-project/src/wire"

	"github.com/gin-gonic/gin"
)

func Run(ginEngine *gin.Engine, app *wire.Application) {
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
