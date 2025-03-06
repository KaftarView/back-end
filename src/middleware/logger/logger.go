package middleware_logger

import (
	"first-project/src/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoggerMiddleware struct{}

func NewLoggerMiddleware() *LoggerMiddleware {
	return &LoggerMiddleware{}
}

func (log *LoggerMiddleware) GinLoggerMiddleware(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery

	c.Next()

	logger := logger.GetLogger()
	end := time.Now()
	latency := end.Sub(start)

	if len(c.Errors) > 0 {
		for _, e := range c.Errors.Errors() {
			logger.Error(e)
		}
	} else {
		logger.Info(
			"Request",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.String("user-agent", c.Request.UserAgent()),
		)
	}
}
