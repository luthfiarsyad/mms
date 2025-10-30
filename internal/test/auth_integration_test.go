package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/luthfiarsyad/mms/internal/domain/user"
	"github.com/luthfiarsyad/mms/internal/infrastructure/persistence/mysql"
	httpInterface "github.com/luthfiarsyad/mms/internal/interface/http"
	"github.com/luthfiarsyad/mms/internal/interface/http/request"
)

func TestAuthIntegration(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Setup test database
	db := SetupTestDatabase(t)
	defer CleanupTestDatabase(t, db)

	// Set the global DB instance for the handlers to use
	mysql.DB = db

	// Setup repositories and services
	userRepo := mysql.NewUserRepo(db)
	userService := user.NewService(userRepo)

	t.Run("Health endpoint works", func(t *testing.T) {
		// Arrange
		router := gin.New()
		httpInterface.SetupRoutes(router)

		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "ok", response["status"])
	})

	t.Run("User registration through handler", func(t *testing.T) {
		// Arrange
		router := gin.New()
		httpInterface.SetupRoutes(router)

		registerReq := request.RegisterRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}

		jsonData, _ := json.Marshal(registerReq)
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "id")
		assert.Contains(t, response, "email")
		assert.Contains(t, response, "created_at")
		assert.Equal(t, "test@example.com", response["email"])
	})

	t.Run("User login through handler", func(t *testing.T) {
		// First, create a user directly in the database
		ctx := context.Background()
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		require.NoError(t, err)

		testUser := &user.User{
			Name:     "Login Test User",
			Email:    "login@example.com",
			Password: string(hashedPassword),
		}

		err = userService.Register(ctx, testUser)
		require.NoError(t, err)

		// Now test login through handler
		router := gin.New()
		httpInterface.SetupRoutes(router)

		loginReq := request.LoginRequest{
			Email:    "login@example.com",
			Password: "password123",
		}

		jsonData, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "access_token")
		assert.Contains(t, response, "token_type")
		assert.Contains(t, response, "expires_in")
		assert.Equal(t, "bearer", response["token_type"])
		assert.Equal(t, float64(86400), response["expires_in"]) // 24 hours in seconds
	})

	t.Run("Login with wrong credentials fails", func(t *testing.T) {
		// Arrange
		router := gin.New()
		httpInterface.SetupRoutes(router)

		loginReq := request.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "wrongpassword",
		}

		jsonData, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "invalid credentials", response["error"])
	})

	t.Run("Registration with duplicate email fails", func(t *testing.T) {
		// First, create a user
		ctx := context.Background()
		testUser := &user.User{
			Name:     "Duplicate Test User",
			Email:    "duplicate@example.com",
			Password: "hashedpassword",
		}

		err := userService.Register(ctx, testUser)
		require.NoError(t, err)

		// Now try to register the same email through handler
		router := gin.New()
		httpInterface.SetupRoutes(router)

		registerReq := request.RegisterRequest{
			Name:     "Another User",
			Email:    "duplicate@example.com", // Same email
			Password: "password123",
		}

		jsonData, _ := json.Marshal(registerReq)
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "error")
	})

	t.Run("Non-implemented endpoints return 501", func(t *testing.T) {
		// Arrange
		router := gin.New()
		httpInterface.SetupRoutes(router)

		testCases := []struct {
			method string
			path   string
		}{
			{"GET", "/api/v1/users"},
			{"POST", "/api/v1/users"},
			{"GET", "/api/v1/users/1"},
			{"PUT", "/api/v1/users/1"},
			{"DELETE", "/api/v1/users/1"},
			{"GET", "/api/v1/transactions"},
			{"POST", "/api/v1/transactions"},
			{"GET", "/api/v1/transactions/1"},
		}

		for _, tc := range testCases {
			t.Run(tc.method+" "+tc.path, func(t *testing.T) {
				req, _ := http.NewRequest(tc.method, tc.path, nil)
				w := httptest.NewRecorder()

				// Act
				router.ServeHTTP(w, req)

				// Assert
				assert.Equal(t, http.StatusNotImplemented, w.Code)

				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, "not implemented", response["error"])
			})
		}
	})
}

func TestUserService_Unit(t *testing.T) {
	// Setup test database
	db := SetupTestDatabase(t)
	defer CleanupTestDatabase(t, db)

	// Setup repositories and services
	userRepo := mysql.NewUserRepo(db)
	userService := user.NewService(userRepo)

	ctx := context.Background()

	t.Run("User Service Register works", func(t *testing.T) {
		// Arrange
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
		require.NoError(t, err)

		testUser := &user.User{
			Name:     "Service Test User",
			Email:    "service@example.com",
			Password: string(hashedPassword),
		}

		// Act
		err = userService.Register(ctx, testUser)

		// Assert
		assert.NoError(t, err)
		assert.NotZero(t, testUser.ID)
		assert.Equal(t, "Service Test User", testUser.Name)
		assert.Equal(t, "service@example.com", testUser.Email)
	})

	t.Run("User Service Authenticate works", func(t *testing.T) {
		// First create a user
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("authpassword"), bcrypt.DefaultCost)
		require.NoError(t, err)

		testUser := &user.User{
			Name:     "Auth Test User",
			Email:    "auth@example.com",
			Password: string(hashedPassword),
		}

		err = userService.Register(ctx, testUser)
		require.NoError(t, err)

		// Test authentication
		foundUser, err := userService.Authenticate(ctx, "auth@example.com", "authpassword")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, testUser.ID, foundUser.ID)
		assert.Equal(t, "Auth Test User", foundUser.Name)
	})

	t.Run("User Service Authenticate fails with wrong credentials", func(t *testing.T) {
		// Test authentication with wrong email
		foundUser, err := userService.Authenticate(ctx, "wrong@example.com", "anypassword")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, user.ErrInvalidCreds, err)
		assert.Nil(t, foundUser)
	})
}
