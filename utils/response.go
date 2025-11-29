package utils

// Response represents a standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
	Meta    *MetaData   `json:"meta,omitempty"`
}

// ErrorDetail represents detailed error information
type ErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
	Detail  string `json:"detail,omitempty"`
}

// MetaData represents metadata for responses (pagination, etc.)
type MetaData struct {
	Page       int `json:"page,omitempty"`
	Limit      int `json:"limit,omitempty"`
	TotalItems int `json:"total_items,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

// SuccessResponse creates a successful response
func SuccessResponse(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

// SuccessResponseWithMeta creates a successful response with metadata
func SuccessResponseWithMeta(data interface{}, meta *MetaData) Response {
	return Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	}
}

// ErrorResponse creates an error response
type ErrorResponse struct {
	Success bool         `json:"success"`
	Error   *ErrorDetail `json:"error"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    MetaData    `json:"meta"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data interface{}, page, limit, totalItems int) PaginatedResponse {
	totalPages := (totalItems + limit - 1) / limit // Ceiling division
	if totalPages < 1 {
		totalPages = 1
	}

	return PaginatedResponse{
		Success: true,
		Data:    data,
		Meta: MetaData{
			Page:       page,
			Limit:      limit,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version,omitempty"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks,omitempty"`
}