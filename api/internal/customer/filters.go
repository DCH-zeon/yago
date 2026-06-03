package customer

import (
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

// CustomerFilterParams represents filtering parameters for user list
type CustomerFilterParams struct {
	Role   string
	Search string
	Sort   string
	Order  string
}

// ParseCustomerFilters parses and validates user filter parameters from request
func ParseCustomerFilters(c *gin.Context) CustomerFilterParams {
	role := c.Query("role")
	if role != "" && role != RoleRegular && role != RoleWholesale {
		role = ""
	}

	// Sanitize search parameter: limit length and strip dangerous characters
	search := c.Query("search")
	if search != "" {
		// Limit search length to prevent DoS
		if utf8.RuneCountInString(search) > 100 {
			search = string([]rune(search)[:100])
		}
		// Trim whitespace
		search = strings.TrimSpace(search)
	}

	sort := c.DefaultQuery("sort", "created_at")
	validSorts := map[string]bool{
		"first_name": true,
		"phone":      true,
		"created_at": true,
		"updated_at": true,
	}
	if !validSorts[sort] {
		sort = "created_at"
	}

	order := c.DefaultQuery("order", "desc")
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	return CustomerFilterParams{
		Role:   role,
		Search: search,
		Sort:   sort,
		Order:  order,
	}
}
