package middleware

import (
	"api/internal/utils/auth"
	"api/internal/utils/response"
	"strings"

	"api/internal/app/constants"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// JWTAuthMiddleware validates JWT tokens and extracts user information from them
func JWTAuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the Authorization header
		authHeader := c.GetHeader(constants.Const.Headers.AuthToken)
		if authHeader == "" {
			// If no token is found, return unauthorized response
			response.RespondUnauthorized(c)
			c.Abort()
			return
		}

		// Remove "Bearer " prefix from the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Initialize claims struct
		claims := &auth.Claims{}

		// Parse and validate JWT token, extracting claims
		claims, err := auth.GetJwtClaims(tokenString)
		if err != nil {
			// If token is invalid, return unauthorized response
			response.RespondUnauthorized(c)
			c.Abort()
			return
		}

		// Store user information in the Gin context
		c.Set(constants.Const.Context.UserID, claims.Sub)
		c.Set(constants.Const.Context.Username, claims.Email)

		// Proceed to the next middleware or handler
		c.Next()
	}
}

// JWTReissueAuthMiddleware handles refresh token validation and token reissuance
func JWTReissueAuthMiddleware(redisClient *redis.Client, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Placeholder for future implementation of refresh token validation

		// Continue processing the request
		c.Next()
	}
}
