package customer

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCustomer_TableName(t *testing.T) {
	customer := Customer{}
	tableName := customer.TableName()

	assert.Equal(t, "customers", tableName)
}

func TestToCustomerResponse_WithDates(t *testing.T) {
	now := time.Now()
	customer := &Customer{
		ID:        uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Phone:     "+380991234567",
		CreatedAt: now,
		UpdatedAt: now,
	}

	response := ToCustomerResponse(customer)

	assert.Equal(t, uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")), response.ID)
	assert.Equal(t, "John Doe", response.FirstName+" "+response.LastName)
	assert.Equal(t, "john@example.com", response.Email)
	assert.Equal(t, "+380991234567", response.Phone)
	assert.NotEmpty(t, response.CreatedAt)
	assert.NotEmpty(t, response.UpdatedAt)
}

func TestCustomer_HasRole(t *testing.T) {
	tests := []struct {
		name     string
		roles    []Role
		roleName string
		expected bool
	}{
		{
			name: "customer has role",
			roles: []Role{
				{Name: "regular"},
				{Name: "wholesale"},
			},
			roleName: "regular",
			expected: true,
		},
		{
			name: "customer does not have role",
			roles: []Role{
				{Name: "wholesale"},
			},
			roleName: "regular",
			expected: false,
		},
		{
			name:     "customer has no roles",
			roles:    []Role{},
			roleName: "regular",
			expected: false,
		},
		{
			name:     "nil roles",
			roles:    nil,
			roleName: "regular",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customer := &Customer{
				ID:    uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
				Roles: tt.roles,
			}
			result := customer.HasRole(tt.roleName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCustomer_IsWholesale(t *testing.T) {
	tests := []struct {
		name     string
		roles    []Role
		expected bool
	}{
		{
			name: "customer is regular",
			roles: []Role{
				{Name: "regular"},
			},
			expected: true,
		},
		{
			name: "customer is not regular",
			roles: []Role{
				{Name: "wholesale"},
			},
			expected: false,
		},
		{
			name:     "customer has no roles",
			roles:    []Role{},
			expected: false,
		},
		{
			name: "customer has multiple roles including admin",
			roles: []Role{
				{Name: "regular"},
				{Name: "wholesale"},
				{Name: "viewer"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customer := &Customer{
				ID:    uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
				Roles: tt.roles,
			}
			result := customer.IsWholesale()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCustomer_GetRoleNames(t *testing.T) {
	tests := []struct {
		name     string
		roles    []Role
		expected []string
	}{
		{
			name: "single role",
			roles: []Role{
				{Name: "regular"},
			},
			expected: []string{"regular"},
		},
		{
			name: "multiple roles",
			roles: []Role{
				{Name: "regular"},
				{Name: "wholesale"},
				{Name: "viewer"},
			},
			expected: []string{"regular", "wholesale", "viewer"},
		},
		{
			name:     "no roles",
			roles:    []Role{},
			expected: []string{},
		},
		{
			name:     "nil roles",
			roles:    nil,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customer := &Customer{
				ID:    uuid.Must(uuid.Parse("66be2ef9-62ed-4582-b6c8-cf6525f1bb3f")),
				Roles: tt.roles,
			}
			result := customer.GetRoleNames()
			assert.Equal(t, tt.expected, result)
		})
	}
}
