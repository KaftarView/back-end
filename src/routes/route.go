package routes

import (
	"first-project/src/bootstrap"
	middleware_exceptions "first-project/src/middleware/exceptions"
	middleware_i18n "first-project/src/middleware/i18n"
	middleware_rate_limit "first-project/src/middleware/rateLimit"
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
	registerSuperAdminRoutes(v1, di, db, rdb)
	registerEventManagerRoutes(v1, di, db, rdb)
	registerContentManagerRoutes(v1, di, db, rdb)
	registerEditorRoutes(v1, di, db, rdb)
	registerModeratorRoutes(v1, di, db, rdb)
	registerViewerRoutes(v1, di, db, rdb)
}

func registerGeneralRoutes(v1 *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) *gin.RouterGroup {
	return routes_http_v1.SetupGeneralRoutes(v1, di, db, rdb)
}

func registerSuperAdminRoutes(v1 *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) *gin.RouterGroup {
	return routes_http_v1.SetupSuperAdminRoutes(v1, di, db, rdb)
}

func registerEventManagerRoutes(v1 *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) *gin.RouterGroup {
	return routes_http_v1.SetupEventManagerRoutes(v1, di, db, rdb)
}

func registerContentManagerRoutes(v1 *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) *gin.RouterGroup {
	return routes_http_v1.SetupContentManagerRoutes(v1, di, db, rdb)
}

func registerEditorRoutes(v1 *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) *gin.RouterGroup {
	return routes_http_v1.SetupEditorRoutes(v1, di, db, rdb)
}

func registerModeratorRoutes(v1 *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) *gin.RouterGroup {
	return routes_http_v1.SetupModeratorRoutes(v1, di, db, rdb)
}

func registerViewerRoutes(v1 *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) *gin.RouterGroup {
	return routes_http_v1.SetupViewerRoutes(v1, di, db, rdb)
}
