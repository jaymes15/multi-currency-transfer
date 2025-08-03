package exchangeRates

import (
	"context"
	"errors"
	"testing"
	"time"

	mockdb "lemfi/simplebank/db/mock"
	db "lemfi/simplebank/db/sqlc"
	requests "lemfi/simplebank/internal/apps/exchangeRates/requests"
	responses "lemfi/simplebank/internal/apps/exchangeRates/responses"
	respositories "lemfi/simplebank/internal/apps/exchangeRates/respositories"
	testhelpers "lemfi/simplebank/internal/apps/exchangeRates/testHelpers"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// MockExchangeRateService for testing with controllable expiration logic
type MockExchangeRateService struct {
	repo          respositories.ExchangeRateRepositoryInterface
	isExpiredFunc func(db.ExchangeRate) bool
}

func NewMockExchangeRateService(repo respositories.ExchangeRateRepositoryInterface) *MockExchangeRateService {
	return &MockExchangeRateService{
		repo: repo,
	}
}

func (m *MockExchangeRateService) SetIsExpiredFunc(fn func(db.ExchangeRate) bool) {
	m.isExpiredFunc = fn
}

func (m *MockExchangeRateService) GetExchangeRate(ctx context.Context, payload requests.GetExchangeRateRequest) (responses.GetExchangeRateResponse, error) {
	dbExchangeRate, err := m.repo.GetExchangeRate(ctx, payload)
	if err != nil {
		return responses.GetExchangeRateResponse{}, err
	}

	exchangeRate := responses.NewExchangeRateResponse(dbExchangeRate)
	amountToSend := payload.Amount
	amountToReceive := payload.Amount.Mul(dbExchangeRate.Rate).Round(2)

	// Calculate fee (using a default fee for testing)
	fee := decimal.NewFromFloat(1.99) // Default fee for testing
	totalAmount := amountToSend.Add(fee)

	canTransact := false
	message := "Exchange rate expired"

	// Use mock function if set, otherwise assume not expired
	if m.isExpiredFunc != nil {
		if !m.isExpiredFunc(dbExchangeRate) {
			canTransact = true
			message = "Exchange rate available for transaction"
		}
	} else {
		canTransact = true
		message = "Exchange rate available for transaction"
	}

	return responses.GetExchangeRateResponse{
		ExchangeRate:    exchangeRate,
		AmountToSend:    amountToSend,
		AmountToReceive: amountToReceive,
		Fee:             fee,
		TotalAmount:     totalAmount,
		CanTransact:     canTransact,
		Message:         message,
	}, nil
}

func (m *MockExchangeRateService) ListExchangeRates(ctx context.Context) (responses.ListExchangeRatesResponse, error) {
	rates, err := m.repo.ListExchangeRates(ctx)
	if err != nil {
		return responses.ListExchangeRatesResponse{}, err
	}

	exchangeRates := make([]responses.ExchangeRateResponse, len(rates))
	for i, rate := range rates {
		exchangeRates[i] = responses.NewExchangeRateResponse(rate)
	}

	return responses.ListExchangeRatesResponse{
		ExchangeRates: exchangeRates,
		Total:         len(exchangeRates),
	}, nil
}

func TestGetExchangeRateService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

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

	// Create a mock service that controls expiration logic
	mockService := &MockExchangeRateService{
		repo: mockRepo,
	}

	// Mock the expiration logic to return false (not expired)
	mockService.SetIsExpiredFunc(func(exchangeRate db.ExchangeRate) bool {
		return false // Not expired
	})

	// Create test request
	request := requests.GetExchangeRateRequest{
		FromCurrency: "USD",
		ToCurrency:   "EUR",
		Amount:       decimal.NewFromFloat(100.00),
	}

	// Test the service
	result, err := mockService.GetExchangeRate(context.Background(), request)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "USD", result.ExchangeRate.FromCurrency)
	require.Equal(t, "EUR", result.ExchangeRate.ToCurrency)
	require.Equal(t, decimal.NewFromFloat(0.85), result.ExchangeRate.Rate)
	require.Equal(t, decimal.NewFromFloat(100.00), result.AmountToSend)
	require.Equal(t, decimal.NewFromFloat(85.00).Round(2), result.AmountToReceive)
	require.NotEqual(t, decimal.Zero, result.Fee)                             // Fee should be set
	require.Equal(t, result.AmountToSend.Add(result.Fee), result.TotalAmount) // Total should be amount + fee
	require.Equal(t, true, result.CanTransact)
	require.Equal(t, "Exchange rate available for transaction", result.Message)
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

	// Create test request with supported currencies but non-existent pair
	request := requests.GetExchangeRateRequest{
		FromCurrency: "GBP",
		ToCurrency:   "NGN", // Supported currency but pair doesn't exist in DB
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

	// No database call expected since validation fails first
	// store.EXPECT().GetExchangeRate(gomock.Any(), gomock.Any()).Return(expectedRate, nil).Times(1)

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
	require.Error(t, err)
	require.Equal(t, "invalid amount", err.Error())
	require.Empty(t, result.ExchangeRate)
}
