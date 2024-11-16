package controller_v1_general

import (
	"github.com/gin-gonic/gin"
)

// Response sends a standard JSON response with status and message
func Response(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, SwaggerResponseMessage{
		Status:  statusCode,
		Message: message,
	})
}

// ResponseMessage represents a standard structure for API responses
type SwaggerResponseMessage struct {
	Status  int    `json:"status"`  // Status code (e.g., 200, 400, etc.)
	Message string `json:"message"` // Message describing the response
}
