package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jatinnsharma/api/handlers"
	"github.com/jatinnsharma/internal/auth"
	"github.com/jatinnsharma/internal/config"
	"github.com/jatinnsharma/internal/middleware"
	// "github.com/jatinnsharma/internal/otp"
	"github.com/jatinnsharma/internal/token"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, redisClient *redis.Client, cfg *config.Config) {
	// Initialize services
	tokenService := token.NewService(cfg)
	// otpService := otp.NewService(db, cfg)
	authService := auth.NewService(db, tokenService, cfg)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Apply global middleware
	router.Use(middleware.RateLimitMiddleware(redisClient, cfg))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/signup", authHandler.Signup)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/refresh", authHandler.RefreshToken)
			authRoutes.POST("/logout", authHandler.Logout)
		}

		// Protected routes
		protected := v1.Group("/user")
		protected.Use(middleware.AuthMiddleware(tokenService))
		{
			protected.GET("/me", func(c *gin.Context) {
				userID := c.GetString("user_id")
				c.JSON(200, gin.H{"user_id": userID})
			})
		}
	}
}