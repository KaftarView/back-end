package routes_http_v1

import (
	"gingool/src/bootstrap"
	controller_v1_general "gingool/src/controller/v1/general"

	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di) *gin.RouterGroup {
	routerGroup.GET("/ping1", controller_v1_general.Pong)
	routerGroup.GET("/ping2", controller_v1_general.Pong)
	routerGroup.GET("/ping3", controller_v1_general.Pong)

	return routerGroup
}
