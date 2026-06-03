package customer

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dch-zeon/yago/internal/auth"
	"github.com/dch-zeon/yago/internal/contextutil"
	apiErrors "github.com/dch-zeon/yago/internal/errors"
)

// Handler handles customer-related HTTP requests
type Handler struct {
	customerService Service
	authService     auth.Service
}

// NewHandler creates a new customer handler
func NewHandler(customerService Service, authService auth.Service) *Handler {
	return &Handler{
		customerService: customerService,
		authService:     authService,
	}
}

// Register godoc
// @Summary Register a new customer
// @Description Register a new customer with phone, name and password, returns access and refresh tokens
// @Tags auth
// @Accept JSON
// @Produce JSON
// @Param request body RegisterRequest true "Registration request"
// @Success 200 {object} errors.Response{success=bool,data=AuthResponse} "Success response with customer data and tokens"
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Validation error"
// @Failure 409 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Phone already exists"
// @Failure 500 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Failed to register customer or generate token"
// @Router /api/v1/auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	customer, err := h.customerService.RegisterCustomer(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrPhoneExists) {
			_ = c.Error(apiErrors.Conflict("Phone already exists"))
			return
		}
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	tokenPair, err := h.authService.GenerateTokenPair(c.Request.Context(), customer.ID, customer.Phone, customer.FirstName+" "+customer.LastName)
	if err != nil {
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, apiErrors.Success(AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		Customer:     ToCustomerResponse(customer),
	}))
}

// Login godoc
// @Summary Login customer
// @Description Authenticate customer with phone and password, returns access and refresh tokens
// @Tags auth
// @Accept JSON
// @Produce JSON
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} errors.Response{success=bool,data=AuthResponse} "Success response with customer data and tokens"
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Validation error"
// @Failure 401 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Invalid phone or password"
// @Failure 500 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Failed to authenticate customer or generate token"
// @Router /api/v1/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	customer, err := h.customerService.AuthenticateCustomer(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			_ = c.Error(apiErrors.Unauthorized("Invalid phone or password"))
			return
		}
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	tokenPair, err := h.authService.GenerateTokenPair(c.Request.Context(), customer.ID, customer.Phone, customer.FirstName+" "+customer.LastName)
	if err != nil {
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, apiErrors.Success(AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		Customer:     ToCustomerResponse(customer),
	}))
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Exchange refresh token for new access and refresh tokens with automatic rotation
// @Tags auth
// @Accept JSON
// @Produce JSON
// @Param request body auth.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} errors.Response{success=bool,data=auth.TokenPairResponse} "Success response with new token pair"
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Validation error"
// @Failure 401 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Invalid or expired refresh token"
// @Failure 403 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Token reuse detected - all tokens revoked"
// @Failure 500 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Failed to refresh token"
// @Router /api/v1/auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var req auth.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	tokenPair, err := h.authService.RefreshAccessToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, auth.ErrExpiredToken) {
			_ = c.Error(apiErrors.Unauthorized("Invalid or expired refresh token"))
			return
		}
		if errors.Is(err, auth.ErrTokenReuse) {
			_ = c.Error(apiErrors.Forbidden("Token reuse detected. All tokens have been revoked for security."))
			return
		}
		if errors.Is(err, auth.ErrTokenRevoked) {
			_ = c.Error(apiErrors.Unauthorized("Token has been revoked"))
			return
		}
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, apiErrors.Success(auth.TokenPairResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	}))
}

// Logout godoc
// @Summary Logout customer
// @Description Revoke refresh token and invalidate customer session
// @Tags auth
// @Accept JSON
// @Produce JSON
// @Security BearerAuth
// @Param request body auth.RefreshTokenRequest true "Refresh token to revoke"
// @Success 200 {object} errors.Response{success=bool,data=object} "Successfully logged out"
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Validation error"
// @Failure 401 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Unauthorized"
// @Failure 403 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Token does not belong to customer"
// @Failure 500 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Failed to logout"
// @Router /api/v1/auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	customerID := contextutil.GetCustomerID(c)
	if customerID == uuid.Nil {
		_ = c.Error(apiErrors.Unauthorized("customer not authenticated"))
		return
	}

	var req auth.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	if err := h.authService.RevokeCustomerRefreshToken(c.Request.Context(), customerID, req.RefreshToken); err != nil {
		if errors.Is(err, auth.ErrTokenDoesNotBelongToCustomer) {
			_ = c.Error(apiErrors.Forbidden("token does not belong to customer"))
			return
		}
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, apiErrors.Success(gin.H{"message": "Successfully logged out"}))
}
