// Package auth provides authentication-related utilities such as JWT management
// and user authentication claims handling.
package auth

import (
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JwtKey is the secret key used for signing JWT access tokens.
var JwtKey []byte

// JwtRefreshKey is the secret key used for signing JWT refresh tokens.
var JwtRefreshKey []byte

// accessTokenExpiry defines the duration for which an access token remains valid.
var accessTokenExpiry time.Duration

// refreshTokenExpiry defines the duration for which a refresh token remains valid.
var refreshTokenExpiry time.Duration

// mutex ensures thread safety when modifying JWT-related variables.
var mutex sync.Mutex

// Claims represents the JWT payload structure containing essential user information.
type Claims struct {
	Email string `json:"email"` // Email of the authenticated user.
	Sub   uint   `json:"sub"`   // Subject identifier, typically the user's unique ID.
	jwt.RegisteredClaims
}
