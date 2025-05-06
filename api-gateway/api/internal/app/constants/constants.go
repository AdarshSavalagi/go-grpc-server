// Description: Defines standard HTTP headers, cookie names, context keys, and messages used in the application.
// This package initializes constants used throughout the application by reading struct tags.
// The struct tags define the constant values for each field.
// This approach allows for easy modification of constants without changing the code.
// The constants are used in API responses, error messages, and HTTP headers.
// The package initializes the constants by reading struct tags and setting the field values.
package constants

import "reflect"

// Headers defines standard HTTP headers used in the application.
type Headers struct {
	Accept                    string `constant:"Accept"`
	AcceptJSON                string `constant:"application/json"`
	AcceptXML                 string `constant:"application/xml"`
	ContentType               string `constant:"Content-Type"`
	ContentTypeJSON           string `constant:"application/json"`
	ContentTypeXML            string `constant:"application/xml"`
	ContentEncoding           string `constant:"Content-Encoding"`
	Idempotency               string `constant:"Idempotency-Key"`
	RequestID                 string `constant:"X-Request-ID"`
	UserAgent                 string `constant:"User-Agent"`
	RefreshToken              string `constant:"X-Refresh-Token"`
	AuthToken                 string `constant:"Authorization"`
	XFF                       string `constant:"X-Forwarded-For"`
	RealIP                    string `constant:"X-Real-IP"`
	MobileClientID            string `constant:"X-Device-ID"`
	WebClientID               string `constant:"X-Client-ID"`
	CacheControl              string `constant:"Cache-Control"`
	CacheControlNoneMatch     string `constant:"If-None-Match"`
	CacheControlETag          string `constant:"ETag"`
	CacheControlLastModified  string `constant:"Last-Modified"`
	CacheControlModifiedSince string `constant:"If-Modified-Since"`
	ErrorCode                 string `constant:"X-Error-Code"`
	ErrorMessage              string `constant:"X-Error-Message"`
	Referer                   string `constant:"Referer"`
	CacheContent              string `constant:"X-Cache-Content"`
	TotalRateLimit            string `constant:"X-RateLimit-Limit"`
	RateLimitRemaining        string `constant:"X-RateLimit-Remaining"`
	RateLimitResetIn          string `constant:"X-RateLimit-Reset"`
}

// Cookies defines constants for cookie names.
type Cookies struct {
	WebClientID string `constant:"client_id"`
}

// Context defines keys for storing request-specific values in the Gin context.
type Context struct {
	UserID        string `constant:"UserID"`
	Username      string `constant:"Username"`
	ClientID      string `constant:"ClientID"`
	UserAgent     string `constant:"UserAgent"`
	CachedContent string `constant:"CachedContent"`
	RequestID     string `constant:"RequestID"`
	DeviceInfo    string `constant:"DeviceInfo"`
	DeviceID      string `constant:"DeviceID"`
	IdempotencyID string `constant:"IdempotencyID"`
}

// Messages defines standard messages used in API responses.
type Messages struct {
	ResourceCreated string `constant:"Resource Created Successfully."`
	ResourceDeleted string `constant:"Resource Deleted"`
	ResourceUpdated string `constant:"Resource Updated"`
}

// Default defines standard error messages and codes.
type Default struct {
	ServerErrorCode             string `constant:"500"`
	ServerErrorMessage          string `constant:"Sorry, the server can't handle the request at this moment."`
	ValidationErrorCode         string `constant:"400"`
	ValidationErrorMessage      string `constant:"Validation Error"`
	UnauthorizedResponseMessage string `constant:"Unauthorized: Authentication is required and has failed or has not yet been provided."`
	ForbiddenResponseMessage    string `constant:"Forbidden: You do not have permission to access this resource."`
}

// Constants aggregates all the defined constant structures.
type Constants struct {
	Headers  Headers
	Cookies  Cookies
	Context  Context
	Messages Messages
	Default  Default
}

// Const holds the initialized constants used throughout the application.
var Const = &Constants{
	Headers:  Headers{},
	Cookies:  Cookies{},
	Context:  Context{},
	Messages: Messages{},
	Default:  Default{},
}

// init initializes constant structures by setting values from struct tags.
func init() {
	initStructFromTags(&Const.Headers)
	initStructFromTags(&Const.Cookies)
	initStructFromTags(&Const.Context)
	initStructFromTags(&Const.Messages)
	initStructFromTags(&Const.Default)
}

// initStructFromTags sets struct field values using their "constant" tag values.
//
// This function reflects over the struct fields, reads the tag values,
// and assigns them as field values.
//
// Parameters:
//   - s: A pointer to the struct to initialize.
func initStructFromTags(s interface{}) {
	val := reflect.ValueOf(s).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		tag := val.Type().Field(i).Tag.Get("constant")
		field.SetString(tag)
	}
}
