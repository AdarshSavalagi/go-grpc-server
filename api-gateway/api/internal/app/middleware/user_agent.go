// Package middleware provides various middleware functions to enhance
// security, logging, and request handling in the application.
package middleware

import (
	"api/internal/app/constants"
	"api/internal/utils/response"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// UserAgentMiddleware extracts and logs the User-Agent from incoming requests.
//
// It checks for specific User-Agents and takes appropriate action, such as blocking
// outdated browsers. The extracted User-Agent is stored in the request context.
//
// Returns:
// - gin.HandlerFunc: A Gin middleware function.
func UserAgentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the User-Agent from request headers.
		userAgent := c.GetHeader(constants.Const.Headers.UserAgent)

		// Log the User-Agent for analytics or debugging purposes.
		log.Printf("User-Agent: %s", userAgent)

		// Example: Reject requests from an outdated browser.
		if strings.Contains(userAgent, "SomeBrowser/1.0") {
			response.RespondError(c, http.StatusBadRequest, "Unsupported browser. Please update your browser.", nil)
			c.Abort()
			return
		}

		// Store User-Agent and device info in context for further processing.
		c.Set(constants.Const.Context.UserAgent, userAgent)
		c.Set(constants.Const.Context.DeviceInfo, fmt.Sprintf("%s:%s", c.ClientIP(), userAgent))

		// Continue processing the request.
		c.Next()
	}
}
