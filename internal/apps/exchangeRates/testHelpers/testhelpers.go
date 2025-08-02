package exchangeRates

import (
	"context"
	db "lemfi/simplebank/db/sqlc"
	requests "lemfi/simplebank/internal/apps/exchangeRates/requests"
	responses "lemfi/simplebank/internal/apps/exchangeRates/responses"
)

// MockExchangeRateRepository implements ExchangeRateRepositoryInterface for testing
type MockExchangeRateRepository struct {
	store db.Store
}

func (m *MockExchangeRateRepository) ListExchangeRates(ctx context.Context) (responses.ListExchangeRatesResponse, error) {
	rates, err := m.store.ListExchangeRates(ctx)
	if err != nil {
		return responses.ListExchangeRatesResponse{}, err
	}

	// Convert database rates to response format
	var exchangeRates []responses.ExchangeRateResponse
	for _, rate := range rates {
		exchangeRates = append(exchangeRates, responses.ExchangeRateResponse{
			ID:           rate.ID,
			FromCurrency: rate.FromCurrency,
			ToCurrency:   rate.ToCurrency,
			Rate:         rate.Rate,
			CreatedAt:    rate.CreatedAt.Time,
		})
	}

	return responses.ListExchangeRatesResponse{
		ExchangeRates: exchangeRates,
		Total:         len(exchangeRates),
	}, nil
}

func (m *MockExchangeRateRepository) GetExchangeRate(ctx context.Context, payload requests.GetExchangeRateRequest) (responses.GetExchangeRateResponse, error) {
	rate, err := m.store.GetExchangeRate(ctx, db.GetExchangeRateParams{
		FromCurrency: payload.FromCurrency,
		ToCurrency:   payload.ToCurrency,
	})
	if err != nil {
		return responses.GetExchangeRateResponse{}, err
	}

	// Calculate amounts
	amountToSend := payload.Amount
	amountToReceive := payload.Amount.Mul(rate.Rate).Round(2)

	return responses.GetExchangeRateResponse{
		ExchangeRate: responses.ExchangeRateResponse{
			ID:           rate.ID,
			FromCurrency: rate.FromCurrency,
			ToCurrency:   rate.ToCurrency,
			Rate:         rate.Rate,
			CreatedAt:    rate.CreatedAt.Time,
		},
		AmountToSend:    amountToSend,
		AmountToReceive: amountToReceive,
		CanTransact:     true,
		Message:         "Exchange rate available for transaction",
	}, nil
}

// NewMockExchangeRateRepository creates a new mock repository that wraps a store
func NewMockExchangeRateRepository(store db.Store) *MockExchangeRateRepository {
	return &MockExchangeRateRepository{store: store}
}
