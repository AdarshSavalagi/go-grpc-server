// Package middleware provides utility functions to enhance request processing.
package middleware

import (
	"api/internal/app/constants"
	"api/internal/utils"
	"log"

	"github.com/gin-gonic/gin"
)

// RequestIDMiddleware generates a unique request ID for each incoming request.
//
// If a request ID is already provided in the headers, it is used; otherwise, a new UUID is generated.
// The request ID is stored in the request context and added to the response headers.
//
// Returns:
// - gin.HandlerFunc: A Gin middleware function.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the request ID header key
		RequestIDKey := constants.Const.Headers.RequestID

		// Check if a request ID is provided by the client
		requestID := c.GetHeader(RequestIDKey)
		if requestID == "" {
			// Generate a new request ID if not provided
			requestID = utils.GenerateUUID()
		}

		// Store the request ID in the context for further use
		c.Set(constants.Const.Context.RequestID, requestID)

		// Add the request ID to the response headers for tracking
		log.Printf("Storing request ID in header %s: %s", RequestIDKey, requestID)
		c.Writer.Header().Set(RequestIDKey, requestID)

		// Continue processing the request
		c.Next()
	}
}
