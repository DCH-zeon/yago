package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dch-zeon/yago/internal/contextutil"
	"github.com/dch-zeon/yago/internal/errors"
)

// RequireRole returns a middleware that checks if the user has the specified role
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !contextutil.HasRole(c, role) {
			c.JSON(http.StatusForbidden, errors.Forbidden("insufficient permissions"))
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireWholesale returns a middleware that checks if the user is an admin
func RequireWholesale() gin.HandlerFunc {
	return RequireRole("wholesale")
}
