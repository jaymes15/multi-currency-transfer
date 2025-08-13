package users

import (
	"context"
	"errors"
	"testing"
	"time"

	db "lemfi/simplebank/db/sqlc"
	userErrors "lemfi/simplebank/internal/apps/users/errors"
	requests "lemfi/simplebank/internal/apps/users/requests"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

// MockStore implements the db.Store interface for testing
type MockStore struct {
	createUserFunc func(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error)
}

func (m *MockStore) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error) {
	return m.createUserFunc(ctx, arg)
}

// Implement other required methods with empty implementations for testing
func (m *MockStore) CreateAccount(ctx context.Context, arg db.CreateAccountParams) (db.Account, error) {
	return db.Account{}, nil
}

func (m *MockStore) GetAccount(ctx context.Context, id int64) (db.Account, error) {
	return db.Account{}, nil
}

func (m *MockStore) GetAccountForUpdate(ctx context.Context, id int64) (db.Account, error) {
	return db.Account{}, nil
}

func (m *MockStore) ListAccounts(ctx context.Context, arg db.ListAccountsParams) ([]db.Account, error) {
	return nil, nil
}

func (m *MockStore) ListAllAccounts(ctx context.Context, arg db.ListAllAccountsParams) ([]db.Account, error) {
	return nil, nil
}

func (m *MockStore) UpdateAccount(ctx context.Context, arg db.UpdateAccountParams) (db.Account, error) {
	return db.Account{}, nil
}

func (m *MockStore) AddAccountBalance(ctx context.Context, arg db.AddAccountBalanceParams) (decimal.Decimal, error) {
	return decimal.Zero, nil
}

func (m *MockStore) DeleteAccount(ctx context.Context, id int64) error {
	return nil
}

func (m *MockStore) CreateEntry(ctx context.Context, arg db.CreateEntryParams) (db.Entry, error) {
	return db.Entry{}, nil
}

func (m *MockStore) GetEntry(ctx context.Context, id int64) (db.Entry, error) {
	return db.Entry{}, nil
}

func (m *MockStore) ListEntries(ctx context.Context, arg db.ListEntriesParams) ([]db.Entry, error) {
	return nil, nil
}

func (m *MockStore) CreateTransfer(ctx context.Context, arg db.CreateTransferParams) (db.CreateTransferRow, error) {
	return db.CreateTransferRow{}, nil
}

func (m *MockStore) GetTransfer(ctx context.Context, id int64) (db.GetTransferRow, error) {
	return db.GetTransferRow{}, nil
}

func (m *MockStore) ListTransfers(ctx context.Context, arg db.ListTransfersParams) ([]db.ListTransfersRow, error) {
	return nil, nil
}

func (m *MockStore) CreateExchangeRate(ctx context.Context, arg db.CreateExchangeRateParams) (db.ExchangeRate, error) {
	return db.ExchangeRate{}, nil
}

func (m *MockStore) GetExchangeRate(ctx context.Context, arg db.GetExchangeRateParams) (db.ExchangeRate, error) {
	return db.ExchangeRate{}, nil
}

func (m *MockStore) ListExchangeRates(ctx context.Context) ([]db.ExchangeRate, error) {
	return nil, nil
}

func (m *MockStore) UpdateExchangeRate(ctx context.Context, arg db.UpdateExchangeRateParams) (db.ExchangeRate, error) {
	return db.ExchangeRate{}, nil
}

func (m *MockStore) GetUser(ctx context.Context, username string) (db.GetUserRow, error) {
	return db.GetUserRow{}, nil
}

func (m *MockStore) ListUsers(ctx context.Context) ([]db.ListUsersRow, error) {
	return nil, nil
}

func (m *MockStore) UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.UpdateUserRow, error) {
	return db.UpdateUserRow{}, nil
}

func (m *MockStore) DeleteUser(ctx context.Context, username string) error {
	return nil
}

func (m *MockStore) TransferTx(ctx context.Context, arg db.TransferTxParams) (db.TransferTxResult, error) {
	return db.TransferTxResult{}, nil
}

func TestCreateUser_Success(t *testing.T) {
	// Create mock store
	mockStore := &MockStore{
		createUserFunc: func(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error) {
			return db.CreateUserRow{
				Username:  arg.Username,
				FullName:  arg.FullName,
				Email:     arg.Email,
				CreatedAt: time.Now(),
			}, nil
		},
	}

	// Create repository with mock store
	repo := &UserRespository{
		context: context.Background(),
		queries: mockStore,
	}

	// Test request
	request := requests.CreateUserRequest{
		Username:       "testuser",
		HashedPassword: "hashedpassword123",
		FullName:       "Test User",
		Email:          "test@example.com",
	}

	// Call repository
	user, err := repo.CreateUser(request)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, request.Username, user.Username)
	require.Equal(t, request.FullName, user.FullName)
	require.Equal(t, request.Email, user.Email)
}

func TestCreateUser_DuplicateUsername(t *testing.T) {
	// Create mock store that returns duplicate key error
	mockStore := &MockStore{
		createUserFunc: func(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error) {
			return db.CreateUserRow{}, errors.New("duplicate key value violates unique constraint \"users_username_key\"")
		},
	}

	// Create repository with mock store
	repo := &UserRespository{
		context: context.Background(),
		queries: mockStore,
	}

	// Test request
	request := requests.CreateUserRequest{
		Username:       "testuser",
		HashedPassword: "hashedpassword123",
		FullName:       "Test User",
		Email:          "test@example.com",
	}

	// Call repository
	user, err := repo.CreateUser(request)

	// Assertions
	require.Error(t, err)
	require.Equal(t, userErrors.ErrDuplicateUsername, err)
	require.Empty(t, user)
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	// Create mock store that returns duplicate key error
	mockStore := &MockStore{
		createUserFunc: func(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error) {
			return db.CreateUserRow{}, errors.New("duplicate key value violates unique constraint \"users_email_key\"")
		},
	}

	// Create repository with mock store
	repo := &UserRespository{
		context: context.Background(),
		queries: mockStore,
	}

	// Test request
	request := requests.CreateUserRequest{
		Username:       "testuser",
		HashedPassword: "hashedpassword123",
		FullName:       "Test User",
		Email:          "test@example.com",
	}

	// Call repository
	user, err := repo.CreateUser(request)

	// Assertions
	require.Error(t, err)
	require.Equal(t, userErrors.ErrDuplicateEmail, err)
	require.Empty(t, user)
}

func TestCreateUser_GeneralError(t *testing.T) {
	// Create mock store that returns general error
	mockStore := &MockStore{
		createUserFunc: func(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error) {
			return db.CreateUserRow{}, errors.New("general database error")
		},
	}

	// Create repository with mock store
	repo := &UserRespository{
		context: context.Background(),
		queries: mockStore,
	}

	// Test request
	request := requests.CreateUserRequest{
		Username:       "testuser",
		HashedPassword: "hashedpassword123",
		FullName:       "Test User",
		Email:          "test@example.com",
	}

	// Call repository
	user, err := repo.CreateUser(request)

	// Assertions
	require.Error(t, err)
	require.Equal(t, "general database error", err.Error())
	require.Empty(t, user)
}
