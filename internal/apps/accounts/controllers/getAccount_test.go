package accounts

import (
	"fmt"
	mockdb "lemfi/simplebank/db/mock"
	db "lemfi/simplebank/db/sqlc"
	services "lemfi/simplebank/internal/apps/accounts/services"
	testhelpers "lemfi/simplebank/internal/apps/accounts/testHelpers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccountHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// Create test data
	testAccounts := []db.Account{
		{
			ID:       1,
			Owner:    "test_owner_1",
			Balance:  1000,
			Currency: "USD",
		},
		{
			ID:       2,
			Owner:    "test_owner_2",
			Balance:  2000,
			Currency: "EUR",
		},
	}

	store.EXPECT().ListAllAccounts(gomock.Any(), gomock.Any()).Return(testAccounts, nil).Times(1)

	// Create mock repository that wraps the mock store
	mockRepo := testhelpers.NewMockAccountRepository(store)

	accountService := services.NewAccountService(mockRepo)
	accountController := NewAccountController(accountService)

	// Create a new Gin router for testing
	router := gin.New()

	// Set up the route
	router.GET("/accounts", accountController.GetAccountsController)

	// Create HTTP test recorder
	recorder := httptest.NewRecorder()

	// Create HTTP request
	url := "/accounts"
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	// Serve the request
	router.ServeHTTP(recorder, request)

	// Assert the response
	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetAccountHTTP_EmptyList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// Expect empty list
	store.EXPECT().ListAllAccounts(gomock.Any(), gomock.Any()).Return([]db.Account{}, nil).Times(1)

	mockRepo := testhelpers.NewMockAccountRepository(store)
	accountService := services.NewAccountService(mockRepo)
	accountController := NewAccountController(accountService)

	router := gin.New()
	router.GET("/accounts", accountController.GetAccountsController)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/accounts", nil)
	require.NoError(t, err)

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetAccountHTTP_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// Expect database error
	store.EXPECT().ListAllAccounts(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("database connection failed")).Times(1)

	mockRepo := testhelpers.NewMockAccountRepository(store)
	accountService := services.NewAccountService(mockRepo)
	accountController := NewAccountController(accountService)

	router := gin.New()
	router.GET("/accounts", accountController.GetAccountsController)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/accounts", nil)
	require.NoError(t, err)

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestGetAccountHTTP_WrongMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	mockRepo := testhelpers.NewMockAccountRepository(store)
	accountService := services.NewAccountService(mockRepo)
	accountController := NewAccountController(accountService)

	router := gin.New()
	router.GET("/accounts", accountController.GetAccountsController)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "/accounts", nil)
	require.NoError(t, err)

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestGetAccountHTTP_ResponseBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	testAccounts := []db.Account{
		{
			ID:       1,
			Owner:    "test_owner_1",
			Balance:  1000,
			Currency: "USD",
		},
		{
			ID:       2,
			Owner:    "test_owner_2",
			Balance:  2000,
			Currency: "EUR",
		},
	}

	store.EXPECT().ListAllAccounts(gomock.Any(), gomock.Any()).Return(testAccounts, nil).Times(1)

	mockRepo := testhelpers.NewMockAccountRepository(store)
	accountService := services.NewAccountService(mockRepo)
	accountController := NewAccountController(accountService)

	router := gin.New()
	router.GET("/accounts", accountController.GetAccountsController)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/accounts", nil)
	require.NoError(t, err)

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Contains(t, recorder.Header().Get("Content-Type"), "application/json")

	// Verify response body contains expected data
	responseBody := recorder.Body.String()
	require.Contains(t, responseBody, "test_owner_1")
	require.Contains(t, responseBody, "test_owner_2")
	require.Contains(t, responseBody, "USD")
	require.Contains(t, responseBody, "EUR")
}
