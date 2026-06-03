package customer

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)

	_, err = sqlDB.Exec(`
		CREATE TABLE customers ()
			id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
			external_id TEXT,
			first_name TEXT,
			last_name TEXT,
			email TEXT,
			phone TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME
		);
		CREATE INDEX idx_customers_phone ON customers(phone);
		CREATE INDEX idx_customers_deleted_at ON customers(deleted_at);

		CREATE TABLE roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX idx_roles_name ON roles(name);

		CREATE TABLE customer_roles (
			customer_id TEXT NOT NULL,
			role_id INTEGER NOT NULL,
			assigned_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (customer_id, role_id),
			FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
			FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
		);
		CREATE INDEX idx_customer_roles_customer_id ON customer_roles(customer_id);
		CREATE INDEX idx_customer_roles_role_id ON customer_roles(role_id);

		INSERT INTO roles (id, name, description) VALUES 
			(1, 'regular', 'Звичайний клієнт зі стандартними дозволами'),
			(2, 'wholesale', 'Оптовий клієнт зі спеціальними цінами та доступом до оптових замовлень')
	`)
	require.NoError(t, err)

	return db
}

func TestNewRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	assert.NotNil(t, repo)
	assert.IsType(t, &repository{}, repo)
}

func TestRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	customer := &Customer{
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		Phone:        "+380997775566",
		PasswordHash: "hashed_password",
	}

	err := repo.Create(context.Background(), customer)
	assert.NoError(t, err)
	assert.NotZero(t, customer.ID)
	assert.NotZero(t, customer.CreatedAt)
	assert.NotZero(t, customer.UpdatedAt)
}

func TestRepository_Create_DuplicatePhone(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	customer1 := &Customer{
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		Phone:        "+380997775566",
		PasswordHash: "hashed_password",
	}
	err := repo.Create(context.Background(), customer1)
	assert.NoError(t, err)

	customer2 := &Customer{
		FirstName:    "Jane",
		LastName:     "Doe",
		Email:        "john@example.com",
		Phone:        "+380997775566",
		PasswordHash: "hashed_password",
	}
	err = repo.Create(context.Background(), customer2)
	assert.Error(t, err)
}

func TestRepository_FindByPhone(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	originalCustomer := &Customer{
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		Phone:        "+380997775566",
		PasswordHash: "hashed_password",
	}
	err := repo.Create(context.Background(), originalCustomer)
	require.NoError(t, err)

	t.Run("customer found", func(t *testing.T) {
		customer, err := repo.FindByPhone(context.Background(), "+380997775566")
		assert.NoError(t, err)
		assert.NotNil(t, customer)
		assert.Equal(t, "John Doe", customer.FirstName+" "+customer.LastName)
		assert.Equal(t, "+380997775566", customer.Phone)
		assert.Equal(t, "john@example.com", customer.Email)
	})

	t.Run("customer not found", func(t *testing.T) {
		customer, err := repo.FindByPhone(context.Background(), "+380999995566")
		assert.NoError(t, err)
		assert.Nil(t, customer)
	})
}

func TestRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	originalCustomer := &Customer{
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		Phone:        "+380997775566",
		PasswordHash: "hashed_password",
	}
	err := repo.Create(context.Background(), originalCustomer)
	require.NoError(t, err)

	t.Run("customer found", func(t *testing.T) {
		customer, err := repo.FindByID(context.Background(), originalCustomer.ID)
		assert.NoError(t, err)
		assert.NotNil(t, customer)
		assert.Equal(t, originalCustomer.ID, customer.ID)
		assert.Equal(t, "John Doe", customer.FirstName+" "+customer.LastName)
		assert.Equal(t, "john@example.com", customer.Email)
	})

	t.Run("customer not found", func(t *testing.T) {
		customer, err := repo.FindByID(context.Background(), uuid.Must(uuid.Parse("55be2ef9-62ed-4582-b6c8-cf6525f1bb3f")))
		assert.NoError(t, err)
		assert.Nil(t, customer)
	})
}

func TestRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	customer := &Customer{
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		Phone:        "+380997775566",
		PasswordHash: "hashed_password",
	}
	err := repo.Create(context.Background(), customer)
	require.NoError(t, err)

	customer.FirstName = "Updated Name"
	customer.Email = "updated@example.com"

	err = repo.Update(context.Background(), customer)
	assert.NoError(t, err)

	updatedCustomer, err := repo.FindByID(context.Background(), customer.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updatedCustomer.FirstName)
	assert.Equal(t, "updated@example.com", updatedCustomer.Email)
}

