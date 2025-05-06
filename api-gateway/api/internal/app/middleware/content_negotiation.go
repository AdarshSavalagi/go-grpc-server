// Package middleware provides middleware functions for handling content negotiation.
package middleware

import (
	"api/internal/app/constants"
	"api/internal/utils/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ContentNegotiationMiddleware ensures the server properly handles requests based on the Accept and Content-Type headers.
//
// It verifies that the client accepts a supported response format (JSON or XML) and ensures the request body is in an accepted format.
// If the request's Accept header is invalid, it responds with "406 Not Acceptable".
// If the Content-Type header is missing, it defaults to JSON.
//
// Returns:
// - gin.HandlerFunc: A middleware function for handling content negotiation.
func ContentNegotiationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract Accept header from request
		accept := c.GetHeader(constants.Const.Headers.Accept)
		if accept == "" {
			// Default to JSON if no Accept header is provided
			accept = constants.Const.Headers.AcceptJSON
		}

		// Extract Content-Type header from request
		contentType := c.ContentType()
		if contentType == "" {
			// Default to JSON if no Content-Type header is provided
			contentType = constants.Const.Headers.ContentTypeJSON
		}

		// Log headers for debugging
		log.Printf("Accept: %s, Content-Type: %s", accept, contentType)

		// Validate Accept header
		if accept != constants.Const.Headers.AcceptJSON && accept != constants.Const.Headers.AcceptXML {
			response.Respond(c, false, http.StatusNotAcceptable, "Not Acceptable")
			c.Abort()
			return
		}

		// Set the response Content-Type based on the Accept header
		c.Writer.Header().Set(constants.Const.Headers.ContentType, accept)

		// Proceed to the next middleware or handler
		c.Next()
	}
}
