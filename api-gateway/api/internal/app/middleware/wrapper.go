// Package middleware provides various middleware functions to enhance
// security, logging, performance, and request handling in the application.
package middleware

import (
	"api/internal/app/config"
	"github.com/gin-gonic/gin"
)

// MiddlewareWrapper applies middleware configurations to the Gin router.
//
// It dynamically enables or disables middlewares based on the application's
// configuration settings.
//
// Parameters:
// - r: A pointer to the Gin engine (router) to which middleware will be applied.
func MiddlewareWrapper(r *gin.Engine) {
	// Apply panic recovery middleware.
	r.Use(RecoveryMiddleware())

	// Apply User-Agent middleware if enabled.
	if config.MiddlewareConfig.UserAgent {
		r.Use(UserAgentMiddleware())
	}

	// Apply ClientID middleware if enabled.
	if config.MiddlewareConfig.ClientID {
		r.Use(ClientIDMiddleware())
	}

	// Add RequestID middleware if enabled.
	if config.MiddlewareConfig.RequestID {
		r.Use(RequestIDMiddleware())
	}

	// Apply Logger Middleware if enabled.
	if config.MiddlewareConfig.Logger {
		r.Use(LoggingMiddleware())
	}

	// Apply Gzip compression middleware if enabled.
	if config.MiddlewareConfig.GzipCompression {
		r.Use(Compression())
	}

	// Add security headers middleware if enabled.
	if config.MiddlewareConfig.SecurityHeaders {
		r.Use(SecurityHeadersMiddleware())
	}

	// Add CORS middleware if enabled.
	if config.MiddlewareConfig.CORS {
		r.Use(CORSMiddleware())
	}

	// Apply default no-cache middleware if enabled.
	if config.MiddlewareConfig.Cache {
		r.Use(CacheControlMiddleware(CacheConfig{}))
	}

	// Apply content negotiation middleware if enabled.
	if config.MiddlewareConfig.ContentNegotiation {
		r.Use(ContentNegotiationMiddleware())
	}

	// Apply Referer middleware if enabled.
	if config.MiddlewareConfig.Referer {
		r.Use(RefererMiddleware())
	}

	// Apply cookies middleware if enabled.
	if config.MiddlewareConfig.Cookies {
		r.Use(CookieMiddleware())
	}

	// Ensure proper exit control if any middleware issues an abort.
	r.Use(ExitControlFlow())
}
