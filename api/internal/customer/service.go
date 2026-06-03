package customer

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrPhoneExists is returned when phone already exists
	ErrPhoneExists = errors.New("phone already exists")
	// ErrInvalidCredentials is returned when credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Service defines customer service interface
type Service interface {
	RegisterCustomer(ctx context.Context, req RegisterRequest) (*Customer, error)
	AuthenticateCustomer(ctx context.Context, req LoginRequest) (*Customer, error)
}

type service struct {
	repo Repository
}

// NewService creates a new customer service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// RegisterCustomer registers a new customer
func (s *service) RegisterCustomer(ctx context.Context, req RegisterRequest) (*Customer, error) {
	existingCustomer, err := s.repo.FindByPhone(ctx, req.Phone)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing phone: %w", err)
	}
	if existingCustomer != nil {
		return nil, ErrPhoneExists
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	customer := &Customer{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
	}

	// Use transaction to ensure atomic customer creation and role assignment
	err = s.repo.Transaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.Create(txCtx, customer); err != nil {
			return fmt.Errorf("failed to create customer: %w", err)
		}

		if err := s.repo.AssignRole(txCtx, customer.ID, RoleRegular); err != nil {
			return fmt.Errorf("failed to assign default role: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload customer with roles after successful transaction
	customer, err = s.repo.FindByID(ctx, customer.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload customer: %w", err)
	}
	if customer == nil {
		return nil, fmt.Errorf("failed to reload customer: customer not found after creation")
	}

	return customer, nil
}

// AuthenticateCustomer authenticates a customer with phone and password
func (s *service) AuthenticateCustomer(ctx context.Context, req LoginRequest) (*Customer, error) {
	customer, err := s.repo.FindByPhone(ctx, req.Phone)
	if err != nil {
		return nil, fmt.Errorf("failed to find customer: %w", err)
	}
	if customer == nil {
		return nil, ErrInvalidCredentials
	}

	if err := verifyPassword(customer.PasswordHash, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	return customer, nil
}

// hashPassword hashes a plain text password using bcrypt
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// verifyPassword verifies a password against a hash
func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
