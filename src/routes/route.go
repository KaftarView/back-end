package routes

import (
	application_jwt "first-project/src/application/jwt"
	"first-project/src/bootstrap"
	middleware_authentication "first-project/src/middleware/Authentication"
	middleware_exceptions "first-project/src/middleware/exceptions"
	middleware_i18n "first-project/src/middleware/i18n"
	middleware_rate_limit "first-project/src/middleware/rateLimit"
	repository_database "first-project/src/repository/database"
	routes_http_v1 "first-project/src/routes/http/v1"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Run(ginEngine *gin.Engine, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) {
	localizationMiddleware := middleware_i18n.NewLocalization(&di.Constants.Context)
	recoveryMiddleware := middleware_exceptions.NewRecovery(&di.Constants.Context)
	rateLimitMiddleware := middleware_rate_limit.NewRateLimit(5, 10)

	ginEngine.Use(localizationMiddleware.Localization)
	ginEngine.Use(recoveryMiddleware.Recovery)
	ginEngine.Use(rateLimitMiddleware.RateLimit)

	v1 := ginEngine.Group("/v1")

	registerGeneralRoutes(v1, di, db, rdb)
	registerCustomerRoutes(v1, di, db, rdb)
	registerAdminRoutes(v1, di, db, rdb)
}

func registerGeneralRoutes(v1 *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) {
	routes_http_v1.SetupGeneralRoutes(v1, di, db, rdb)
}

func registerCustomerRoutes(v1 *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) {
	userRepository := repository_database.NewUserRepository(db)
	jwtService := application_jwt.NewJWTToken()
	authMiddleware := middleware_authentication.NewAuthMiddleware(di.Constants, userRepository, jwtService)

	customerGroup := v1.Group("")
	customerGroup.Use(authMiddleware.AuthRequired)
	routes_http_v1.SetupCustomerRoutes(customerGroup, di, db, rdb)
}

func registerAdminRoutes(v1 *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) {
	userRepository := repository_database.NewUserRepository(db)
	jwtService := application_jwt.NewJWTToken()
	authMiddleware := middleware_authentication.NewAuthMiddleware(di.Constants, userRepository, jwtService)

	adminGroup := v1.Group("/admin")
	adminGroup.Use(authMiddleware.AuthRequired)
	routes_http_v1.SetupAdminRoutes(adminGroup, di, db, rdb)
}
