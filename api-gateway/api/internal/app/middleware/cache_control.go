package middleware

import (
	"api/internal/app/constants"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CacheConfig defines caching options for HTTP responses.
type CacheConfig struct {
	Public       bool          // If true, the response can be cached by public caches (e.g., CDNs)
	MaxAge       time.Duration // The duration for which the response should be cached
	ETag         string        // Unique identifier for the response content (used for cache validation)
	LastModified time.Time     // Timestamp of the last modification (used for cache validation)
}

// DefaultCacheConfig defines the default cache settings (no caching by default).
var DefaultCacheConfig = CacheConfig{
	Public: false, // Default to private caching
	MaxAge: 0,     // No caching unless explicitly set
}

// CacheControlMiddleware applies cache control headers based on the provided configuration.
func CacheControlMiddleware(config CacheConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ========== BEFORE REQUEST HANDLING ==========

		// Apply the default config if no specific config is provided
		if config == (CacheConfig{}) {
			config = DefaultCacheConfig
		}

		// Determine the Cache-Control header value
		cacheControl := "no-cache" // Default to no caching
		if config.MaxAge > 0 {     // If a Max-Age is set, allow caching
			cacheType := "private" // Default to private caching (per-user cache)
			if config.Public {     // If public caching is allowed (e.g., for static assets)
				cacheType = "public"
			}
			cacheControl = cacheType + ", max-age=" + strconv.Itoa(int(config.MaxAge.Seconds()))
		}

		// Set the Cache-Control header in the response
		c.Writer.Header().Set(constants.Const.Headers.CacheControl, cacheControl)

		// ========== ETag (Entity Tag) CACHE VALIDATION ==========
		if config.ETag != "" {
			// Set ETag header in the response
			c.Writer.Header().Set(constants.Const.Headers.CacheControlETag, config.ETag)

			// Check if the client already has this ETag (sent via "If-None-Match" header)
			if match := c.Request.Header.Get(constants.Const.Headers.CacheControlNoneMatch); match == config.ETag {
				// If ETag matches, return 304 (Not Modified) to prevent unnecessary data transfer
				c.Writer.WriteHeader(http.StatusNotModified)
				c.Abort()
				return
			}
		}

		// ========== Last-Modified CACHE VALIDATION ==========
		if !config.LastModified.IsZero() {
			// Set the Last-Modified header in the response
			lastModifiedStr := config.LastModified.Format(time.RFC1123)
			c.Writer.Header().Set(constants.Const.Headers.CacheControlLastModified, lastModifiedStr)

			// Check if the client has a cached version using "If-Modified-Since"
			if since := c.Request.Header.Get(constants.Const.Headers.CacheControlModifiedSince); since != "" {
				// Parse the timestamp from the client's request
				if t, err := time.Parse(time.RFC1123, since); err == nil && config.LastModified.Before(t.Add(time.Second)) {
					// If the cached version is still valid, return 304 (Not Modified)
					c.Writer.WriteHeader(http.StatusNotModified)
					c.Abort()
					return
				}
			}
		}

		// Proceed to the next middleware/handler
		c.Next()

		// ========== AFTER REQUEST HANDLING (Post-processing) ==========
	}
}
