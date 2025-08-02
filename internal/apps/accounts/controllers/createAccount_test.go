package accounts

import (
	"bytes"
	"encoding/json"
	"fmt"
	mockdb "lemfi/simplebank/db/mock"
	db "lemfi/simplebank/db/sqlc"
	services "lemfi/simplebank/internal/apps/accounts/services"
	testhelpers "lemfi/simplebank/internal/apps/accounts/testHelpers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateAccountHTTP_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// Create test data
	createRequest := map[string]interface{}{
		"owner":    "test_owner",
		"currency": "USD",
	}

	expectedAccount := db.Account{
		ID:       1,
		Owner:    "test_owner",
		Balance:  decimal.Zero,
		Currency: "USD",
	}

	// Expect account creation
	store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(expectedAccount, nil).Times(1)

	mockRepo := testhelpers.NewMockAccountRepository(store)
	accountService := services.NewAccountService(mockRepo)
	accountController := NewAccountController(accountService)

	router := gin.New()
	router.POST("/accounts", accountController.CreateAccountController)

	// Create request body
	requestBody, err := json.Marshal(createRequest)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusCreated, recorder.Code)
}

func TestCreateAccountHTTP_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	mockRepo := testhelpers.NewMockAccountRepository(store)
	accountService := services.NewAccountService(mockRepo)
	accountController := NewAccountController(accountService)

	router := gin.New()
	router.POST("/accounts", accountController.CreateAccountController)

	// Create invalid request (missing required fields)
	invalidRequest := map[string]interface{}{
		"currency": "USD",
		// missing "owner" field
	}

	requestBody, err := json.Marshal(invalidRequest)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCreateAccountHTTP_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	createRequest := map[string]interface{}{
		"owner":    "test_owner",
		"currency": "USD",
	}

	// Expect database error
	store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(db.Account{}, fmt.Errorf("database error")).Times(1)

	mockRepo := testhelpers.NewMockAccountRepository(store)
	accountService := services.NewAccountService(mockRepo)
	accountController := NewAccountController(accountService)

	router := gin.New()
	router.POST("/accounts", accountController.CreateAccountController)

	requestBody, err := json.Marshal(createRequest)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
}
