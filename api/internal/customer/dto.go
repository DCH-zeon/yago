package customer

import "github.com/google/uuid"

// RegisterRequest represents registration request payload
type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"omitempty,min=2,max=36"`
	LastName  string `json:"last_name" binding:"omitempty,min=2,max=36"`
	Email     string `json:"email" binding:"omitempty,email"`
	Phone     string `json:"phone" binding:"required,e164"`
	Password  string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents login request payload
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required,e164"`
	Password string `json:"password" binding:"required"`
}

// CustomerResponse represents customer response (without sensitive fields)
type CustomerResponse struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email,omitempty"`
	Phone     string    `json:"phone"`
	Roles     []string  `json:"roles"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	TokenType    string           `json:"token_type"`
	ExpiresIn    int64            `json:"expires_in"`
	Customer     CustomerResponse `json:"customer"`
}

// ToCustomerResponse converts Customer model to CustomerResponse DTO
func ToCustomerResponse(customer *Customer) CustomerResponse {
	return CustomerResponse{
		ID:        customer.ID,
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		Email:     customer.Email,
		Phone:     customer.Phone,
		Roles:     customer.GetRoleNames(),
		CreatedAt: customer.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: customer.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
