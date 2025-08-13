package users

import (
	"errors"
	"testing"
	"time"

	db "lemfi/simplebank/db/sqlc"
	requests "lemfi/simplebank/internal/apps/users/requests"

	"github.com/stretchr/testify/require"
)

// MockUserRepository for testing
type MockUserRepository struct {
	createUserFunc func(payload requests.CreateUserRequest) (db.CreateUserRow, error)
}

func (m *MockUserRepository) CreateUser(payload requests.CreateUserRequest) (db.CreateUserRow, error) {
	return m.createUserFunc(payload)
}

func TestCreateUser_Success(t *testing.T) {
	// Create mock repository
	mockRepo := &MockUserRepository{
		createUserFunc: func(payload requests.CreateUserRequest) (db.CreateUserRow, error) {
			// Return a mock user response
			return db.CreateUserRow{
				Username:  payload.Username,
				FullName:  payload.FullName,
				Email:     payload.Email,
				CreatedAt: time.Now(),
			}, nil
		},
	}

	userService := &UserService{
		userRespository: mockRepo,
	}

	// Test request
	request := requests.CreateUserRequest{
		Username: "testuser",
		Password: "password123",
		FullName: "Test User",
		Email:    "test@example.com",
	}

	// Call service
	response, err := userService.CreateUser(request)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, request.Username, response.Username)
	require.Equal(t, request.FullName, response.FullName)
	require.Equal(t, request.Email, response.Email)
}

func TestCreateUser_RepositoryError(t *testing.T) {
	// Create mock repository that returns an error
	mockRepo := &MockUserRepository{
		createUserFunc: func(payload requests.CreateUserRequest) (db.CreateUserRow, error) {
			return db.CreateUserRow{}, errors.New("database error")
		},
	}

	userService := &UserService{
		userRespository: mockRepo,
	}

	// Test request
	request := requests.CreateUserRequest{
		Username: "testuser",
		Password: "password123",
		FullName: "Test User",
		Email:    "test@example.com",
	}

	// Call service
	response, err := userService.CreateUser(request)

	// Assertions
	require.Error(t, err)
	require.Equal(t, "database error", err.Error())
	require.Empty(t, response)
}
