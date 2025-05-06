package response

import (
	"api/internal/app/constants"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Response creates a generic response structure.
//
// Parameters:
//   - flag: A boolean indicating success (true) or failure (false).
//   - message: A descriptive message.
//   - requestID: The unique identifier for the request.
//
// Returns:
//   - A pointer to TResponse containing response details.
func Response(flag bool, message string, requestID string) *TResponse {
	return &TResponse{
		Flag:      flag,
		Message:   message,
		RequestID: requestID,
	}
}

// SuccessResponse creates a success response with data.
//
// Parameters:
//   - data: The data to include in the response.
//   - requestID: The unique identifier for the request.
//
// Returns:
//   - A pointer to TSuccessResponse containing the response details.
func SuccessResponse(data interface{}, requestID string) *TSuccessResponse {
	return &TSuccessResponse{
		Flag:      true,
		Data:      data,
		RequestID: requestID,
	}
}

// ErrorResponse creates an error response with details.
//
// Parameters:
//   - message: The error message.
//   - errors: A slice of error messages.
//   - trace: A slice of error traces (for debugging).
//   - requestID: The unique identifier for the request.
//
// Returns:
//   - A pointer to TErrorResponse containing error details.
func ErrorResponse(message string, errors []string, trace []string, requestID string) *TErrorResponse {
	return &TErrorResponse{
		Flag:      false,
		Message:   message,
		Errors:    errors,
		Trace:     trace,
		RequestID: requestID,
	}
}

// RespondError sends an error response with appropriate status code and error details.
//
// Parameters:
//   - c: The Gin context.
//   - statusCode: The HTTP status code.
//   - message: The error message.
//   - errors: The error object.
//
// Behavior:
//   - If validation errors exist, they are formatted accordingly.
//   - In production, internal server errors do not expose details.
func RespondError(c *gin.Context, statusCode int, message string, errors error) {
	requestID := c.GetString(constants.Const.Context.RequestID)

	var errorMessages []string
	var errorTraces []string
	validationError := false

	if errors != nil {
		if os.Getenv("APP_ENV") == "development" {
			errorTraces = append(errorTraces, errors.Error())
			errorMessages = append(errorMessages, errors.Error())
		}

		if ve, ok := errors.(validator.ValidationErrors); ok {
			validationError = true
			for _, fieldErr := range ve {
				errorMessages = append(errorMessages, fieldErr.Error())
			}
		}

		if statusCode == http.StatusInternalServerError && os.Getenv("APP_ENV") != "development" {
			message = constants.Const.Default.ServerErrorMessage
		}
	}

	response := ErrorResponse(message, errorMessages, errorTraces, requestID)

	if errors == nil {
		response.Errors = []string{}
	}

	if statusCode == http.StatusInternalServerError && os.Getenv("APP_ENV") != "development" {
		response.Message = constants.Const.Default.ServerErrorMessage
		response.Errors = nil
		response.Trace = nil
	}

	if validationError {
		statusCode = http.StatusBadRequest
		response.Message = constants.Const.Default.ValidationErrorMessage
	}

	c.JSON(statusCode, response)
}

// Respond sends a generic response with a status, status code, and message.
//
// Parameters:
//   - c: The Gin context.
//   - status: A boolean indicating success (true) or failure (false).
//   - statusCode: The HTTP status code.
//   - message: The response message.
func Respond(c *gin.Context, status bool, statusCode int, message string) {
	requestID := c.GetString(constants.Const.Context.RequestID)
	response := Response(status, message, requestID)
	c.JSON(statusCode, response)
}

// RespondUnauthorized sends a predefined 401 Unauthorized response.
//
// Parameters:
//   - c: The Gin context.
func RespondUnauthorized(c *gin.Context) {
	Respond(c, false, http.StatusUnauthorized, constants.Const.Default.UnauthorizedResponseMessage)
}

// RespondSuccessWithData sends a success response with data.
//
// Parameters:
//   - c: The Gin context.
//   - data: The data to be returned in the response.
func RespondSuccessWithData(c *gin.Context, data interface{}) {
	requestId := c.GetString(constants.Const.Context.RequestID)
	response := SuccessResponse(data, requestId)
	c.JSON(http.StatusOK, response)
}

// SuccessPaginationResponse returns a successful response with paginated data.
//
// Parameters:
//   - data: The actual data payload
//   - pagination: Pagination metadata (page, limit, total, pages)
//   - requestID: Unique ID for this request
//
// Returns:
//   - A paginated success response object
func SuccessPaginationResponse(data interface{}, pagination PaginationMeta, requestID string) *TPaginatedSuccessResponse {
	return &TPaginatedSuccessResponse{
		Flag:       true,
		Data:       data,
		Pagination: pagination,
		RequestID:  requestID,
	}
}

// RespondWithPagination writes a JSON response with paginated data.
//
// Parameters:
//   - c: Gin context
//   - status: success/failure flag (mostly `true` here)
//   - message: optional message (currently unused but kept for future flexibility)
//   - data: actual paginated list
//   - pagination: pagination object containing meta info
func RespondWithPagination(c *gin.Context, status bool, message string, data interface{}, pagination *Pagination) {
	requestID := c.GetString(constants.Const.Context.RequestID)

	meta := PaginationMeta{
		Page:  pagination.Page,
		Limit: pagination.Limit,
		Total: pagination.Total,
		Pages: pagination.Pages,
	}

	response := SuccessPaginationResponse(data, meta, requestID)

	c.JSON(http.StatusOK, response)
}
