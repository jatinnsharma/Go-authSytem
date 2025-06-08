package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jatinnsharma/internal/token"
	"github.com/jatinnsharma/internal/utils"
)

func AuthMiddleware(tokenService *token.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header required")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.UnauthorizedResponse(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		claims, err := tokenService.ValidateAccessToken(tokenString)
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}