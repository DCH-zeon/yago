package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// AuthorizationHeader is the header key for authorization
	AuthorizationHeader = "Authorization"
	// CustomerIDKey is the context key for customer ID
	CustomerIDKey = "customerID"
	// KeyCustomer is the context key for JWT claims
	KeyCustomer = "customer"
)

// AuthMiddleware creates a middleware that validates JWT tokens
func AuthMiddleware(authService Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header required",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set(KeyCustomer, claims)
		c.Next()
	}
}

// GetCustomerIDFromContext extracts customer ID from gin context
func GetCustomerIDFromContext(c *gin.Context) (uint, bool) {
	customerID, exists := c.Get(CustomerIDKey)
	if !exists {
		return 0, false
	}

	id, ok := customerID.(uint)
	return id, ok
}
