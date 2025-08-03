package exchangeRates

import (
	"context"
	db "lemfi/simplebank/db/sqlc"
	requests "lemfi/simplebank/internal/apps/exchangeRates/requests"
)

// MockExchangeRateRepository implements ExchangeRateRepositoryInterface for testing
type MockExchangeRateRepository struct {
	store db.Store
}

func (m *MockExchangeRateRepository) ListExchangeRates(ctx context.Context) ([]db.ExchangeRate, error) {
	rates, err := m.store.ListExchangeRates(ctx)
	if err != nil {
		return []db.ExchangeRate{}, err
	}

	return rates, nil
}

func (m *MockExchangeRateRepository) GetExchangeRate(ctx context.Context, payload requests.GetExchangeRateRequest) (db.ExchangeRate, error) {
	rate, err := m.store.GetExchangeRate(ctx, db.GetExchangeRateParams{
		FromCurrency: payload.FromCurrency,
		ToCurrency:   payload.ToCurrency,
	})
	if err != nil {
		return db.ExchangeRate{}, err
	}

	return rate, nil
}

// NewMockExchangeRateRepository creates a new mock repository that wraps a store
func NewMockExchangeRateRepository(store db.Store) *MockExchangeRateRepository {
	return &MockExchangeRateRepository{store: store}
}
