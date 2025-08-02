package exchangeRates

import (
	"bytes"
	"encoding/json"
	"errors"
	mockdb "lemfi/simplebank/db/mock"
	db "lemfi/simplebank/db/sqlc"
	services "lemfi/simplebank/internal/apps/exchangeRates/services"
	testhelpers "lemfi/simplebank/internal/apps/exchangeRates/testHelpers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetExchangeRateHTTP_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// Create test data
	expectedRate := db.ExchangeRate{
		ID:           1,
		FromCurrency: "USD",
		ToCurrency:   "EUR",
		Rate:         decimal.NewFromFloat(0.85),
	}

	// Expect exchange rate to be fetched
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

	// Create request body
	requestBody := map[string]interface{}{
		"from_currency": "USD",
		"to_currency":   "XYZ", // Non-existent currency pair
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
	require.Contains(t, response, "amount_to_send")
	require.Contains(t, response, "amount_to_receive")
	require.Contains(t, response, "can_transact")
	require.Contains(t, response, "message")

	// Verify exchange rate details
	exchangeRate := response["exchange_rate"].(map[string]interface{})
	require.Equal(t, "USD", exchangeRate["from_currency"])
	require.Equal(t, "EUR", exchangeRate["to_currency"])
	require.Equal(t, "0.85", exchangeRate["rate"])
}
