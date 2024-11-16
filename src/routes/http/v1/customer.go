package routes_http_v1

import (
	"first-project/src/bootstrap"
	controller_v1_general "first-project/src/controller/v1/general"

	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(routerGroup *gin.RouterGroup, di *bootstrap.Di) *gin.RouterGroup {
	routerGroup.GET("/ping1", controller_v1_general.Pong)
	routerGroup.GET("/ping2", controller_v1_general.Pong)
	routerGroup.GET("/ping3", controller_v1_general.Pong)

	return routerGroup
}
