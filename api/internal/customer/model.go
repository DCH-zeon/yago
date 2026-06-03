package customer

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Customer represents a customer in the system
type Customer struct {
	ID           uuid.UUID      `Gorm:"primaryKey" json:"id"`
	FirstName    string         `json:"first_name,omitempty"`
	LastName     string         `json:"last_name,omitempty"`
	Email        string         `json:"email,omitempty"`
	Phone        string         `Gorm:"uniqueIndex;not null" json:"phone"`
	PasswordHash string         `Gorm:"not null" json:"-"`
	Roles        []Role         `Gorm:"many2many:customer_roles;" json:"-"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `Gorm:"index" json:"-"`
}

// TableName specifies the table name for a Customer model
func (c *Customer) TableName() string {
	return "customers"
}

// HasRole checks if a customer has a specific role
func (c *Customer) HasRole(roleName string) bool {
	for _, role := range c.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

// IsWholesale checks if a customer has a wholesale role
func (c *Customer) IsWholesale() bool {
	return c.HasRole(RoleWholesale)
}

// GetRoleNames returns a list of role names
func (c *Customer) GetRoleNames() []string {
	roleNames := make([]string, len(c.Roles))
	for i, role := range c.Roles {
		roleNames[i] = role.Name
	}
	return roleNames
}
