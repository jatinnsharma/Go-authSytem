package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jatinnsharma/internal/config"
	"github.com/jatinnsharma/internal/utils"
)

func RateLimitMiddleware(redisClient *redis.Client, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", ip)
		
		ctx := context.Background()
		
		// Get current count
		val, err := redisClient.Get(ctx, key).Result()
		if err != nil && err != redis.Nil {
			// If Redis is down, allow the request
			c.Next()
			return
		}

		var count int
		if val != "" {
			count, _ = strconv.Atoi(val)
		}

		if count >= cfg.RateLimitPerMin {
			utils.ErrorResponse(c, http.StatusTooManyRequests, "Rate limit exceeded", nil)
			c.Abort()
			return
		}

		// Increment counter
		pipe := redisClient.Pipeline()
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, time.Minute)
		pipe.Exec(ctx)

		c.Next()
	}
}