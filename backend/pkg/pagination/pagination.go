package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Params holds parsed pagination values from query string.
// Usage in Postman: GET /jobs?page=2&limit=20
type Params struct {
	Page   int
	Limit  int
	Offset int
}

// Parse extracts ?page= and ?limit= from the request.
// Node analogy: const { page = 1, limit = 20 } = req.query
func Parse(c *gin.Context) Params {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return Params{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}

// Meta is returned in API responses for list endpoints.
// Node analogy: { total, page, limit, totalPages } in your list response.
type Meta struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"totalPages"`
}

func NewMeta(total int64, params Params) Meta {
	totalPages := int(total) / params.Limit
	if int(total)%params.Limit != 0 {
		totalPages++
	}
	return Meta{
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}
}
