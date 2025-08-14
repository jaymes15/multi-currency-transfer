package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestLoginUserHTTP_InvalidRequest(t *testing.T) {
	// Create a simple test without mocks for now
	userController := &UserController{
		userService: nil, // nil for now
		tokenMaker:  nil, // nil for now
	}

	router := gin.New()
	router.POST("/users/login", userController.LoginUserController)

	// Create invalid request (missing required fields)
	invalidRequest := map[string]interface{}{
		"username": "testuser",
		// missing "password" field
	}

	requestBody, err := json.Marshal(invalidRequest)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	// Should return bad request due to missing password
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}
