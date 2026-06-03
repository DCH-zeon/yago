package tests

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/dch-zeon/yago/internal/auth"
	"github.com/dch-zeon/yago/internal/contextutil"
)

func TestGetCustomer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns customer claims when present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		claims := &auth.Claims{CustomerID: uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1cc4r")), Phone: "+380997775533"}
		c.Set(auth.KeyCustomer, claims)

		result := contextutil.GetCustomer(c)
		assert.NotNil(t, result)
		assert.Equal(t, uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1cc4r")), result.CustomerID)
		assert.Equal(t, "+380997775533", result.Phone)
	})

	t.Run("returns nil when not present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)

		result := contextutil.GetCustomer(c)
		assert.Nil(t, result)
	})

	t.Run("returns nil when wrong type", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Set(auth.KeyCustomer, "not-a-claims-struct")

		result := contextutil.GetCustomer(c)
		assert.Nil(t, result)
	})
}

func TestMustGetCustomer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns customer claims when present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		claims := &auth.Claims{CustomerID: uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1cc4r")), Phone: "+380997775533"}
		c.Set(auth.KeyCustomer, claims)

		result, err := contextutil.MustGetCustomer(c)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1cc4r")), result.CustomerID)
		assert.Equal(t, "+380997775533", result.Phone)
	})

	t.Run("returns error when not present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)

		result, err := contextutil.MustGetCustomer(c)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "customer not found in context")
	})
}

func TestGetCustomerID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns customer ID when present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		claims := &auth.Claims{CustomerID: uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1cc4r")), Phone: "+380997775533"}
		c.Set(auth.KeyCustomer, claims)

		customerID := contextutil.GetCustomerID(c)
		assert.Equal(t, uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1cc4r")), customerID)
	})

	t.Run("returns 0 when not present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)

		customerID := contextutil.GetCustomerID(c)
		assert.Equal(t, uuid.Nil, customerID)
	})

	t.Run("returns 0 when wrong type", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Set(auth.KeyCustomer, "not-a-claims-struct")

		customerID := contextutil.GetCustomerID(c)
		assert.Equal(t, uuid.Nil, customerID)
	})
}

func TestMustGetCustomerID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns customer ID when present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		claims := &auth.Claims{CustomerID: uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1cc4r")), Phone: "+380997775533"}
		c.Set(auth.KeyCustomer, claims)

		customerID, err := contextutil.MustGetCustomerID(c)
		assert.NoError(t, err)
		assert.Equal(t, uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1cc4r")), customerID)
	})

	t.Run("returns error when not present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)

		customerID, err := contextutil.MustGetCustomerID(c)
		assert.Error(t, err)
		assert.Equal(t, uint(0), customerID)
		assert.Contains(t, err.Error(), "customer ID not found in context")
	})

	t.Run("returns error when customer ID is nil", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		claims := &auth.Claims{CustomerID: uuid.Nil, Phone: "+380997775533"}
		c.Set(auth.KeyCustomer, claims)

		customerID, err := contextutil.MustGetCustomerID(c)
		assert.Error(t, err)
		assert.Equal(t, uint(0), customerID)
		assert.Contains(t, err.Error(), "customer ID not found in context")
	})
}

func TestGetPhone(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns email when present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		claims := &auth.Claims{CustomerID: uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1cc4r")), Phone: "+380997775533"}
		c.Set(auth.KeyCustomer, claims)

		email := contextutil.GetPhone(c)
		assert.Equal(t, "+380997775533", email)
	})

	t.Run("returns empty string when not present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)

		email := contextutil.GetPhone(c)
		assert.Equal(t, "", email)
	})

	t.Run("returns empty string when wrong type", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Set(auth.KeyCustomer, "not-a-claims-struct")

		email := contextutil.GetPhone(c)
		assert.Equal(t, "", email)
	})
}

func TestIsAuthenticated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns true when customer is present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		claims := &auth.Claims{CustomerID: uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1cc4r")), Phone: "+380997775533"}
		c.Set(auth.KeyCustomer, claims)

		authenticated := contextutil.IsAuthenticated(c)
		assert.True(t, authenticated)
	})

	t.Run("returns false when customer is not present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)

		authenticated := contextutil.IsAuthenticated(c)
		assert.False(t, authenticated)
	})

	t.Run("returns false when wrong type", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Set(auth.KeyCustomer, "not-a-claims-struct")

		authenticated := contextutil.IsAuthenticated(c)
		assert.False(t, authenticated)
	})
}

func TestCanAccessCustomer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns true when customer can access own resource", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		claims := &auth.Claims{CustomerID: uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1cc4r")), Phone: "+380997775533"}
		c.Set(auth.KeyCustomer, claims)

		canAccess := contextutil.CanAccessCustomer(c, uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1cc4r")))
		assert.True(t, canAccess)
	})

	t.Run("returns false when customer cannot access other resource", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		claims := &auth.Claims{CustomerID: uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1cc4r")), Phone: "+380997775533"}
		c.Set(auth.KeyCustomer, claims)

		canAccess := contextutil.CanAccessCustomer(c, uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1cc4r")))
		assert.False(t, canAccess)
	})

	t.Run("returns false when customer is not authenticated", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)

		canAccess := contextutil.CanAccessCustomer(c, uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1cc4r")))
		assert.False(t, canAccess)
	})
}

func TestGetCustomerName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns customer name when present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		claims := &auth.Claims{CustomerID: uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1cc4r")), Phone: "+380997775533", Name: "John Doe"}
		c.Set(auth.KeyCustomer, claims)

		customerName := contextutil.GetCustomerName(c)
		assert.Equal(t, "John Doe", customerName)
	})

	t.Run("returns empty string when not present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)

		customerName := contextutil.GetCustomerName(c)
		assert.Equal(t, "", customerName)
	})

	t.Run("returns empty string when wrong type", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Set(auth.KeyCustomer, "not-a-claims-struct")

		customerName := contextutil.GetCustomerName(c)
		assert.Equal(t, "", customerName)
	})
}

func TestHasRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns false when customer is present (roles not implemented)", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		claims := &auth.Claims{CustomerID: uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1cc4r")), Phone: "+380997775533", Name: "John Doe"}
		c.Set(auth.KeyCustomer, claims)

		hasRole := contextutil.HasRole(c, "regular")
		assert.False(t, hasRole) // Currently returns false as roles are not implemented
	})

	t.Run("returns false when customer is not present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)

		hasRole := contextutil.HasRole(c, "regular")
		assert.False(t, hasRole)
	})

	t.Run("returns false when wrong type", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Set(auth.KeyCustomer, "not-a-claims-struct")

		hasRole := contextutil.HasRole(c, "regular")
		assert.False(t, hasRole)
	})
}
