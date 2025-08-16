package users

import (
	"lemfi/simplebank/config"
	"lemfi/simplebank/internal/apps/core"
	requests "lemfi/simplebank/internal/apps/users/requests"
	userValidation "lemfi/simplebank/internal/apps/users/validationMessages"
	errorResponse "lemfi/simplebank/pkg/errorResponse"
	"lemfi/simplebank/pkg/requestHandler"
	"lemfi/simplebank/pkg/responseHandler"

	"github.com/gin-gonic/gin"
)

func (userController *UserController) LogoutController(c *gin.Context) {
	config.Logger.Info("Processing logout request", "method", "POST", "endpoint", "/users/logout")

	var req requests.LogoutRequest

	err := requestHandler.ReadJSONGin(c, &req, userValidation.LogoutValidationMessages)
	if err != nil {
		config.Logger.Error("Failed to read logout request", "error", err.Error())
		errorResponse.BadRequestResponse(c, err)
		return
	}

	config.Logger.Info("Logout request validated successfully")

	err = userController.userService.Logout(req)
	if err != nil {
		config.Logger.Error("Failed to logout user", "error", err.Error())
		if clientErr, isClient := core.IsClientError(err); isClient {
			errorResponse.BadRequestResponse(c, clientErr)
		} else {
			errorResponse.ServerErrorResponse(c, err)
		}
		return
	}

	config.Logger.Info("User logged out successfully")

	responseData := responseHandler.Envelope{
		"message": "User logged out successfully",
	}

	err = responseHandler.WriteJSON(c.Writer, 200, responseData, nil)
	if err != nil {
		config.Logger.Error("Failed to write JSON response", "error", err.Error())
		errorResponse.ServerErrorResponse(c, err)
		return
	}

	config.Logger.Info("Logout completed successfully")
}