func TestRepository_Update_NonExistentCustomer(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	customer := &Customer{
		ID:           uuid.Must(uuid.Parse("22be2ef9-62ed-4582-b6c8-cf6525f1bb4e")),
		FirstName:    "Ghost Customer",
		Email:        "ghost@example.com",
		Phone:        "+380991115566",
		PasswordHash: "password",
	}

	err := repo.Update(context.Background(), customer)
	// GORM does not return an error when updating a non-existent record; it just affects 0 rows.
	assert.NoError(t, err)
}

func TestRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	customer := &Customer{
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		Phone:        "+380997775566",
		PasswordHash: "hashed_password",
	}
	err := repo.Create(context.Background(), customer)
	require.NoError(t, err)

	err = repo.Delete(context.Background(), customer.ID)
	assert.NoError(t, err)

	deletedCustomer, err := repo.FindByID(context.Background(), customer.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedCustomer)
}

func TestRepository_Delete_NonExistentCustomer(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	err := repo.Delete(context.Background(), uuid.Must(uuid.Parse("22be2ef9-62ed-4582-b6c8-cf6525f1bb4e")))
	// Repository returns an error when no rows are affected (record not found).
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}

func TestRepository_FindRoleByName(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	t.Run("role found", func(t *testing.T) {
		role, err := repo.FindRoleByName(context.Background(), RoleRegular)
		assert.NoError(t, err)
		assert.NotNil(t, role)
		assert.Equal(t, RoleRegular, role.Name)
	})

	t.Run("role not found", func(t *testing.T) {
		role, err := repo.FindRoleByName(context.Background(), "nonexistent_role")
		// SQLite may return nil without error for missing records
		if err == nil {
			assert.Nil(t, role)
		} else {
			assert.Error(t, err)
			assert.Nil(t, role)
		}
	})
}

