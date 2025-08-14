package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	mockdb "lemfi/simplebank/db/mock"
	db "lemfi/simplebank/db/sqlc"
	services "lemfi/simplebank/internal/apps/users/services"
	testhelpers "lemfi/simplebank/internal/apps/users/testHelpers"
	"lemfi/simplebank/pkg/cipher"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// eqCreateUserParamsMatcher is a custom matcher for CreateUserParams
type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	// Check if the password hash is valid
	err := cipher.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	// Set the hashed password to match for comparison
	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

// EqCreateUserParams creates a matcher for CreateUserParams
func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserHTTP_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// Create test data
	createRequest := map[string]interface{}{
		"username":  "testuser",
		"password":  "password123",
		"full_name": "Test User",
		"email":     "test@example.com",
	}

	expectedUser := db.CreateUserRow{
		Username:  "testuser",
		FullName:  "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
	}

	// Create expected parameters for matching
	expectedParams := db.CreateUserParams{
		Username:       "testuser",
		FullName:       "Test User",
		Email:          "test@example.com",
		HashedPassword: "", // This will be set by the matcher
	}

	// Expect user creation with proper parameter matching
	store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(expectedParams, "password123")).Return(expectedUser, nil).Times(1)

	mockRepo := testhelpers.NewMockUserRepository(store)
	userService := services.NewUserService(mockRepo, nil) // nil tokenMaker for now
	userController := NewUserController(userService, nil) // nil tokenMaker for now

	router := gin.New()
	router.POST("/users", userController.CreateUserController)

	// Create request body
	requestBody, err := json.Marshal(createRequest)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusCreated, recorder.Code)
}

func TestCreateUserHTTP_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	mockRepo := testhelpers.NewMockUserRepository(store)
	userService := services.NewUserService(mockRepo, nil) // nil tokenMaker for now
	userController := NewUserController(userService, nil) // nil tokenMaker for now

	router := gin.New()
	router.POST("/users", userController.CreateUserController)

	// Create invalid request (missing required fields)
	invalidRequest := map[string]interface{}{
		"username": "testuser",
		// missing "password", "full_name", "email" fields
	}

	requestBody, err := json.Marshal(invalidRequest)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCreateUserHTTP_InvalidEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	mockRepo := testhelpers.NewMockUserRepository(store)
	userService := services.NewUserService(mockRepo, nil) // nil tokenMaker for now
	userController := NewUserController(userService, nil) // nil tokenMaker for now

	router := gin.New()
	router.POST("/users", userController.CreateUserController)

	// Create invalid request (invalid email format)
	invalidRequest := map[string]interface{}{
		"username":  "testuser",
		"password":  "password123",
		"full_name": "Test User",
		"email":     "invalid-email",
	}

	requestBody, err := json.Marshal(invalidRequest)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCreateUserHTTP_ShortPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	mockRepo := testhelpers.NewMockUserRepository(store)
	userService := services.NewUserService(mockRepo, nil) // nil tokenMaker for now
	userController := NewUserController(userService, nil) // nil tokenMaker for now

	router := gin.New()
	router.POST("/users", userController.CreateUserController)

	// Create invalid request (password too short)
	invalidRequest := map[string]interface{}{
		"username":  "testuser",
		"password":  "123",
		"full_name": "Test User",
		"email":     "test@example.com",
	}

	requestBody, err := json.Marshal(invalidRequest)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCreateUserHTTP_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	createRequest := map[string]interface{}{
		"username":  "testuser",
		"password":  "password123",
		"full_name": "Test User",
		"email":     "test@example.com",
	}

	// Create expected parameters for matching
	expectedParams := db.CreateUserParams{
		Username:       "testuser",
		FullName:       "Test User",
		Email:          "test@example.com",
		HashedPassword: "", // This will be set by the matcher
	}

	// Expect database error with proper parameter matching
	store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(expectedParams, "password123")).Return(db.CreateUserRow{}, fmt.Errorf("database error")).Times(1)

	mockRepo := testhelpers.NewMockUserRepository(store)
	userService := services.NewUserService(mockRepo, nil) // nil tokenMaker for now
	userController := NewUserController(userService, nil) // nil tokenMaker for now

	router := gin.New()
	router.POST("/users", userController.CreateUserController)

	requestBody, err := json.Marshal(createRequest)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
}
