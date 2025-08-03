package exchangeRates

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	mockdb "lemfi/simplebank/db/mock"
	db "lemfi/simplebank/db/sqlc"
	requests "lemfi/simplebank/internal/apps/exchangeRates/requests"
	responses "lemfi/simplebank/internal/apps/exchangeRates/responses"
	respositories "lemfi/simplebank/internal/apps/exchangeRates/respositories"
	services "lemfi/simplebank/internal/apps/exchangeRates/services"
	testhelpers "lemfi/simplebank/internal/apps/exchangeRates/testHelpers"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
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

func TestGetExchangeRateHTTP_Success(t *testing.T) {
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
	mockService := NewMockExchangeRateService(mockRepo) // Using the local mock service

	// Mock the expiration logic to return false (not expired)
	mockService.SetIsExpiredFunc(func(exchangeRate db.ExchangeRate) bool {
		return false // Not expired
	})

	exchangeRateController := NewExchangeRateController(mockService)

	router := gin.New()
	router.POST("/exchange-rates/calculate", exchangeRateController.GetExchangeRateController)

	// Create request body
	requestBody := map[string]interface{}{
		"from_currency": "USD",
		"to_currency":   "EUR",
		"amount":        "100.00",
	}

	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "/exchange-rates/calculate", bytes.NewBuffer(body))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	// Assert the response structure
	require.NotNil(t, response["exchange_rate"])
	exchangeRateData := response["exchange_rate"].(map[string]interface{})
	require.Equal(t, "100", exchangeRateData["amount_to_send"])
	require.Equal(t, "85", exchangeRateData["amount_to_receive"])
	require.NotNil(t, exchangeRateData["fee"])
	require.NotNil(t, exchangeRateData["total_amount"])
	require.Equal(t, true, exchangeRateData["can_transact"])
	require.Equal(t, "Exchange rate available for transaction", exchangeRateData["message"])
}

func TestGetExchangeRateHTTP_ExpiredRate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	expectedRate := db.ExchangeRate{
		ID:           1,
		FromCurrency: "USD",
		ToCurrency:   "EUR",
		Rate:         decimal.NewFromFloat(0.85),
		CreatedAt:    pgtype.Timestamptz{Time: time.Now().Add(-24 * time.Hour), Valid: true}, // Old rate
		UpdatedAt:    pgtype.Timestamptz{Time: time.Now().Add(-24 * time.Hour), Valid: true}, // Old rate
	}

	store.EXPECT().GetExchangeRate(gomock.Any(), gomock.Any()).Return(expectedRate, nil).Times(1)

	mockRepo := testhelpers.NewMockExchangeRateRepository(store)
	mockService := NewMockExchangeRateService(mockRepo) // Using the local mock service

	// Mock the expiration logic to return true (expired)
	mockService.SetIsExpiredFunc(func(exchangeRate db.ExchangeRate) bool {
		return true // Expired
	})

	exchangeRateController := NewExchangeRateController(mockService)

	router := gin.New()
	router.POST("/exchange-rates/calculate", exchangeRateController.GetExchangeRateController)

	// Create request body
	requestBody := map[string]interface{}{
		"from_currency": "USD",
		"to_currency":   "EUR",
		"amount":        "100.00",
	}

	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "/exchange-rates/calculate", bytes.NewBuffer(body))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	// Assert the response structure
	require.NotNil(t, response["exchange_rate"])
	exchangeRateData := response["exchange_rate"].(map[string]interface{})
	require.Equal(t, "100", exchangeRateData["amount_to_send"])
	require.Equal(t, "85", exchangeRateData["amount_to_receive"])
	require.NotNil(t, exchangeRateData["fee"])
	require.NotNil(t, exchangeRateData["total_amount"])
	require.Equal(t, false, exchangeRateData["can_transact"])
	require.Equal(t, "Exchange rate expired", exchangeRateData["message"])
}

