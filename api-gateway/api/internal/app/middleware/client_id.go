// Package middleware provides middleware functions for handling client identification.
package middleware

import (
	"api/internal/app/constants"
	"api/internal/utils"
	"api/internal/utils/auth"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// getClientIP extracts the client's IP address from request headers.
func getClientIP(c *gin.Context) string {
	xff := c.GetHeader(constants.Const.Headers.XFF)
	if xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	realIP := c.GetHeader(constants.Const.Headers.RealIP)
	if realIP != "" {
		return realIP
	}
	return c.ClientIP()
}

// ClientIDMiddleware assigns a unique client ID for web clients.
//
// This middleware retrieves an existing client ID from cookies or headers.
// If no ID is found, a new one is generated and stored in a cookie and response header.
//
// Returns:
// - gin.HandlerFunc: A middleware function for managing client IDs.
func ClientIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var clientID string
		var err error

		clientIP := getClientIP(c)
		userAgent := c.GetString(constants.Const.Context.UserAgent)
		userID := auth.GetUserIDFromToken(c)
		clientID, err = c.Cookie(constants.Const.Cookies.WebClientID)

		if err != nil || clientID == "" {
			clientID = c.GetHeader(constants.Const.Headers.WebClientID)
			if clientID == "" {
				// Generate a new client ID
				clientID = utils.GenerateUUID()
				// Set the client ID as a cookie if the request has a User-Agent (browser request)
				if userAgent != "" {
					c.SetCookie(constants.Const.Cookies.WebClientID, clientID, int(24*time.Hour.Seconds()), "/", "", false, true)
				}
				// Include the client ID in the response headers
				c.Header(constants.Const.Headers.WebClientID, clientID)
				log.Printf("Generated new client ID: %s", clientID)
			} else {
				log.Printf("Received client ID from header: %s", clientID)
			}
		} else {
			log.Printf("Found existing client ID in cookie: %s", clientID)
		}

		// Fallback to IP if client ID is not available
		if clientID == "" {
			clientID = clientIP
		}

		// Store the client ID in the request context
		c.Set(constants.Const.Context.ClientID, clientID)

		// Generate device identifiers
		deviceID := utils.GetUserKey("", clientIP, userID, clientID)
		deviceKey := utils.GetDeviceKey(deviceID, userAgent)
		c.Set(constants.Const.Context.DeviceID, deviceID)
		c.Set(constants.Const.Context.DeviceInfo, deviceKey)

		// Proceed with the request
		c.Next()
	}
}

// MobileClientIDMiddleware manages client IDs for mobile applications.
//
// This middleware retrieves the client ID from headers or generates a new one if missing.
//
// Returns:
// - gin.HandlerFunc: A middleware function for handling mobile client IDs.
func MobileClientIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.GetHeader(constants.Const.Headers.MobileClientID)
		if clientID == "" {
			clientID = utils.GenerateUUID()
			log.Printf("Generated new client ID for mobile: %s", clientID)
		} else {
			log.Printf("Found existing client ID for mobile: %s", clientID)
		}

		// Store the client ID in the request context
		c.Set(constants.Const.Context.ClientID, clientID)

		// Proceed with the request
		c.Next()
	}
}
