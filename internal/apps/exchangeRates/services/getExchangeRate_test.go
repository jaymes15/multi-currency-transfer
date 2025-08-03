package exchangeRates

import (
	"context"
	"errors"
	mockdb "lemfi/simplebank/db/mock"
	db "lemfi/simplebank/db/sqlc"
	requests "lemfi/simplebank/internal/apps/exchangeRates/requests"
	testhelpers "lemfi/simplebank/internal/apps/exchangeRates/testHelpers"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetExchangeRateService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// Create test data
	expectedRate := db.ExchangeRate{
		ID:           1,
		FromCurrency: "USD",
		ToCurrency:   "EUR",
		Rate:         decimal.NewFromFloat(0.85),
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	// Expect exchange rate to be fetched
	store.EXPECT().GetExchangeRate(gomock.Any(), gomock.Any()).Return(expectedRate, nil).Times(1)

	mockRepo := testhelpers.NewMockExchangeRateRepository(store)
	exchangeRateService := NewExchangeRateService(mockRepo)

	// Create test request
	request := requests.GetExchangeRateRequest{
		FromCurrency: "USD",
		ToCurrency:   "EUR",
		Amount:       decimal.NewFromFloat(100.00),
	}

	// Test the service
	result, err := exchangeRateService.GetExchangeRate(context.Background(), request)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "USD", result.ExchangeRate.FromCurrency)
	require.Equal(t, "EUR", result.ExchangeRate.ToCurrency)
	require.Equal(t, decimal.NewFromFloat(0.85), result.ExchangeRate.Rate)
	require.Equal(t, decimal.NewFromFloat(100.00), result.AmountToSend)
	require.Equal(t, decimal.NewFromFloat(85.00).Round(2), result.AmountToReceive)
}

func TestGetExchangeRateService_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// Expect database error
	store.EXPECT().GetExchangeRate(gomock.Any(), gomock.Any()).Return(db.ExchangeRate{}, errors.New("database error")).Times(1)

	mockRepo := testhelpers.NewMockExchangeRateRepository(store)
	exchangeRateService := NewExchangeRateService(mockRepo)

	// Create test request
	request := requests.GetExchangeRateRequest{
		FromCurrency: "USD",
		ToCurrency:   "EUR",
		Amount:       decimal.NewFromFloat(100.00),
	}

	// Test the service
	result, err := exchangeRateService.GetExchangeRate(context.Background(), request)

	// Assertions
	require.Error(t, err)
	require.Equal(t, "database error", err.Error())
	require.Empty(t, result.ExchangeRate)
}

func TestGetExchangeRateService_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// Expect not found error
	store.EXPECT().GetExchangeRate(gomock.Any(), gomock.Any()).Return(db.ExchangeRate{}, errors.New("no rows in result set")).Times(1)

	mockRepo := testhelpers.NewMockExchangeRateRepository(store)
	exchangeRateService := NewExchangeRateService(mockRepo)

	// Create test request
	request := requests.GetExchangeRateRequest{
		FromCurrency: "USD",
		ToCurrency:   "XYZ", // Non-existent currency pair
		Amount:       decimal.NewFromFloat(100.00),
	}

	// Test the service
	result, err := exchangeRateService.GetExchangeRate(context.Background(), request)

	// Assertions
	require.Error(t, err)
	require.Equal(t, "no rows in result set", err.Error())
	require.Empty(t, result.ExchangeRate)
}

func TestGetExchangeRateService_ZeroAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// Expect exchange rate to be fetched
	expectedRate := db.ExchangeRate{
		ID:           1,
		FromCurrency: "USD",
		ToCurrency:   "EUR",
		Rate:         decimal.NewFromFloat(0.85),
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	store.EXPECT().GetExchangeRate(gomock.Any(), gomock.Any()).Return(expectedRate, nil).Times(1)

	mockRepo := testhelpers.NewMockExchangeRateRepository(store)
	exchangeRateService := NewExchangeRateService(mockRepo)

	// Create test request with zero amount
	request := requests.GetExchangeRateRequest{
		FromCurrency: "USD",
		ToCurrency:   "EUR",
		Amount:       decimal.Zero,
	}

	// Test the service
	result, err := exchangeRateService.GetExchangeRate(context.Background(), request)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, decimal.Zero, result.AmountToSend)
	require.Equal(t, decimal.Zero.Round(2), result.AmountToReceive)
}
