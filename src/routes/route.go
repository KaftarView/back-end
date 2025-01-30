package routes

import (
	application_jwt "first-project/src/application/jwt"
	"first-project/src/bootstrap"
	middleware_authentication "first-project/src/middleware/Authentication"
	middleware_exceptions "first-project/src/middleware/exceptions"
	middleware_i18n "first-project/src/middleware/i18n"
	middleware_rate_limit "first-project/src/middleware/rateLimit"
	middleware_websocket "first-project/src/middleware/websocket"
	repository_database "first-project/src/repository/database"
	routes_http_v1 "first-project/src/routes/http/v1"
	routes_websocket_v1 "first-project/src/routes/ws/v1"
	"first-project/src/websocket"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Run(ginEngine *gin.Engine, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client, hub *websocket.Hub) {
	localizationMiddleware := middleware_i18n.NewLocalization(&di.Constants.Context)
	recoveryMiddleware := middleware_exceptions.NewRecovery(&di.Constants.Context)
	rateLimitMiddleware := middleware_rate_limit.NewRateLimit()

	ginEngine.Use(localizationMiddleware.Localization)
	ginEngine.Use(recoveryMiddleware.Recovery)
	ginEngine.Use(rateLimitMiddleware.RateLimit)

	v1 := ginEngine.Group("/v1")

	registerGeneralRoutes(v1, di, db, rdb)
	registerCustomerRoutes(v1, di, db, rdb, hub)
	registerAdminRoutes(v1, di, db, rdb)
}

func registerGeneralRoutes(v1 *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) {
	routes_http_v1.SetupGeneralRoutes(v1, di, db, rdb)
}

func registerCustomerRoutes(v1 *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client, hub *websocket.Hub) {
	userRepository := repository_database.NewUserRepository()
	jwtService := application_jwt.NewJWTToken()
	authMiddleware := middleware_authentication.NewAuthMiddleware(di.Constants, userRepository, jwtService, db)
	wsMiddleware := middleware_websocket.NewWebsocketMiddleware(di.Constants)

	customerGroup := v1.Group("")
	customerGroup.Use(authMiddleware.AuthRequired)
	routes_http_v1.SetupCustomerRoutes(customerGroup, di, db, rdb, hub)

	ws := v1.Group("/ws")
	ws.Use(wsMiddleware.UpgradeToWebSocket)
	routes_websocket_v1.SetupCustomerRoutes(ws, di, db, rdb, hub)
}

func registerAdminRoutes(v1 *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) {
	userRepository := repository_database.NewUserRepository()
	jwtService := application_jwt.NewJWTToken()
	authMiddleware := middleware_authentication.NewAuthMiddleware(di.Constants, userRepository, jwtService, db)

	adminGroup := v1.Group("/admin")
	adminGroup.Use(authMiddleware.AuthRequired)
	routes_http_v1.SetupAdminRoutes(adminGroup, di, db, rdb)
}
