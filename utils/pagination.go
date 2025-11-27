package utils

import (
	"math"
	"net/http"
	"strconv"
)

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page     int `json:"page"`
	Limit    int `json:"limit"`
	Offset   int `json:"-"`
	MaxLimit int `json:"-"`
}

// PaginationMeta holds pagination metadata
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// NewPaginationParams creates and validates pagination parameters
func NewPaginationParams(r *http.Request) *PaginationParams {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	maxLimit := 100
	if limit > maxLimit {
		limit = maxLimit
	}

	offset := (page - 1) * limit

	return &PaginationParams{
		Page:     page,
		Limit:    limit,
		Offset:   offset,
		MaxLimit: maxLimit,
	}
}

// CalculateMeta calculates pagination metadata
func (p *PaginationParams) CalculateMeta(total int) *PaginationMeta {
	totalPages := int(math.Ceil(float64(total) / float64(p.Limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	return &PaginationMeta{
		Page:       p.Page,
		Limit:      p.Limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    p.Page < totalPages,
		HasPrev:    p.Page > 1,
	}
}

// PaginatedResponse wraps data with pagination metadata
type PaginatedResponse struct {
	Data       interface{}      `json:"data"`
	Pagination *PaginationMeta `json:"pagination"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data interface{}, pagination *PaginationMeta) *PaginatedResponse {
	return &PaginatedResponse{
		Data:       data,
		Pagination: pagination,
	}
}