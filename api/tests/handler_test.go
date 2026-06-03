package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/dch-zeon/yago/internal/auth"
	"github.com/dch-zeon/yago/internal/config"
	"github.com/dch-zeon/yago/internal/customer"
	"github.com/dch-zeon/yago/internal/db"
	"github.com/dch-zeon/yago/internal/server"
)

// createTestSchema creates the SQLite test schema using GORM AutoMigrate for consistency
func createTestSchema(t *testing.T, database *gorm.DB) {
	t.Helper()

	err := database.AutoMigrate(&customer.Customer{}, &customer.Role{}, &auth.RefreshToken{})
	assert.NoError(t, err)

	// Drop the auto-created customer_roles table (created by GORM for many2many)
	// and recreate it with our custom schema including assigned_at column
	database.Exec("DROP TABLE IF EXISTS customer_roles")

	// Manually create the customer_roles junction table with the assigned_at column
	err = database.Exec(`
		CREATE TABLE customer_roles (
            customer_id UUID NOT NULL,
			role_id INTEGER NOT NULL,
			assigned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (customer_id, role_id),
			FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
			FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
		)
	`).Error
	assert.NoError(t, err)

	// Seed role data - use FirstOrCreate to avoid duplicate errors
	roles := []customer.Role{
		{ID: 1, Name: "regular", Description: "Звичайний клієнт зі стандартними дозволами"},
		{ID: 2, Name: "wholesale", Description: "Оптовий клієнт зі спеціальними цінами та доступом до оптових замовлень"},
	}
	for _, role := range roles {
		var existingRole customer.Role
		result := database.Where("name = ?", role.Name).FirstOrCreate(&existingRole, &role)
		if result.Error != nil {
			t.Fatalf("Failed to create role %s: %v", role.Name, result.Error)
		}
	}
}

func setupTestRouter(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)

	testCfg := config.NewTestConfig()

	database, err := db.NewSQLiteDB(":memory:")
	assert.NoError(t, err)

	createTestSchema(t, database)

	authService := auth.NewServiceWithRepo(&testCfg.JWT, database)
	customerRepo := customer.NewRepository(database)
	customerService := customer.NewService(customerRepo)
	customerHandler := customer.NewHandler(customerService, authService)

	router := server.SetupRouter(customerHandler, authService, testCfg, database)

	return router
}

func setupRateLimitTestRouter(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)

	testCfg := config.NewTestConfig()
	testCfg.Ratelimit.Enabled = true
	testCfg.Ratelimit.Requests = 10
	testCfg.Ratelimit.Window = time.Minute

	database, err := db.NewSQLiteDB(":memory:")
	assert.NoError(t, err)

	createTestSchema(t, database)

	authService := auth.NewServiceWithRepo(&testCfg.JWT, database)
	customerRepo := customer.NewRepository(database)
	customerService := customer.NewService(customerRepo)
	customerHandler := customer.NewHandler(customerService, authService)

	return server.SetupRouter(customerHandler, authService, testCfg, database)
}

