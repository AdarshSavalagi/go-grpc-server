// Package middleware provides various middleware functions for handling request headers, cookies, and caching.
package middleware

import (
	"api/internal/app/constants"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RefererMiddleware logs the referer header from incoming requests.
//
// If a referer header is present, it is logged for debugging or tracking purposes.
//
// Returns:
// - gin.HandlerFunc: A Gin middleware function.
func RefererMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the Referer header
		referer := c.GetHeader(constants.Const.Headers.Referer)
		if referer != "" {
			log.Printf("Referer: %s", referer)
		}

		// Continue processing the request
		c.Next()
	}
}

// CookieMiddleware logs all cookies received in the request.
//
// It iterates through all cookies and logs their names and values for debugging purposes.
//
// Returns:
// - gin.HandlerFunc: A Gin middleware function.
func CookieMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve and log cookies
		cookies := c.Request.Cookies()
		for _, cookie := range cookies {
			log.Printf("Cookie: %s = %s", cookie.Name, cookie.Value)
		}

		// Continue processing the request
		c.Next()
	}
}

// ExitControlFlow checks if a cached response exists and serves it if available.
//
// If a cached response is found in the request context, it is returned immediately, preventing further processing.
// Otherwise, the request continues normally.
//
// Returns:
// - gin.HandlerFunc: A Gin middleware function.
func ExitControlFlow() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if cached content exists
		if cachedResponse, exists := c.Get(constants.Const.Context.CachedContent); exists && cachedResponse != nil {
			// Set cache response header
			c.Writer.Header().Set(constants.Const.Headers.CacheContent, "true")

			// Respond with cached content and stop request processing
			c.AbortWithStatusJSON(http.StatusOK, cachedResponse)
			return
		}

		// Continue processing if no cached content is found
		c.Next()
	}
}
