// Package middleware provides middleware functions for handling request compression.
package middleware

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// Compression applies GZIP compression to responses, except for the "/metrics" endpoint.
//
// This middleware uses the default compression level to optimize response size while maintaining performance.
// The "/metrics" endpoint is excluded to prevent issues with Prometheus scraping.
//
// Returns:
// - gin.HandlerFunc: A middleware function for compressing responses.
func Compression() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Apply GZIP compression to all responses except "/metrics"
		if c.FullPath() != "/metrics" {
			gzip.Gzip(gzip.DefaultCompression)(c)
		}

		// Proceed to the next middleware or handler
		c.Next()
	}
}
