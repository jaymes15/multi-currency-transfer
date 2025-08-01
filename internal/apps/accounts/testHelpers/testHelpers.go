package accounts

import (
	"context"
	db "lemfi/simplebank/db/sqlc"
	requests "lemfi/simplebank/internal/apps/accounts/requests"
)

// MockAccountRepository implements AccountRespositoryInterface for testing
type MockAccountRepository struct {
	store db.Store
}

func (m *MockAccountRepository) CreateAccount(payload requests.CreateAccountRequest) (db.Account, error) {
	return m.store.CreateAccount(context.Background(), db.CreateAccountParams{
		Owner:    payload.Owner,
		Currency: payload.Currency,
	})
}

func (m *MockAccountRepository) GetAccounts() ([]db.Account, error) {
	return m.store.ListAllAccounts(context.Background(), db.ListAllAccountsParams{
		Limit:  10,
		Offset: 0,
	})
}

// NewMockAccountRepository creates a new mock repository that wraps a store
func NewMockAccountRepository(store db.Store) *MockAccountRepository {
	return &MockAccountRepository{store: store}
}
