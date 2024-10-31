package middleware_rate_limit

import (
	"first-project/src/exceptions"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimitMiddleware struct {
	limit rate.Limit
	burst int
}

func NewRateLimit(limit rate.Limit, burst int) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limit: limit,
		burst: burst,
	}
}

func (rl *RateLimitMiddleware) RateLimit(c *gin.Context) {
	limiter := rate.NewLimiter(rl.limit, rl.burst)
	if !limiter.Allow() {
		rateLimitError := exceptions.NewRateLimitError()
		panic(rateLimitError)
	}
	c.Next()
}
