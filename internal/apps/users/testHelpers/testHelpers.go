package users

import (
	"context"
	db "lemfi/simplebank/db/sqlc"
	requests "lemfi/simplebank/internal/apps/users/requests"
	"time"

	"github.com/google/uuid"
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

func (m *MockUserRepository) GetUser(username string) (db.GetUserRow, error) {
	return m.store.GetUser(context.Background(), username)
}

func (m *MockUserRepository) CreateSession(username string, refreshTokenID uuid.UUID, refreshToken string, expiresAt time.Time) error {
	_, err := m.store.CreateSession(context.Background(), db.CreateSessionParams{
		ID:           refreshTokenID,
		Username:     username,
		RefreshToken: refreshToken,
		UserAgent:    "", // not relevant for tests
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    expiresAt,
	})
	return err
}

func (m *MockUserRepository) GetSession(refreshTokenID uuid.UUID) (db.GetSessionRow, error) {
	return m.store.GetSession(context.Background(), refreshTokenID)
}

func (m *MockUserRepository) BlockSession(sessionID uuid.UUID) error {
	// For tests we don't need to actually hit DB here
	return nil
}

// NewMockUserRepository creates a new mock repository that wraps a store
func NewMockUserRepository(store db.Store) *MockUserRepository {
	return &MockUserRepository{store: store}
}
