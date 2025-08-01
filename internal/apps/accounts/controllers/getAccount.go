package accounts

import (
	"lemfi/simplebank/config"
	"net/http"

	"github.com/gin-gonic/gin"

	errorResponse "lemfi/simplebank/pkg/errorResponse"
	"lemfi/simplebank/pkg/responseHandler"
)

func (accountController *AccountController) GetAccountsController(c *gin.Context) {
	config.Logger.Info("Fetching all accounts", "method", "GET", "endpoint", "/accounts")

	accounts, err := accountController.accountService.GetAccounts()
	if err != nil {
		config.Logger.Error("Failed to fetch accounts", "error", err.Error())
		errorResponse.ServerErrorResponse(c, err)
		return
	}

	config.Logger.Info("Accounts fetched successfully", "count", len(accounts))

	response := responseHandler.Envelope{
		"accounts": accounts,
	}

	err = responseHandler.WriteJSON(c.Writer, http.StatusOK, response, nil)
	if err != nil {
		config.Logger.Error("Failed to write JSON response", "error", err.Error())
		errorResponse.ServerErrorResponse(c, err)
		return
	}

	config.Logger.Info("Get accounts request completed successfully", "count", len(accounts))
}
