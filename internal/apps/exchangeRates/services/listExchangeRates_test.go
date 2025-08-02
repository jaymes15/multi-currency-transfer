package exchangeRates

import (
	"context"
	"errors"
	mockdb "lemfi/simplebank/db/mock"
	db "lemfi/simplebank/db/sqlc"
	testhelpers "lemfi/simplebank/internal/apps/exchangeRates/testHelpers"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListExchangeRatesService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// Create test data
	expectedRates := []db.ExchangeRate{
		{
			ID:           1,
			FromCurrency: "USD",
			ToCurrency:   "EUR",
			Rate:         decimal.NewFromFloat(0.85),
		},
		{
			ID:           2,
			FromCurrency: "EUR",
			ToCurrency:   "USD",
			Rate:         decimal.NewFromFloat(1.18),
		},
	}

	// Expect exchange rates to be listed
	store.EXPECT().ListExchangeRates(gomock.Any()).Return(expectedRates, nil).Times(1)

	mockRepo := testhelpers.NewMockExchangeRateRepository(store)
	exchangeRateService := NewExchangeRateService(mockRepo)

	// Test the service
	result, err := exchangeRateService.ListExchangeRates(context.Background())

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.ExchangeRates, 2)
	require.Equal(t, 2, result.Total)
	require.Equal(t, "USD", result.ExchangeRates[0].FromCurrency)
	require.Equal(t, "EUR", result.ExchangeRates[0].ToCurrency)
	require.Equal(t, "EUR", result.ExchangeRates[1].FromCurrency)
	require.Equal(t, "USD", result.ExchangeRates[1].ToCurrency)
}

func TestListExchangeRatesService_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// Expect database error
	store.EXPECT().ListExchangeRates(gomock.Any()).Return(nil, errors.New("database error")).Times(1)

	mockRepo := testhelpers.NewMockExchangeRateRepository(store)
	exchangeRateService := NewExchangeRateService(mockRepo)

	// Test the service
	result, err := exchangeRateService.ListExchangeRates(context.Background())

	// Assertions
	require.Error(t, err)
	require.Equal(t, "database error", err.Error())
	require.Empty(t, result.ExchangeRates)
}

func TestListExchangeRatesService_EmptyResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// Expect empty result
	store.EXPECT().ListExchangeRates(gomock.Any()).Return([]db.ExchangeRate{}, nil).Times(1)

	mockRepo := testhelpers.NewMockExchangeRateRepository(store)
	exchangeRateService := NewExchangeRateService(mockRepo)

	// Test the service
	result, err := exchangeRateService.ListExchangeRates(context.Background())

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.ExchangeRates, 0)
	require.Equal(t, 0, result.Total)
}
