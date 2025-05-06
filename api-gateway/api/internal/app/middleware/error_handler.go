// Package middleware provides middleware functions for error recovery and handling panics.
package middleware

import (
	"api/internal/app/constants"
	"api/internal/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware is a custom middleware that recovers from panics to prevent the application from crashing.
//
// If a panic occurs, it logs the error and returns a structured JSON response with a relevant error message.
// The error message varies based on the environment:
// - In "development" mode, the panic message is included in the response for debugging.
// - In other environments, a generic error message is returned.
//
// Returns:
// - gin.HandlerFunc: A Gin middleware function.
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				var errMsg string
				// if os.Getenv("APP_ENV") == "development" {
				errMsg = "Internal Server Error: " + r.(string)
				// } else {
				// 	errMsg = constants.Const.Default.ServerErrorMessage
				// }

				// Log the panic for debugging purposes
				log.Printf("Panic recovered: %s", r)

				// Ensure response is not already written before modifying headers
				if !c.Writer.Written() {
					c.Writer.Header().Set(constants.Const.Headers.ErrorCode, utils.ConvertUINTtoString(http.StatusInternalServerError))
					c.Writer.Header().Set(constants.Const.Headers.ErrorMessage, errMsg)
				}

				// Return JSON response with error message
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errMsg})
				c.Abort()
				return
			}
		}()
		// Proceed with the next middleware/handler
		c.Next()
	}
}
