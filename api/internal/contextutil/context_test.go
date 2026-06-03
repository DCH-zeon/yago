package contextutil

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/dch-zeon/yago/internal/auth"
)

func TestGetCustomer(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*gin.Context)
		expected *auth.Claims
	}{
		{
			name: "successful get user",
			setup: func(c *gin.Context) {
				claims := &auth.Claims{
					CustomerID: uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
					Phone:      "+380997775533",
					Name:       "Test Customer",
				}
				c.Set(auth.KeyCustomer, claims)
			},
			expected: &auth.Claims{
				CustomerID: uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
				Phone:      "+380997775533",
				Name:       "Test Customer",
			},
		},
		{
			name:     "user not found in context",
			setup:    func(c *gin.Context) {}, // Don't set anything
			expected: nil,
		},
		{
			name: "invalid type in context",
			setup: func(c *gin.Context) {
				c.Set(auth.KeyCustomer, "invalid-type")
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.setup(c)

			result := GetCustomer(c)

			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected.CustomerID, result.CustomerID)
				assert.Equal(t, tt.expected.Phone, result.Phone)
				assert.Equal(t, tt.expected.Name, result.Name)
			}
		})
	}
}

func TestMustGetCustomer(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*gin.Context)
		expectError bool
		expected    *auth.Claims
	}{
		{
			name: "successful get user",
			setup: func(c *gin.Context) {
				claims := &auth.Claims{
					CustomerID: uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
					Phone:      "+380997775533",
					Name:       "Test Customer",
				}
				c.Set(auth.KeyCustomer, claims)
			},
			expectError: false,
			expected: &auth.Claims{
				CustomerID: uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
				Phone:      "+380997775533",
				Name:       "Test Customer",
			},
		},
		{
			name:        "user not found",
			setup:       func(c *gin.Context) {}, // Don't set anything
			expectError: true,
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.setup(c)

			result, err := MustGetCustomer(c)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, "user not found in context", err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected.CustomerID, result.CustomerID)
				assert.Equal(t, tt.expected.Phone, result.Phone)
				assert.Equal(t, tt.expected.Name, result.Name)
			}
		})
	}
}

