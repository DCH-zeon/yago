package contextutil

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dch-zeon/yago/internal/auth"
)

// GetCustomer retrieves the authenticated customer claims from context
// Returns nil if not found or invalid type
func GetCustomer(c *gin.Context) *auth.Claims {
	value, exists := c.Get(auth.KeyCustomer)
	if !exists {
		return nil
	}

	claims, ok := value.(*auth.Claims)
	if !ok {
		return nil
	}

	return claims
}

// MustGetCustomer retrieves customer claims or returns error
func MustGetCustomer(c *gin.Context) (*auth.Claims, error) {
	claims := GetCustomer(c)
	if claims == nil {
		return nil, fmt.Errorf("customer not found in context")
	}
	return claims, nil
}

// GetCustomerID retrieves the authenticated customer's ID from context
// Returns uuid.Nil if not found
func GetCustomerID(c *gin.Context) uuid.UUID {
	claims := GetCustomer(c)
	if claims == nil {
		return uuid.Nil
	}
	return claims.CustomerID
}

// MustGetCustomerID retrieves customer ID or returns error
func MustGetCustomerID(c *gin.Context) (uuid.UUID, error) {
	customerID := GetCustomerID(c)
	if customerID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("customer ID not found in context")
	}
	return customerID, nil
}

// GetPhone retrieves the authenticated customer's email from context
func GetPhone(c *gin.Context) string {
	claims := GetCustomer(c)
	if claims == nil {
		return ""
	}
	return claims.Phone
}

// IsAuthenticated checks if the request has valid authentication
func IsAuthenticated(c *gin.Context) bool {
	return GetCustomer(c) != nil
}

// CanAccessCustomer checks if an authenticated customer can access the target customer
func CanAccessCustomer(c *gin.Context, targetCustomerID uuid.UUID) bool {
	if IsWholesale(c) {
		return true
	}
	authenticatedCustomerID := GetCustomerID(c)
	return authenticatedCustomerID == targetCustomerID
}

// GetCustomerName retrieves the authenticated customer's name from context
func GetCustomerName(c *gin.Context) string {
	claims := GetCustomer(c)
	if claims == nil {
		return ""
	}
	return claims.Name
}

// HasRole checks if a customer has a specific role
func HasRole(c *gin.Context, role string) bool {
	claims := GetCustomer(c)
	if claims == nil {
		return false
	}
	for _, r := range claims.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// GetRoles retrieves customer roles from context
func GetRoles(c *gin.Context) []string {
	claims := GetCustomer(c)
	if claims == nil {
		return []string{}
	}
	return claims.Roles
}

// IsWholesale checks if a customer has a wholesale role
func IsWholesale(c *gin.Context) bool {
	return HasRole(c, "wholesale")
}
