package exchangeRates

import (
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

func TestListExchangeRatesHTTP_Success(t *testing.T) {
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
	exchangeRateService := services.NewExchangeRateService(mockRepo)
	exchangeRateController := NewExchangeRateController(exchangeRateService)

	router := gin.New()
	router.GET("/exchange-rates", exchangeRateController.ListExchangeRatesController)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/exchange-rates", nil)
	require.NoError(t, err)

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestListExchangeRatesHTTP_EmptyList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// Expect empty list
	store.EXPECT().ListExchangeRates(gomock.Any()).Return([]db.ExchangeRate{}, nil).Times(1)

	mockRepo := testhelpers.NewMockExchangeRateRepository(store)
	exchangeRateService := services.NewExchangeRateService(mockRepo)
	exchangeRateController := NewExchangeRateController(exchangeRateService)

	router := gin.New()
	router.GET("/exchange-rates", exchangeRateController.ListExchangeRatesController)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/exchange-rates", nil)
	require.NoError(t, err)

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestListExchangeRatesHTTP_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// Expect database error
	store.EXPECT().ListExchangeRates(gomock.Any()).Return(nil, errors.New("database error")).Times(1)

	mockRepo := testhelpers.NewMockExchangeRateRepository(store)
	exchangeRateService := services.NewExchangeRateService(mockRepo)
	exchangeRateController := NewExchangeRateController(exchangeRateService)

	router := gin.New()
	router.GET("/exchange-rates", exchangeRateController.ListExchangeRatesController)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/exchange-rates", nil)
	require.NoError(t, err)

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestListExchangeRatesHTTP_WrongMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	mockRepo := testhelpers.NewMockExchangeRateRepository(store)
	exchangeRateService := services.NewExchangeRateService(mockRepo)
	exchangeRateController := NewExchangeRateController(exchangeRateService)

	router := gin.New()
	router.GET("/exchange-rates", exchangeRateController.ListExchangeRatesController)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "/exchange-rates", nil)
	require.NoError(t, err)

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestListExchangeRatesHTTP_ResponseBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

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

	store.EXPECT().ListExchangeRates(gomock.Any()).Return(expectedRates, nil).Times(1)

	mockRepo := testhelpers.NewMockExchangeRateRepository(store)
	exchangeRateService := services.NewExchangeRateService(mockRepo)
	exchangeRateController := NewExchangeRateController(exchangeRateService)

	router := gin.New()
	router.GET("/exchange-rates", exchangeRateController.ListExchangeRatesController)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/exchange-rates", nil)
	require.NoError(t, err)

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Contains(t, recorder.Header().Get("Content-Type"), "application/json")

	// Verify response body contains expected data
	responseBody := recorder.Body.String()
	require.Contains(t, responseBody, "USD")
	require.Contains(t, responseBody, "EUR")
	require.Contains(t, responseBody, "0.85")
	require.Contains(t, responseBody, "1.18")

	// Parse response to verify structure
	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify response structure
	require.Contains(t, response, "exchange_rates")
	require.Contains(t, response, "total")
	rates := response["exchange_rates"].([]interface{})
	require.Len(t, rates, 2)
	require.Equal(t, float64(2), response["total"])
}
