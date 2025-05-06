package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 🔑 Setters for JWT Keys
func SetJwtKey(key []byte) {
	mutex.Lock()
	JwtKey = key
	mutex.Unlock()
}

func SetJwtRefreshKey(key []byte) {
	mutex.Lock()
	JwtRefreshKey = key
	mutex.Unlock()
}

// 🔑 Getters for JWT Keys
func GetJwtKey() []byte {
	mutex.Lock()
	defer mutex.Unlock()
	return JwtKey
}

func GetJwtRefreshKey() []byte {
	mutex.Lock()
	defer mutex.Unlock()
	return JwtRefreshKey
}

// ⏳ Setters for Expiry Time
func SetAccessTokenExpiry(expiry int) {
	mutex.Lock()
	accessTokenExpiry = time.Duration(expiry) * time.Minute
	mutex.Unlock()
}

func SetRefreshTokenExpiry(expiry int) {
	mutex.Lock()
	refreshTokenExpiry = time.Duration(expiry) * time.Minute
	mutex.Unlock()
}

// ⏳ Getters for Expiry Time
func GetAccessTokenExpiry() time.Duration {
	mutex.Lock()
	defer mutex.Unlock()
	return accessTokenExpiry
}

func GetRefreshTokenExpiry() time.Duration {
	mutex.Lock()
	defer mutex.Unlock()
	return refreshTokenExpiry
}

// 🛠 Generate JWT Tokens
func GenerateJwtToken(email string, userId uint) (string, string, time.Time, time.Time, error) {
	// Ensure secret keys are not empty
	jwtKey := GetJwtKey()
	if len(jwtKey) < 32 {
		return "", "", time.Time{}, time.Time{}, fmt.Errorf("invalid JWT key: must be at least 32 bytes")
	}

	refreshKey := GetJwtRefreshKey()
	if len(refreshKey) < 32 {
		return "", "", time.Time{}, time.Time{}, fmt.Errorf("invalid Refresh JWT key: must be at least 32 bytes")
	}

	// Set token expiration times
	accessExpiration := time.Now().Add(10 * time.Minute)
	refreshExpiration := time.Now().Add(1000 * time.Minute)
	issued := jwt.NewNumericDate(time.Now())

	// Create Access Token Claims
	accessClaims := &Claims{
		Email: email,
		Sub:   userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiration),
			IssuedAt:  issued,
		},
	}

	// Generate Access Token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	access, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Create Refresh Token Claims
	refreshClaims := &Claims{
		Email: email,
		Sub:   userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiration),
			IssuedAt:  issued,
		},
	}

	// Generate Refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refresh, err := refreshToken.SignedString(refreshKey)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return access, refresh, accessExpiration, refreshExpiration, nil
}

// 🛠 Extract Claims from JWT Token
func GetJwtClaims(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return GetJwtKey(), nil
	})

	// ❌ FIX: Return error if token is invalid
	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
func GetJwtRefreshClaims(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return GetJwtRefreshKey(), nil
	})

	// ❌ FIX: Return error if token is invalid
	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

// 🔄 Generate New Access Token from Refresh Token
func GenerateAccessFromRefresh(refreshToken string) (string, error) {
	claims, err := GetJwtRefreshClaims(refreshToken)
	if err != nil {
		return "", err
	}

	// ❌ FIX: If refresh token is expired, return error
	if time.Now().After(claims.ExpiresAt.Time) {
		return "", jwt.ErrTokenExpired
	}

	// ✅ Generate new access token
	access, _, _, _, err := GenerateJwtToken(claims.Email, claims.Sub)
	if err != nil {
		return "", err
	}
	return access, nil
}

// 🔎 Get User ID from JWT Token
func GetUserIDFromToken(c *gin.Context) interface{} {
	authHeader := c.GetHeader("Authorization") // 🔄 Using string instead of constants
	decryptKey := GetJwtKey()

	if authHeader == "" {
		authHeader = c.GetHeader("X-Refresh-Token")
		if authHeader == "" {
			return nil
		}
		decryptKey = GetJwtRefreshKey()
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return decryptKey, nil
	})

	if err != nil || !token.Valid {
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}

	// ✅ Return the User ID from claims
	return uint(claims["sub"].(float64))
}
