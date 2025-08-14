package users

import (
	"context"
	db "lemfi/simplebank/db/sqlc"
	requests "lemfi/simplebank/internal/apps/users/requests"
)

// MockUserRepository implements UserRespositoryInterface for testing
type MockUserRepository struct {
	store db.Store
}

func (m *MockUserRepository) CreateUser(payload requests.CreateUserRequest) (db.CreateUserRow, error) {
	return m.store.CreateUser(context.Background(), db.CreateUserParams{
		Username:       payload.Username,
		HashedPassword: payload.HashedPassword,
		FullName:       payload.FullName,
		Email:          payload.Email,
	})
}

func (m *MockUserRepository) GetUserHashedPassword(username string) (string, error) {
	// Mock implementation - return a hashed password for testing
	return "$2a$10$hashedpassword123", nil
}

// NewMockUserRepository creates a new mock repository that wraps a store
func NewMockUserRepository(store db.Store) *MockUserRepository {
	return &MockUserRepository{store: store}
}
