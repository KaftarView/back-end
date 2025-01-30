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

func NewRateLimit() *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limit: 5,
		burst: 10,
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
