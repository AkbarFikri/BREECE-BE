package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

)

func RateLimiter() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		limiter := rate.NewLimiter(100, 10)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   true,
				"message": "Too many request",
				"data":    nil,
			})
			return
		}

		c.Next()
	})
}