func TestGetCustomerID(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*gin.Context)
		expected uuid.UUID
	}{
		{
			name: "successful get user ID",
			setup: func(c *gin.Context) {
				claims := &auth.Claims{
					CustomerID: uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1cc4r")),
					Phone:      "+380997775533",
					Name:       "Test Customer",
				}
				c.Set(auth.KeyCustomer, claims)
			},
			expected: uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1cc4r")),
		},
		{
			name:     "customer not found",
			setup:    func(c *gin.Context) {}, // Don't set anything
			expected: uuid.UUID{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.setup(c)

			result := GetCustomerID(c)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMustGetCustomerID(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*gin.Context)
		expectError bool
		expected    uuid.UUID
	}{
		{
			name: "successful get customer ID",
			setup: func(c *gin.Context) {
				claims := &auth.Claims{
					CustomerID: uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1be4d")),
					Phone:      "+380997775533",
					Name:       "Test Customer",
				}
				c.Set(auth.KeyCustomer, claims)
			},
			expectError: false,
			expected:    uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1be4d")),
		},
		{
			name:        "customer not found",
			setup:       func(c *gin.Context) {}, // Don't set anything
			expectError: true,
			expected:    uuid.UUID{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.setup(c)

			result, err := MustGetCustomerID(c)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, uint(0), result)
				assert.Equal(t, "customer ID not found in context", err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestGetPhone(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*gin.Context)
		expected string
	}{
		{
			name: "successful get phone",
			setup: func(c *gin.Context) {
				claims := &auth.Claims{
					CustomerID: uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
					Phone:      "+380997775533",
					Name:       "Test Customer",
				}
				c.Set(auth.KeyCustomer, claims)
			},
			expected: "+380997775533",
		},
		{
			name:     "customer not found",
			setup:    func(c *gin.Context) {}, // Don't set anything
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.setup(c)

			result := GetPhone(c)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsAuthenticated(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*gin.Context)
		expected bool
	}{
		{
			name: "authenticated customer",
			setup: func(c *gin.Context) {
				claims := &auth.Claims{
					CustomerID: uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
					Phone:      "+380997775533",
					Name:       "Test Customer",
				}
				c.Set(auth.KeyCustomer, claims)
			},
			expected: true,
		},
		{
			name:     "unauthenticated customer",
			setup:    func(c *gin.Context) {}, // Don't set anything
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.setup(c)

			result := IsAuthenticated(c)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCanAccessUser(t *testing.T) {
	tests := []struct {
		name             string
		setup            func(*gin.Context)
		targetCustomerID uuid.UUID
		expected         bool
	}{
		{
			name: "can access own customer",
			setup: func(c *gin.Context) {
				claims := &auth.Claims{
					CustomerID: uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
					Phone:      "+380997775533",
					Name:       "Test Customer",
				}
				c.Set(auth.KeyCustomer, claims)
			},
			targetCustomerID: uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
			expected:         true,
		},
		{
			name: "cannot access other customer",
			setup: func(c *gin.Context) {
				claims := &auth.Claims{
					CustomerID: uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
					Phone:      "+380997775533",
					Name:       "Test Customer",
				}
				c.Set(auth.KeyCustomer, claims)
			},
			targetCustomerID: uuid.Must(uuid.Parse("77be2ef9-62ed-4582-b6c8-cf6525f1bd4e")),
			expected:         false,
		},
		{
			name:             "unauthenticated customer",
			setup:            func(c *gin.Context) {}, // Don't set anything
			targetCustomerID: uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
			expected:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.setup(c)

			result := CanAccessCustomer(c, tt.targetCustomerID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetCustomerName(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*gin.Context)
		expected string
	}{
		{
			name: "successful get customer name",
			setup: func(c *gin.Context) {
				claims := &auth.Claims{
					CustomerID: uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
					Phone:      "+380997775533",
					Name:       "Test Customer",
				}
				c.Set(auth.KeyCustomer, claims)
			},
			expected: "Test Customer",
		},
		{
			name:     "customer not found",
			setup:    func(c *gin.Context) {}, // Don't set anything
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.setup(c)

			result := GetCustomerName(c)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHasRole(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*gin.Context)
		role     string
		expected bool
	}{
		{
			name: "customer has wholesale role",
			setup: func(c *gin.Context) {
				claims := &auth.Claims{
					CustomerID: uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
					Phone:      "+380997775533",
					Name:       "Wholesale Customer",
					Roles:      []string{"wholesale", "regular"},
				}
				c.Set(auth.KeyCustomer, claims)
			},
			role:     "wholesale",
			expected: true,
		},
		{
			name: "customer does not have role",
			setup: func(c *gin.Context) {
				claims := &auth.Claims{
					CustomerID: uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
					Phone:      "+380997775533",
					Name:       "Regular Customer",
					Roles:      []string{"regular"},
				}
				c.Set(auth.KeyCustomer, claims)
			},
			role:     "wholesale",
			expected: false,
		},
		{
			name:     "unauthenticated customer",
			setup:    func(c *gin.Context) {},
			role:     "wholesale",
			expected: false,
		},
		{
			name: "customer with no roles",
			setup: func(c *gin.Context) {
				claims := &auth.Claims{
					CustomerID: uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
					Phone:      "+380997775533",
					Name:       "User",
					Roles:      []string{},
				}
				c.Set(auth.KeyCustomer, claims)
			},
			role:     "wholesale",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.setup(c)

			result := HasRole(c, tt.role)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetRoles(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*gin.Context)
		expected []string
	}{
		{
			name: "customer with multiple roles",
			setup: func(c *gin.Context) {
				claims := &auth.Claims{
					CustomerID: uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
					Phone:      "+380997775533",
					Roles:      []string{"regular", "wholesale", "viewer"},
				}
				c.Set(auth.KeyCustomer, claims)
			},
			expected: []string{"regular", "wholesale", "viewer"},
		},
		{
			name: "customer with single role",
			setup: func(c *gin.Context) {
				claims := &auth.Claims{
					CustomerID: uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
					Phone:      "+380997775533",
					Roles:      []string{"regular"},
				}
				c.Set(auth.KeyCustomer, claims)
			},
			expected: []string{"regular"},
		},
		{
			name:     "unauthenticated customer",
			setup:    func(c *gin.Context) {},
			expected: []string{},
		},
		{
			name: "customer with no roles",
			setup: func(c *gin.Context) {
				claims := &auth.Claims{
					CustomerID: uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
					Phone:      "+380997775533",
					Roles:      []string{},
				}
				c.Set(auth.KeyCustomer, claims)
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.setup(c)

			result := GetRoles(c)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsWholesale(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*gin.Context)
		expected bool
	}{
		{
			name: "customer is regular",
			setup: func(c *gin.Context) {
				claims := &auth.Claims{
					CustomerID: uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
					Phone:      "+380997775533",
					Roles:      []string{"regular", "wholesale"},
				}
				c.Set(auth.KeyCustomer, claims)
			},
			expected: true,
		},
		{
			name: "customer is not regular",
			setup: func(c *gin.Context) {
				claims := &auth.Claims{
					CustomerID: uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
					Phone:      "+380997775533",
					Roles:      []string{"wholesale", "viewer"},
				}
				c.Set(auth.KeyCustomer, claims)
			},
			expected: false,
		},
		{
			name:     "unauthenticated customer",
			setup:    func(c *gin.Context) {},
			expected: false,
		},
		{
			name: "customer with no roles",
			setup: func(c *gin.Context) {
				claims := &auth.Claims{
					CustomerID: uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
					Phone:      "+380997775533",
					Roles:      []string{},
				}
				c.Set(auth.KeyCustomer, claims)
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.setup(c)

			result := IsWholesale(c)
			assert.Equal(t, tt.expected, result)
		})
	}
}