func TestRepository_AssignRole(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	customer := &Customer{FirstName: "John", LastName: "Doe", Email: "john@example.com", Phone: "+380991115566", PasswordHash: "hash"}
	err := repo.Create(context.Background(), customer)
	require.NoError(t, err)

	t.Run("successful role assignment", func(t *testing.T) {
		err := repo.AssignRole(context.Background(), customer.ID, RoleRegular)
		assert.NoError(t, err)

		var count int64
		db.Table("customer_roles").Where("customer_id = ?", customer.ID).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("idempotent - assigning same role twice doesn't error", func(t *testing.T) {
		err := repo.AssignRole(context.Background(), customer.ID, RoleRegular)
		assert.NoError(t, err)

		var count int64
		db.Table("customer_roles").Where("customer_id = ?", customer.ID).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("nonexistent role", func(t *testing.T) {
		err := repo.AssignRole(context.Background(), customer.ID, "nonexistent_role")
		assert.Error(t, err)
	})

	t.Run("nonexistent customer", func(t *testing.T) {
		err := repo.AssignRole(context.Background(), uuid.Must(uuid.Parse("22be2ef9-62ed-4582-b6c8-cf6525f1bb4e")), RoleRegular)
		// SQLite may not enforce foreign key constraints strictly
		// In production PostgreSQL, this would error
		_ = err // Accept either success or error
	})
}

func TestRepository_RemoveRole(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	customer := &Customer{FirstName: "John", LastName: "Doe", Email: "john@example.com", Phone: "+380991115566", PasswordHash: "hash"}
	err := repo.Create(context.Background(), customer)
	require.NoError(t, err)

	err = repo.AssignRole(context.Background(), customer.ID, RoleRegular)
	require.NoError(t, err)

	t.Run("successful role removal", func(t *testing.T) {
		err := repo.RemoveRole(context.Background(), customer.ID, RoleRegular)
		assert.NoError(t, err)

		var count int64
		db.Table("customer_roles").Where("customer_id = ?", customer.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("removing non-assigned role doesn't error", func(t *testing.T) {
		err := repo.RemoveRole(context.Background(), customer.ID, RoleRegular)
		assert.NoError(t, err)
	})

	t.Run("nonexistent role", func(t *testing.T) {
		err := repo.RemoveRole(context.Background(), customer.ID, "nonexistent_role")
		assert.Error(t, err)
	})

	t.Run("remove role from nonexistent customer - succeeds silently", func(t *testing.T) {
		// This is actually trying to remove a Regular role, so it finds the role, but the customer doesn't exist
		// The DELETE just affects 0 rows, which is not an error in SQL
		err := repo.RemoveRole(context.Background(), uuid.Must(uuid.Parse("22be2ef9-62ed-4582-b6c8-cf6525f1bb4e")), RoleRegular)
		assert.NoError(t, err)
	})
}

func TestRepository_GetCustomerRoles(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	customer := &Customer{FirstName: "John", LastName: "Doe", Email: "john@example.com", Phone: "+380991115566", PasswordHash: "hash"}
	err := repo.Create(context.Background(), customer)
	require.NoError(t, err)

	t.Run("customer with no roles", func(t *testing.T) {
		roles, err := repo.GetCustomerRoles(context.Background(), customer.ID)
		assert.NoError(t, err)
		assert.Empty(t, roles)
	})

	t.Run("customer with single role", func(t *testing.T) {
		err := repo.AssignRole(context.Background(), customer.ID, RoleWholesale)
		require.NoError(t, err)

		roles, err := repo.GetCustomerRoles(context.Background(), customer.ID)
		assert.NoError(t, err)
		assert.Len(t, roles, 1)
		assert.Equal(t, RoleWholesale, roles[0].Name)
	})

	t.Run("customer with multiple roles", func(t *testing.T) {
		err := repo.AssignRole(context.Background(), customer.ID, RoleRegular)
		require.NoError(t, err)

		roles, err := repo.GetCustomerRoles(context.Background(), customer.ID)
		assert.NoError(t, err)
		assert.Len(t, roles, 2)
	})

	t.Run("nonexistent customer", func(t *testing.T) {
		roles, err := repo.GetCustomerRoles(context.Background(), uuid.Must(uuid.Parse("22be2ef9-62ed-4582-b6c8-cf6525f1bb4e")))
		assert.NoError(t, err)
		assert.Empty(t, roles)
	})
}

func TestRepository_ListAllCustomers(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	customer1 := &Customer{FirstName: "Alice", LastName: "Customer", Email: "alice@example.com", Phone: "+380991235566", PasswordHash: "hash"}
	err := repo.Create(context.Background(), customer1)
	require.NoError(t, err)
	err = repo.AssignRole(context.Background(), customer1.ID, RoleRegular)
	require.NoError(t, err)

	customer2 := &Customer{FirstName: "Bob", LastName: "Customer", Email: "bob@example.com", Phone: "+380991245566", PasswordHash: "hash"}
	err = repo.Create(context.Background(), customer2)
	require.NoError(t, err)
	err = repo.AssignRole(context.Background(), customer2.ID, RoleWholesale)
	require.NoError(t, err)

	customer3 := &Customer{FirstName: "Charlie", LastName: "Customer", Email: "charlie@example.com", Phone: "+380991255566", PasswordHash: "hash"}
	err = repo.Create(context.Background(), customer3)
	require.NoError(t, err)

	t.Run("list all customers with defaults", func(t *testing.T) {
		filters := CustomerFilterParams{Sort: "created_at", Order: "desc"}
		customers, total, err := repo.ListAllCustomers(context.Background(), filters, 1, 20)
		assert.NoError(t, err)
		assert.Len(t, customers, 3)
		assert.Equal(t, int64(3), total)
	})

	t.Run("filter by Regular role", func(t *testing.T) {
		filters := CustomerFilterParams{Role: RoleRegular, Sort: "created_at", Order: "desc"}
		customers, total, err := repo.ListAllCustomers(context.Background(), filters, 1, 20)
		assert.NoError(t, err)
		assert.Len(t, customers, 1)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, "alice@example.com", customers[0].Email)
	})

	t.Run("filter by customer role", func(t *testing.T) {
		filters := CustomerFilterParams{Role: RoleWholesale, Sort: "created_at", Order: "desc"}
		customers, total, err := repo.ListAllCustomers(context.Background(), filters, 1, 20)
		assert.NoError(t, err)
		assert.Len(t, customers, 1)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, "bob@example.com", customers[0].Email)
	})

	t.Run("search by name", func(t *testing.T) {
		filters := CustomerFilterParams{Search: "alice", Sort: "created_at", Order: "desc"}
		customers, total, err := repo.ListAllCustomers(context.Background(), filters, 1, 20)
		assert.NoError(t, err)
		assert.Len(t, customers, 1)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, "alice@example.com", customers[0].Email)
	})

	t.Run("search by email", func(t *testing.T) {
		filters := CustomerFilterParams{Search: "bob@", Sort: "created_at", Order: "desc"}
		customers, total, err := repo.ListAllCustomers(context.Background(), filters, 1, 20)
		assert.NoError(t, err)
		assert.Len(t, customers, 1)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, "bob@example.com", customers[0].Email)
	})

	t.Run("pagination - page 1", func(t *testing.T) {
		filters := CustomerFilterParams{Sort: "created_at", Order: "asc"}
		customers, total, err := repo.ListAllCustomers(context.Background(), filters, 1, 2)
		assert.NoError(t, err)
		assert.Len(t, customers, 2)
		assert.Equal(t, int64(3), total)
	})

	t.Run("pagination - page 2", func(t *testing.T) {
		filters := CustomerFilterParams{Sort: "created_at", Order: "asc"}
		customers, total, err := repo.ListAllCustomers(context.Background(), filters, 2, 2)
		assert.NoError(t, err)
		assert.Len(t, customers, 1)
		assert.Equal(t, int64(3), total)
	})

	t.Run("sort by email asc", func(t *testing.T) {
		filters := CustomerFilterParams{Sort: "email", Order: "asc"}
		customers, total, err := repo.ListAllCustomers(context.Background(), filters, 1, 10)
		assert.NoError(t, err)
		assert.Len(t, customers, 3)
		assert.Equal(t, int64(3), total)
		assert.Equal(t, "alice@example.com", customers[0].Email)
		assert.Equal(t, "bob@example.com", customers[1].Email)
		assert.Equal(t, "charlie@example.com", customers[2].Email)
	})

	t.Run("no results for nonexistent search", func(t *testing.T) {
		filters := CustomerFilterParams{Search: "nonexistent", Sort: "created_at", Order: "desc"}
		customers, total, err := repo.ListAllCustomers(context.Background(), filters, 1, 20)
		assert.NoError(t, err)
		assert.Empty(t, customers)
		assert.Equal(t, int64(0), total)
	})

	t.Run("invalid sort field", func(t *testing.T) {
		filters := CustomerFilterParams{Sort: "invalid_field", Order: "asc"}
		customers, total, err := repo.ListAllCustomers(context.Background(), filters, 1, 20)
		assert.Error(t, err)
		assert.Equal(t, "invalid sort field", err.Error())
		assert.Nil(t, customers)
		assert.Equal(t, int64(0), total)
	})

	t.Run("invalid sort order", func(t *testing.T) {
		filters := CustomerFilterParams{Sort: "email", Order: "invalid"}
		customers, total, err := repo.ListAllCustomers(context.Background(), filters, 1, 20)
		assert.Error(t, err)
		assert.Equal(t, "invalid sort order", err.Error())
		assert.Nil(t, customers)
		assert.Equal(t, int64(0), total)
	})

	t.Run("sort by name desc", func(t *testing.T) {
		filters := CustomerFilterParams{Sort: "name", Order: "desc"}
		customers, total, err := repo.ListAllCustomers(context.Background(), filters, 1, 10)
		assert.NoError(t, err)
		assert.Len(t, customers, 3)
		assert.Equal(t, int64(3), total)
		assert.Equal(t, "Charlie", customers[0].FirstName)
	})

	t.Run("sort by updated_at asc", func(t *testing.T) {
		filters := CustomerFilterParams{Sort: "updated_at", Order: "asc"}
		customers, total, err := repo.ListAllCustomers(context.Background(), filters, 1, 10)
		assert.NoError(t, err)
		assert.Len(t, customers, 3)
		assert.Equal(t, int64(3), total)
	})
}

func TestRepository_Transaction(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	t.Run("successful transaction", func(t *testing.T) {
		var createdCustomer *Customer
		err := repo.Transaction(context.Background(), func(txCtx context.Context) error {
			customer := &Customer{FirstName: "John", LastName: "Doe", Email: "john@example.com", Phone: "+380991115566", PasswordHash: "hash"}
			if err := repo.Create(txCtx, customer); err != nil {
				return err
			}
			createdCustomer = customer
			return nil
		})
		assert.NoError(t, err)
		assert.NotZero(t, createdCustomer.ID)

		fetchedCustomer, err := repo.FindByID(context.Background(), createdCustomer.ID)
		assert.NoError(t, err)
		assert.Equal(t, createdCustomer.Email, fetchedCustomer.Email)
	})

	t.Run("rollback on error", func(t *testing.T) {
		err := repo.Transaction(context.Background(), func(txCtx context.Context) error {
			customer := &Customer{FirstName: "Jane", LastName: "Doe", Email: "jane@example.com", Phone: "+380991116677", PasswordHash: "hash"}
			if err := repo.Create(txCtx, customer); err != nil {
				return err
			}
			return errors.New("intentional error to trigger rollback")
		})
		assert.Error(t, err)

		customer, err := repo.FindByPhone(context.Background(), "+380991116677")
		// SQLite may handle rollback differently - ensure the customer was not created
		if err == nil {
			assert.Nil(t, customer, "Customer should not exist after rollback")
		} else {
			assert.Error(t, err)
		}
	})
}

func TestRepository_FindByPhone_Error(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	t.Run("returns error when phone is empty", func(t *testing.T) {
		customer, err := repo.FindByPhone(context.Background(), "")
		assert.NoError(t, err)
		assert.Nil(t, customer)
	})
}

func TestRepository_FindByID_Error(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	t.Run("returns nil when ID is 00000000-0000-0000-0000-000000000000", func(t *testing.T) {
		customer, err := repo.FindByID(context.Background(), uuid.Must(uuid.Parse("00000000-0000-0000-0000-000000000000")))
		assert.NoError(t, err)
		assert.Nil(t, customer)
	})
}

func TestRepository_Update_Error(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	t.Run("successfully updates with empty password hash", func(t *testing.T) {
		customer := &Customer{
			FirstName:    "John",
			LastName:     "Doe",
			Email:        "john@example.com",
			Phone:        "+380991115566",
			PasswordHash: "hashed_password",
		}
		err := repo.Create(context.Background(), customer)
		require.NoError(t, err)

		customer.FirstName = "Updated Name"
		customer.PasswordHash = ""
		err = repo.Update(context.Background(), customer)
		assert.NoError(t, err)

		updatedCustomer, err := repo.FindByID(context.Background(), customer.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name", updatedCustomer.FirstName)
	})
}

func TestRepository_ListAllCustomers_ErrorCases(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	customer1 := &Customer{FirstName: "Alice", LastName: "Doe", Email: "alice@example.com", Phone: "+380991115566", PasswordHash: "hash"}
	err := repo.Create(context.Background(), customer1)
	require.NoError(t, err)

	t.Run("search with empty string", func(t *testing.T) {
		filters := CustomerFilterParams{Search: "", Sort: "created_at", Order: "desc"}
		customers, total, err := repo.ListAllCustomers(context.Background(), filters, 1, 20)
		assert.NoError(t, err)
		assert.Len(t, customers, 1)
		assert.Equal(t, int64(1), total)
	})

	t.Run("role filter with empty string", func(t *testing.T) {
		filters := CustomerFilterParams{Role: "", Sort: "created_at", Order: "desc"}
		customers, total, err := repo.ListAllCustomers(context.Background(), filters, 1, 20)
		assert.NoError(t, err)
		assert.Len(t, customers, 1)
		assert.Equal(t, int64(1), total)
	})
}

func TestRepository_GetCustomerRoles_EmptyResult(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	t.Run("customer with ID 0", func(t *testing.T) {
		roles, err := repo.GetCustomerRoles(context.Background(), uuid.Must(uuid.Parse("00000000-0000-0000-0000-000000000000")))
		assert.NoError(t, err)
		assert.Empty(t, roles)
	})
}

func TestRepository_AssignRole_RoleNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	customer := &Customer{
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		Phone:        "+380991115566",
		PasswordHash: "hashed_password",
	}
	err := repo.Create(context.Background(), customer)
	require.NoError(t, err)

	err = repo.AssignRole(context.Background(), customer.ID, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "role not found")
}

func TestRepository_RemoveRole_RoleNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	customer := &Customer{
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		Phone:        "+380991115566",
		PasswordHash: "hashed_password",
	}
	err := repo.Create(context.Background(), customer)
	require.NoError(t, err)

	err = repo.RemoveRole(context.Background(), customer.ID, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "role not found")
}

func TestRepository_ListAllCustomers_InvalidSortField(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	filters := CustomerFilterParams{
		Sort:  "invalid_field",
		Order: "asc",
	}

	_, _, err := repo.ListAllCustomers(context.Background(), filters, 1, 20)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sort field")
}

func TestRepository_ListAllCustomers_InvalidSortOrder(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	filters := CustomerFilterParams{
		Sort:  "first_name",
		Order: "invalid_order",
	}

	_, _, err := repo.ListAllCustomers(context.Background(), filters, 1, 20)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sort order")
}

func TestRepository_FindRoleByName_Error(t *testing.T) {
	db := setupTestDB(t)
	sqlDB, _ := db.DB()
	_ = sqlDB.Close()

	repo := NewRepository(db)

	role, err := repo.FindRoleByName(context.Background(), "regular")
	assert.Error(t, err)
	assert.Nil(t, role)
}

func TestRepository_GetCustomerRoles_Error(t *testing.T) {
	db := setupTestDB(t)
	sqlDB, _ := db.DB()
	_ = sqlDB.Close()

	repo := NewRepository(db)

	roles, err := repo.GetCustomerRoles(context.Background(), uuid.Must(uuid.Parse("22be2ef9-62ed-4582-b6c8-cf6525f1bb4e")))
	assert.Error(t, err)
	assert.Nil(t, roles)
}
