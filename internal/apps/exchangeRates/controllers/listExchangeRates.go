package exchangeRates

import (
	"net/http"

	"lemfi/simplebank/config"
	"lemfi/simplebank/internal/apps/core"
	"lemfi/simplebank/pkg/errorResponse"
	"lemfi/simplebank/pkg/responseHandler"

	"github.com/gin-gonic/gin"
)

// ListExchangeRatesController returns all exchange rates
func (exchangeRateController *ExchangeRateController) ListExchangeRatesController(c *gin.Context) {
	config.Logger.Info("Listing all exchange rates", "method", "GET", "endpoint", "/exchange-rates")

	result, err := exchangeRateController.exchangeRateService.ListExchangeRates(c.Request.Context())
	if err != nil {
		config.Logger.Error("Failed to list exchange rates", "error", err.Error())
		if clientErr, isClient := core.IsClientError(err); isClient {
			errorResponse.BadRequestResponse(c, clientErr)
		} else {
			errorResponse.ServerErrorResponse(c, err)
		}
		return
	}

	config.Logger.Info("Exchange rates listed successfully", "total", result.Total)

	response := responseHandler.Envelope{
		"exchange_rates": result,
	}

	err = responseHandler.WriteJSON(c.Writer, http.StatusOK, response, nil)
	if err != nil {
		config.Logger.Error("Failed to write JSON response", "error", err.Error())
		errorResponse.ServerErrorResponse(c, err)
		return
	}

	config.Logger.Info("Exchange rates response written successfully", "total", result.Total)
}
