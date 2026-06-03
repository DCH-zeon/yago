package customer

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type txKey struct{}

// Repository defines customer repository interface
type Repository interface {
	Create(ctx context.Context, customer *Customer) error
	FindByPhone(ctx context.Context, phone string) (*Customer, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Customer, error)
	Update(ctx context.Context, user *Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListAllCustomers(ctx context.Context, filters CustomerFilterParams, page, perPage int) ([]Customer, int64, error)
	AssignRole(ctx context.Context, customerID uuid.UUID, roleName string) error
	RemoveRole(ctx context.Context, userID uuid.UUID, roleName string) error
	FindRoleByName(ctx context.Context, name string) (*Role, error)
	GetCustomerRoles(ctx context.Context, userID uuid.UUID) ([]Role, error)
	Transaction(ctx context.Context, fn func(context.Context) error) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new customer repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// getDB returns the DB from context if in transaction, otherwise returns the repository's DB
func (r *repository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return r.db
}

// Create creates a new customer in the database
func (r *repository) Create(ctx context.Context, customer *Customer) error {
	result := r.getDB(ctx).WithContext(ctx).Create(customer)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindByPhone finds a customer by phone
func (r *repository) FindByPhone(ctx context.Context, phone string) (*Customer, error) {
	var customer Customer
	result := r.getDB(ctx).WithContext(ctx).Preload("Roles").Where("phone = ?", phone).First(&customer)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &customer, nil
}

// FindByID finds a customer by ID
func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*Customer, error) {
	var customer Customer
	result := r.getDB(ctx).WithContext(ctx).Preload("Roles").First(&customer, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &customer, nil
}

// Update updates a customer in the database
func (r *repository) Update(ctx context.Context, customer *Customer) error {
	// WHY: Save() syncs associations, potentially clearing roles
	result := r.getDB(ctx).WithContext(ctx).Select("first_name", "phone", "password_hash", "updated_at").Save(customer)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete soft deletes a customer from the database
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.getDB(ctx).WithContext(ctx).Delete(&Customer{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ListAllCustomers retrieves a paginated list of customers with filters
func (r *repository) ListAllCustomers(ctx context.Context, filters CustomerFilterParams, page, perPage int) ([]Customer, int64, error) {
	var customers []Customer
	var total int64

	query := r.getDB(ctx).WithContext(ctx).Model(&Customer{}).Preload("Roles")

	if filters.Role != "" {
		query = query.Joins("JOIN customer_roles ON customer_roles.customer_id = customers.id").
			Joins("JOIN roles ON roles.id = customer_roles.role_id").
			Where("roles.name = ?", filters.Role)
	}

	if filters.Search != "" {
		// WHY: Escape SQL LIKE wildcards to prevent incorrect matches
		escapedSearch := strings.ReplaceAll(filters.Search, "%", "\\%")
		escapedSearch = strings.ReplaceAll(escapedSearch, "_", "\\_")
		searchPattern := "%" + escapedSearch + "%"
		query = query.Where("customers.first_name LIKE ? OR customers.phone LIKE ? OR customers.email LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// WHY: Count distinct customer IDs when using JOINs to avoid inflated totals
	if err := query.Distinct("customers.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage

	// Defense-in-depth: Validate sort parameters at the repository layer
	validSorts := map[string]bool{
		"first_name": true, "phone": true, "created_at": true, "updated_at": true,
	}
	if !validSorts[filters.Sort] {
		return nil, 0, errors.New("invalid sort field")
	}
	if filters.Order != "asc" && filters.Order != "desc" {
		return nil, 0, errors.New("invalid sort order")
	}

	// Use type-safe GORM clause to prevent SQL injection
	orderColumn := clause.OrderByColumn{
		Column: clause.Column{Table: "customers", Name: filters.Sort},
		Desc:   filters.Order == "desc",
	}

	// WHY: Use Distinct with explicit columns to avoid duplicate customers with JOINs
	if err := query.Distinct("customers.*").Order(orderColumn).Limit(perPage).Offset(offset).Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

// AssignRole assigns a role to a customer
func (r *repository) AssignRole(ctx context.Context, customerID uuid.UUID, roleName string) error {
	role, err := r.FindRoleByName(ctx, roleName)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Use database-level conflict handling for race-safe, idempotent role assignment
	// Works with both PostgreSQL and SQLite
	return r.getDB(ctx).WithContext(ctx).Exec(`
		INSERT INTO customer_roles (customer_id, role_id, assigned_at)
		VALUES (?, ?, ?)
		ON CONFLICT (customer_id, role_id) DO NOTHING
	`, customerID, role.ID, time.Now()).Error
}

// RemoveRole removes a role from a user
func (r *repository) RemoveRole(ctx context.Context, customerID uuid.UUID, roleName string) error {
	role, err := r.FindRoleByName(ctx, roleName)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	return r.getDB(ctx).WithContext(ctx).Exec(
		"DELETE FROM customer_roles WHERE customer_id = ? AND role_id = ?",
		customerID, role.ID,
	).Error
}

// FindRoleByName finds a role by name
func (r *repository) FindRoleByName(ctx context.Context, name string) (*Role, error) {
	var role Role
	result := r.getDB(ctx).WithContext(ctx).Where("name = ?", name).First(&role)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &role, nil
}

// GetCustomerRoles retrieves all roles for a customer
func (r *repository) GetCustomerRoles(ctx context.Context, customerID uuid.UUID) ([]Role, error) {
	var roles []Role
	err := r.getDB(ctx).WithContext(ctx).
		Table("roles").
		Joins("JOIN customer_roles ON customer_roles.role_id = roles.id").
		Where("customer_roles.customer_id = ?", customerID).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// Transaction executes a function within a database transaction
func (r *repository) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Inject transaction into context
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}
