package pagination

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var (
	DefaultPageSize = 10
	MaxPageSize     = 100
	// PageVar specifies the query parameter name for page number
	PageVar = "page"
	// LimitVar specifies the query parameter name for page size
	LimitVar = "limit"
)

// Pages represents a paginated list of data items
type Pages struct {
	TotalCount  int `json:"total"`
	Limit       int `json:"perPage"`
	CurrentPage int `json:"currentPage"`
	LastPage    int `json:"lastPage"`
}

// New creates a new Pages instance.
// The page parameter is 1-based and refers to the current page index/number.
// The limit parameter refers to the number of items on each page.
// And the total parameter specifies the total number of data items.
// If total is less than 0, it means total is unknown.
func New(currentPage, limit, total int) *Pages {
	if limit <= 0 {
		limit = DefaultPageSize
	}

	if limit > MaxPageSize {
		limit = MaxPageSize
	}

	lastPage := -1
	if total >= 0 {
		lastPage = (total + limit - 1) / limit
		if currentPage > lastPage {
			currentPage = lastPage
		}
	}

	if currentPage < 1 {
		currentPage = 1
	}

	return &Pages{
		CurrentPage: currentPage,
		Limit:       limit,
		TotalCount:  total,
		LastPage:    lastPage,
	}
}

// NewFromRequest creates a Pages object using the query parameters found in the given HTTP request.
// count stands for the total number of items. Use -1 if this is unknown.
func NewFromRequest(c *fiber.Ctx, count int) *Pages {
	currentPage := parseInt(c.Query(PageVar), 1)
	limit := parseInt(c.Query(LimitVar), DefaultPageSize)

	return New(currentPage, limit, count)
}

func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}

	if result, err := strconv.Atoi(value); err != nil {
		return result
	}

	return defaultValue
}

// Offset returns the OFFSET value that can be used in a SQL statement
func (p *Pages) Offset() int {
	return (p.CurrentPage - 1) * p.Limit
}

// Size returns the LIMIT value tahat cannot be used in a SQL statement
func (p *Pages) Size() int {
	return p.Limit
}

type PaginationParams struct {
	Offset int
	Limit  int
}
