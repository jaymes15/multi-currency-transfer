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

func (userController *UserController) CreateUserController(c *gin.Context) {
	config.Logger.Info("Creating new user", "method", "POST", "endpoint", "/users")

	var req requests.CreateUserRequest

	err := requestHandler.ReadJSONGin(c, &req, userValidation.CreateUserValidationMessages)
	if err != nil {
		config.Logger.Error("Failed to read user request", "error", err.Error())
		errorResponse.BadRequestResponse(c, err)
		return
	}

	config.Logger.Info("User request validated successfully", "username", req.Username, "email", req.Email)

	user, err := userController.userService.CreateUser(req)
	if err != nil {
		config.Logger.Error("Failed to create user", "error", err.Error(), "username", req.Username)
		if clientErr, isClient := core.IsClientError(err); isClient {
			errorResponse.BadRequestResponse(c, clientErr)
		} else {
			errorResponse.ServerErrorResponse(c, err)
		}
		return
	}

	config.Logger.Info("User created successfully", "username", user.Username, "email", user.Email)

	response := responseHandler.Envelope{
		"user": user,
	}

	err = responseHandler.WriteJSON(c.Writer, http.StatusCreated, response, nil)
	if err != nil {
		config.Logger.Error("Failed to write JSON response", "error", err.Error())
		errorResponse.ServerErrorResponse(c, err)
		return
	}

	config.Logger.Info("User creation completed successfully", "username", user.Username)
}
