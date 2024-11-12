package routes_http_v1

import (
	"first-project/src/bootstrap"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupViewerRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di, db *gorm.DB, rdb *redis.Client) *gin.RouterGroup {
	return routerGroup
}