func TestRegisterHandler(t *testing.T) {
	router := setupTestRouter(t)

	tests := []struct {
		name           string
		payload        map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "successful registration",
			payload: map[string]string{
				"firstName": "John",
				"lastName":  "Doe",
				"phone":     "+380997775566",
				"password":  "password123",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if success, ok := body["success"].(bool); !ok || !success {
					t.Error("Expected success to be true in response")
				}
				data, ok := body["data"].(map[string]interface{})
				if !ok {
					t.Fatal("Expected data object in response")
				}
				if accessToken, ok := data["access_token"].(string); !ok || accessToken == "" {
					t.Error("Expected access_token in response data")
				}
				if refreshToken, ok := data["refresh_token"].(string); !ok || refreshToken == "" {
					t.Error("Expected refresh_token in response data")
				}
				if customerData, ok := data["customer"].(map[string]interface{}); !ok {
					t.Error("Expected customer object in response data")
				} else {
					if phone, ok := customerData["phone"].(string); !ok || phone != "+380997775566" {
						t.Errorf("Expected phone '+380997775566', got '%v'", phone)
					}
				}
			},
		},
		{
			name: "duplicate phone",
			payload: map[string]string{
				"firstName": "Jane",
				"lastName":  "Doe",
				"phone":     "+380997775566",
				"password":  "password123",
			},
			expectedStatus: http.StatusConflict,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if success, ok := body["success"].(bool); !ok || success {
					t.Error("Expected success to be false for error response")
				}
				errorInfo, ok := body["error"].(map[string]interface{})
				if !ok {
					t.Fatal("Expected error object in response")
				}
				if errorMsg, ok := errorInfo["message"].(string); !ok || errorMsg == "" {
					t.Error("Expected error message in response")
				}
			},
		},
		{
			name: "invalid phone format",
			payload: map[string]string{
				"firstName": "Invalid",
				"lastName":  "User",
				"phone":     "not-an-phone",
				"password":  "password123",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if success, ok := body["success"].(bool); !ok || success {
					t.Error("Expected success to be false for error response")
				}
				errorInfo, ok := body["error"].(map[string]interface{})
				if !ok {
					t.Fatal("Expected error object in response")
				}
				if errorMsg, ok := errorInfo["message"].(string); !ok || errorMsg == "" {
					t.Error("Expected error message in response")
				}
			},
		},
		{
			name: "missing required fields",
			payload: map[string]string{
				"firstName": "Incomplete Customer",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if success, ok := body["success"].(bool); !ok || success {
					t.Error("Expected success to be false for error response")
				}
				errorInfo, ok := body["error"].(map[string]interface{})
				if !ok {
					t.Fatal("Expected error object in response")
				}
				if errorMsg, ok := errorInfo["message"].(string); !ok || errorMsg == "" {
					t.Error("Expected error message in response")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonPayload, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Logf("Response body: %s", w.Body.String())
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestLoginHandler(t *testing.T) {
	router := setupTestRouter(t)

	registerPayload := map[string]string{
		"firstName": "Test",
		"lastName":  "Customer",
		"phone":     "+380997775566",
		"password":  "testpassword123",
	}
	jsonPayload, _ := json.Marshal(registerPayload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	tests := []struct {
		name           string
		payload        map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "successful login",
			payload: map[string]string{
				"phone":    "+380997775566",
				"password": "testpassword123",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if success, ok := body["success"].(bool); !ok || !success {
					t.Error("Expected success to be true in response")
				}
				data, ok := body["data"].(map[string]interface{})
				if !ok {
					t.Fatal("Expected data object in response")
				}
				if accessToken, ok := data["access_token"].(string); !ok || accessToken == "" {
					t.Error("Expected access_token in response data")
				}
				if refreshToken, ok := data["refresh_token"].(string); !ok || refreshToken == "" {
					t.Error("Expected refresh_token in response data")
				}
			},
		},
		{
			name: "invalid password",
			payload: map[string]string{
				"phone":    "+380997775566",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if success, ok := body["success"].(bool); !ok || success {
					t.Error("Expected success to be false for error response")
				}
				errorInfo, ok := body["error"].(map[string]interface{})
				if !ok {
					t.Fatal("Expected error object in response")
				}
				if errorMsg, ok := errorInfo["message"].(string); !ok || errorMsg == "" {
					t.Error("Expected error message in response")
				}
			},
		},
		{
			name: "non-existent customer",
			payload: map[string]string{
				"phone":    "+380998885566",
				"password": "password123",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if success, ok := body["success"].(bool); !ok || success {
					t.Error("Expected success to be false for error response")
				}
				errorInfo, ok := body["error"].(map[string]interface{})
				if !ok {
					t.Fatal("Expected error object in response")
				}
				if errorMsg, ok := errorInfo["message"].(string); !ok || errorMsg == "" {
					t.Error("Expected error message in response")
				}
			},
		},
		{
			name: "missing credentials",
			payload: map[string]string{
				"phone": "+380997775566",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if success, ok := body["success"].(bool); !ok || success {
					t.Error("Expected success to be false for error response")
				}
				errorInfo, ok := body["error"].(map[string]interface{})
				if !ok {
					t.Fatal("Expected error object in response")
				}
				if errorMsg, ok := errorInfo["message"].(string); !ok || errorMsg == "" {
					t.Error("Expected error message in response")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonPayload, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestHealthEndpoint(t *testing.T) {
	router := setupTestRouter(t)

	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal health response: %v", err)
	}
	if status, ok := response["status"].(string); !ok || status != "healthy" {
		t.Errorf("Expected status 'healthy' in health check response, got %v", status)
	}
}

func TestRateLimit_BlocksThenAllows(t *testing.T) {
	r := setupRateLimitTestRouter(t)

	testIP := fmt.Sprintf("192.168.1.%d", time.Now().UnixNano()%255)

	registerBody, _ := json.Marshal(map[string]string{
		"name":     "Rate Test",
		"phone":    fmt.Sprintf("rate%d@example.com", time.Now().UnixNano()),
		"password": "secret123",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(registerBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", testIP)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Logf("Response body: %s", rr.Body.String())
		t.Fatalf("register expected 200, got %d", rr.Code)
	}

	var registerResp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &registerResp); err != nil {
		t.Fatalf("Failed to unmarshal register response: %v", err)
	}
	dataResp, ok := registerResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected data object in register response")
	}
	customerResp, ok := dataResp["customer"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected customer object in register response data")
	}
	phone := customerResp["phone"].(string)

	loginBody, _ := json.Marshal(map[string]string{
		"phone":    phone,
		"password": "secret123",
	})

	successCount := 0
	for i := 0; i < 15; i++ {
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(loginBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Forwarded-For", testIP)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code == http.StatusOK {
			successCount++
		} else if rr.Code == http.StatusTooManyRequests {
			retryAfterStr := rr.Header().Get("Retry-After")
			if retryAfterStr == "" {
				t.Fatalf("expected Retry-After header on 429")
			}
			retryAfterSec, err := strconv.Atoi(retryAfterStr)
			if err != nil || retryAfterSec <= 0 {
				t.Fatalf("Retry-After should be positive integer seconds, got %q (err=%v)", retryAfterStr, err)
			}
			t.Logf("Rate limit triggered after %d successful requests (including register)", successCount+1)
			return
		} else {
			t.Fatalf("login #%d expected 200 or 429, got %d", i+1, rr.Code)
		}
	}

	// If we get here, rate limiting didn't work
	t.Fatalf("expected rate limiting to trigger, but completed %d requests without 429", successCount)
}
