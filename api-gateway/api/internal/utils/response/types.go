package response

// TResponse represents a basic API response without data.
type TResponse struct {
	Flag      bool   `json:"flag"`                // Indicates success or failure
	Message   string `json:"message"`             // Descriptive message
	RequestID string `json:"request_id,omitempty"` // Unique request identifier
}

// TSuccessResponse is used for successful responses with data.
type TSuccessResponse struct {
	Flag      bool        `json:"flag"`                 // Indicates success
	Message   string      `json:"message"`              // Descriptive success message
	Data      interface{} `json:"data,omitempty"`       // Actual payload
	RequestID string      `json:"request_id,omitempty"` // Unique request ID
}

// TErrorResponse represents errors in a consistent format.
type TErrorResponse struct {
	Flag      bool     `json:"flag"`                 // Indicates failure
	Message   string   `json:"message"`              // Error message
	Errors    []string `json:"errors,omitempty"`     // Validation or application errors
	Trace     []string `json:"trace,omitempty"`      // Stack trace or debug info
	RequestID string   `json:"request_id,omitempty"` // Request ID
}

// PaginationMeta holds metadata about paginated results.
type PaginationMeta struct {
	Page  int   `json:"page"`  // Current page number
	Limit int   `json:"limit"` // Items per page
	Total int64 `json:"total"` // Total items
	Pages int   `json:"pages"` // Total pages
}

// TPaginatedSuccessResponse formats paginated data with metadata.
type TPaginatedSuccessResponse struct {
	Flag       bool           `json:"flag"`       // Success indicator
	Data       interface{}    `json:"data"`       // Result list
	Pagination PaginationMeta `json:"pagination"` // Pagination metadata
	RequestID  string         `json:"request_id"` // Request ID
}

// Pagination represents incoming query params and metadata for response.
type Pagination struct {
	Page  int   `json:"-" form:"page,default=1"`   // Incoming page number (from query param)
	Limit int   `json:"-" form:"limit,default=10"` // Incoming limit per page
	Total int64 `json:"-"`                         // Total number of records
	Pages int   `json:"-"`                         // Total number of pages (calculated)
}