package controller_v1_general

import (
	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, SwaggerResponseMessage{
		Status:  statusCode,
		Message: message,
	})
}

type SwaggerResponseMessage struct {
	Status  int    `json:"status"`  // Status code (e.g., 200, 400, etc.)
	Message string `json:"message"` // Message describing the response
}