func TestGetExchangeRateHTTP_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	mockRepo := testhelpers.NewMockExchangeRateRepository(store)
	exchangeRateService := services.NewExchangeRateService(mockRepo)
	exchangeRateController := NewExchangeRateController(exchangeRateService)

	router := gin.New()
	router.POST("/exchange-rates/calculate", exchangeRateController.GetExchangeRateController)

	// Create invalid request (missing required fields)
	invalidRequest := map[string]interface{}{
		"from_currency": "USD",
		// missing "to_currency" and "amount" fields
	}

	body, err := json.Marshal(invalidRequest)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "/exchange-rates/calculate", bytes.NewBuffer(body))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestGetExchangeRateHTTP_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// Expect exchange rate not found
	store.EXPECT().GetExchangeRate(gomock.Any(), gomock.Any()).Return(db.ExchangeRate{}, errors.New("no rows in result set")).Times(1)

	mockRepo := testhelpers.NewMockExchangeRateRepository(store)
	exchangeRateService := services.NewExchangeRateService(mockRepo)
	exchangeRateController := NewExchangeRateController(exchangeRateService)

	router := gin.New()
	router.POST("/exchange-rates/calculate", exchangeRateController.GetExchangeRateController)

	// Create request body with supported currencies but non-existent pair
	requestBody := map[string]interface{}{
		"from_currency": "GBP",
		"to_currency":   "NGN", // Supported currency but pair doesn't exist in DB
		"amount":        "100.00",
	}

	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "/exchange-rates/calculate", bytes.NewBuffer(body))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	// The error is returned as 500 because it's not a client error
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestGetExchangeRateHTTP_WrongMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	mockRepo := testhelpers.NewMockExchangeRateRepository(store)
	exchangeRateService := services.NewExchangeRateService(mockRepo)
	exchangeRateController := NewExchangeRateController(exchangeRateService)

	router := gin.New()
	router.POST("/exchange-rates/calculate", exchangeRateController.GetExchangeRateController)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/exchange-rates/calculate", nil)
	require.NoError(t, err)

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestGetExchangeRateHTTP_ResponseBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	expectedRate := db.ExchangeRate{
		ID:           1,
		FromCurrency: "USD",
		ToCurrency:   "EUR",
		Rate:         decimal.NewFromFloat(0.85),
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	store.EXPECT().GetExchangeRate(gomock.Any(), gomock.Any()).Return(expectedRate, nil).Times(1)

	mockRepo := testhelpers.NewMockExchangeRateRepository(store)
	exchangeRateService := services.NewExchangeRateService(mockRepo)
	exchangeRateController := NewExchangeRateController(exchangeRateService)

	router := gin.New()
	router.POST("/exchange-rates/calculate", exchangeRateController.GetExchangeRateController)

	// Create request body
	requestBody := map[string]interface{}{
		"from_currency": "USD",
		"to_currency":   "EUR",
		"amount":        "100.00",
	}

	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "/exchange-rates/calculate", bytes.NewBuffer(body))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Contains(t, recorder.Header().Get("Content-Type"), "application/json")

	// Verify response body contains expected data
	responseBody := recorder.Body.String()
	require.Contains(t, responseBody, "USD")
	require.Contains(t, responseBody, "EUR")
	require.Contains(t, responseBody, "0.85")
	require.Contains(t, responseBody, "100")
	require.Contains(t, responseBody, "85")

	// Parse response to verify structure
	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify response structure
	require.Contains(t, response, "exchange_rate")

	// Verify exchange rate details
	exchangeRateData := response["exchange_rate"].(map[string]interface{})
	require.Contains(t, exchangeRateData, "amount_to_send")
	require.Contains(t, exchangeRateData, "amount_to_receive")
	require.Contains(t, exchangeRateData, "can_transact")
	require.Contains(t, exchangeRateData, "message")
	require.Contains(t, exchangeRateData, "fee")
	require.Contains(t, exchangeRateData, "total_amount")

	// Verify nested exchange rate details
	exchangeRate := exchangeRateData["exchange_rate"].(map[string]interface{})
	require.Equal(t, "USD", exchangeRate["from_currency"])
	require.Equal(t, "EUR", exchangeRate["to_currency"])
	require.Equal(t, "0.85", exchangeRate["rate"])
}
