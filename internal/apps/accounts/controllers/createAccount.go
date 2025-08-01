package accounts

import (
	"net/http"

	errorResponse "lemfi/simplebank/pkg/errorResponse"
	"lemfi/simplebank/pkg/requestHandler"
	"lemfi/simplebank/pkg/responseHandler"

	requests "lemfi/simplebank/internal/apps/accounts/requests"
	accountValidation "lemfi/simplebank/internal/apps/accounts/validationMessages"

	"github.com/gin-gonic/gin"
)

func (accountController *AccountController) CreateAccountController(c *gin.Context) {
	var req requests.CreateAccountRequest

	err := requestHandler.ReadJSONGin(c, &req, accountValidation.CreateAccountValidationMessages)
	if err != nil {
		errorResponse.BadRequestResponse(c, err)
		return
	}

	data := responseHandler.Envelope{
		"account": "Account created successfully",
	}

	err = responseHandler.WriteJSON(c.Writer, http.StatusOK, data, nil)
	if err != nil {
		errorResponse.ServerErrorResponse(c, err)
		return
	}
}
