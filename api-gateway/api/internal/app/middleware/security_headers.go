// Package middleware provides security enhancements for HTTP responses.
package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware applies security-related HTTP headers to all responses.
//
// It sets headers to prevent content sniffing, clickjacking, and cross-site scripting (XSS).
// In production mode, it enforces HTTPS and applies strict transport security.
//
// Returns:
// - gin.HandlerFunc: A Gin middleware function.
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set security headers
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

		// Apply additional security policies in production mode
		if os.Getenv("APP_ENV") == "production" {
			c.Writer.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

			// Define Content Security Policy (CSP) to restrict sources
			c.Writer.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; form-action 'self'; base-uri 'self';")
		}

		// Enforce HTTPS in production mode
		if c.Request.TLS == nil && os.Getenv("APP_ENV") == "production" {
			c.JSON(http.StatusUpgradeRequired, gin.H{"error": "HTTPS required"})
			c.Abort()
			return
		}

		// Continue request processing
		c.Next()
	}
}
