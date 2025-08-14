package users

import (
	"lemfi/simplebank/internal/middleware"
	"lemfi/simplebank/pkg/errorResponse"
	"lemfi/simplebank/pkg/responseHandler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *UserController) GetUserController(ctx *gin.Context) {
	// Get user info from JWT token context
	userClaims := middleware.ContextGetUser(ctx)
	if userClaims == nil || userClaims.Username == "" {
		errorResponse.UnAuthorizedRequestResponse(ctx)
		return
	}

	// Call service to get user details
	response, err := c.userService.GetUser(userClaims.Username)
	if err != nil {
		errorResponse.ServerErrorResponse(ctx, err)
		return
	}

	responseData := responseHandler.Envelope{
		"user": response,
	}

	err = responseHandler.WriteJSON(ctx.Writer, http.StatusOK, responseData, nil)
	if err != nil {
		errorResponse.ServerErrorResponse(ctx, err)
		return
	}
}
