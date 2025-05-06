// Package middleware provides various middleware functions for logging requests and responses.
package middleware

import (
	"api/internal/app/constants"
	"api/internal/app/logger"
	"bytes"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggingMiddleware logs incoming HTTP requests and their responses.
//
// It captures request details (method, URL, client IP, request ID, client ID),
// logs the request before processing, and then logs the response along with its duration.
//
// Returns:
// - gin.HandlerFunc: A Gin middleware function.
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize logger
		init := logger.InitLogger()
		startTime := time.Now()

		// Log request details
		init.WithFields(logrus.Fields{
			"Method":    c.Request.Method,
			"URL":       c.Request.URL.String(),
			"ClientIP":  c.ClientIP(),
			"RequestID": c.GetString(constants.Const.Context.RequestID),
			"ClientID":  c.GetString(constants.Const.Context.ClientID),
		}).Info("Request received")

		// Capture response body
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Continue processing request
		c.Next()

		// Measure response time
		responseTime := time.Since(startTime)

		// Prepare log fields
		fields := logrus.Fields{
			"Method":       c.Request.Method,
			"URL":          c.Request.URL.String(),
			"HttpStatus":   c.Writer.Status(),
			"ResponseTime": responseTime.Seconds(),
			"RequestID":    c.GetString(constants.Const.Context.RequestID),
			"ClientID":     c.GetString(constants.Const.Context.ClientID),
		}

		// Handle special cases before logging response body
		if c.FullPath() == "/metrics" {
			fields["response"] = "[metrics data]"
		} else {
			// Determine whether to log response body based on content type and encoding
			contentType := c.Writer.Header().Get(constants.Const.Headers.ContentType)
			contentEncoding := c.Writer.Header().Get(constants.Const.Headers.ContentEncoding)
			if contentEncoding != "gzip" && (contentType == "application/json" || contentType == "text/plain") {
				fields["response"] = blw.body.String()
			} else {
				fields["response"] = "[binary data or compressed data]"
			}
		}

		// Log response details
		init.WithFields(fields).Info("Response sent")
	}
}

// bodyLogWriter is a wrapper for gin.ResponseWriter that captures the response body.
//
// This helps in logging the response data without modifying the actual response flow.
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}
