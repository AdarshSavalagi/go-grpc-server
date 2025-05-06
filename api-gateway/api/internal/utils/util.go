package utils

import (
	"api/internal/utils/response"
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID (Universally Unique Identifier) and returns it as a string.
func GenerateUUID() string {
	newUUID := uuid.New()
	return newUUID.String()
}

// ConvertUINTtoString converts an unsigned integer (uint) to its string representation.
func ConvertUINTtoString(value uint) string {
	return strconv.FormatUint(uint64(value), 10)
}

// GetUserKey generates a unique key for identifying a user session.
//
// Parameters:
//   - key: A base key (optional).
//   - clientIP: The IP address of the client.
//   - userID: The user ID (optional).
//   - clientID: The client ID (optional).
//
// Returns:
//   - A formatted string key that combines the given parameters to uniquely identify the user.
func GetUserKey(key string, clientIP string, userID interface{}, clientID interface{}) string {
	if key == "" {
		if userID != nil {
			return fmt.Sprintf("%v:%s:%v", userID, clientIP, clientID)
		}
		return fmt.Sprintf("%s:%v", clientIP, clientID)
	}

	if userID != nil {
		return fmt.Sprintf("%s:%v:%s:%v", key, userID, clientIP, clientID)
	}
	return fmt.Sprintf("%s:%s:%v", key, clientIP, clientID)
}

// GetDeviceKey generates a unique key for identifying a device along with user agent information.
//
// Parameters:
//   - deviceId: The unique identifier of the device.
//   - userAgent: The user agent string (optional).
//
// Returns:
//   - A formatted string key that uniquely represents the device.
func GetDeviceKey(deviceId string, userAgent string) string {
	if userAgent != "" {
		return fmt.Sprintf("DeviceInfo:%s:%s", deviceId, userAgent)
	}
	return fmt.Sprintf("DeviceInfo:%s", deviceId)
}

func SetDefaultPagination(p *response.Pagination) {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
}
