package users

import (
	"lemfi/simplebank/config"
	"lemfi/simplebank/internal/apps/core"
	"net/http"

	errorResponse "lemfi/simplebank/pkg/errorResponse"
	"lemfi/simplebank/pkg/requestHandler"
	"lemfi/simplebank/pkg/responseHandler"

	requests "lemfi/simplebank/internal/apps/users/requests"
	userValidation "lemfi/simplebank/internal/apps/users/validationMessages"

	"github.com/gin-gonic/gin"
)

func (userController *UserController) LoginUserController(c *gin.Context) {
	config.Logger.Info("User login attempt", "method", "POST", "endpoint", "/users/login")

	var req requests.LoginUserRequest

	err := requestHandler.ReadJSONGin(c, &req, userValidation.LoginUserValidationMessages)
	if err != nil {
		config.Logger.Error("Failed to read login request", "error", err.Error())
		errorResponse.BadRequestResponse(c, err)
		return
	}

	config.Logger.Info("Login request validated successfully", "username", req.Username)

	response, err := userController.userService.LoginUser(req)
	if err != nil {
		config.Logger.Error("Failed to login user", "error", err.Error(), "username", req.Username)
		if clientErr, isClient := core.IsClientError(err); isClient {
			errorResponse.BadRequestResponse(c, clientErr)
		} else {
			errorResponse.ServerErrorResponse(c, err)
		}
		return
	}

	config.Logger.Info("User logged in successfully", "username", req.Username)

	responseData := responseHandler.Envelope{
		"tokens": response,
	}

	err = responseHandler.WriteJSON(c.Writer, http.StatusOK, responseData, nil)
	if err != nil {
		config.Logger.Error("Failed to write JSON response", "error", err.Error())
		errorResponse.ServerErrorResponse(c, err)
		return
	}

	config.Logger.Info("User login completed successfully", "username", req.Username)
}
