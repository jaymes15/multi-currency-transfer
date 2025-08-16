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

func (userController *UserController) RefreshTokenController(c *gin.Context) {
	config.Logger.Info("Processing refresh token request", "method", "POST", "endpoint", "/users/refresh")

	var req requests.RefreshTokenRequest

	err := requestHandler.ReadJSONGin(c, &req, userValidation.RefreshTokenValidationMessages)
	if err != nil {
		config.Logger.Error("Failed to read refresh token request", "error", err.Error())
		errorResponse.BadRequestResponse(c, err)
		return
	}

	config.Logger.Info("Refresh token request validated successfully")

	response, err := userController.userService.RefreshToken(req)
	if err != nil {
		config.Logger.Error("Failed to refresh token", "error", err.Error())
		if clientErr, isClient := core.IsClientError(err); isClient {
			errorResponse.BadRequestResponse(c, clientErr)
		} else {
			errorResponse.ServerErrorResponse(c, err)
		}
		return
	}

	config.Logger.Info("Token refreshed successfully")

	responseData := responseHandler.Envelope{
		"access_token":             response.AccessToken,
		"access_token_expires_at":  response.AccessTokenExpiresAt,
		"refresh_token":            response.RefreshToken,
		"refresh_token_expires_at": response.RefreshTokenExpiresAt,
	}

	err = responseHandler.WriteJSON(c.Writer, 200, responseData, nil)
	if err != nil {
		config.Logger.Error("Failed to write JSON response", "error", err.Error())
		errorResponse.ServerErrorResponse(c, err)
		return
	}

	config.Logger.Info("Token refresh completed successfully")
}
