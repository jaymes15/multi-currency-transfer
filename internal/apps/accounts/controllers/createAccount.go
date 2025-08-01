package accounts

import (
	"lemfi/simplebank/config"
	"net/http"

	errorResponse "lemfi/simplebank/pkg/errorResponse"
	"lemfi/simplebank/pkg/requestHandler"
	"lemfi/simplebank/pkg/responseHandler"

	requests "lemfi/simplebank/internal/apps/accounts/requests"
	accountValidation "lemfi/simplebank/internal/apps/accounts/validationMessages"

	"github.com/gin-gonic/gin"
)

func (accountController *AccountController) CreateAccountController(c *gin.Context) {
	config.Logger.Info("Creating new account", "method", "POST", "endpoint", "/accounts")

	var req requests.CreateAccountRequest

	err := requestHandler.ReadJSONGin(c, &req, accountValidation.CreateAccountValidationMessages)
	if err != nil {
		config.Logger.Error("Failed to read account request", "error", err.Error())
		errorResponse.BadRequestResponse(c, err)
		return
	}

	config.Logger.Info("Account request validated successfully", "owner", req.Owner, "currency", req.Currency)

	account, err := accountController.accountService.CreateAccount(req)
	if err != nil {
		config.Logger.Error("Failed to create account", "error", err.Error(), "owner", req.Owner)
		errorResponse.ServerErrorResponse(c, err)
		return
	}

	config.Logger.Info("Account created successfully", "accountID", account.ID, "owner", account.Owner, "currency", account.Currency)

	response := responseHandler.Envelope{
		"account": account,
	}

	err = responseHandler.WriteJSON(c.Writer, http.StatusCreated, response, nil)
	if err != nil {
		config.Logger.Error("Failed to write JSON response", "error", err.Error())
		errorResponse.ServerErrorResponse(c, err)
		return
	}

	config.Logger.Info("Account creation completed successfully", "accountID", account.ID)
}
