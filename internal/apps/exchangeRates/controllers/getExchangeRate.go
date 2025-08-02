package exchangeRates

import (
	"net/http"

	"lemfi/simplebank/config"
	"lemfi/simplebank/internal/apps/core"
	requests "lemfi/simplebank/internal/apps/exchangeRates/requests"
	exchangeRateValidation "lemfi/simplebank/internal/apps/exchangeRates/validationMessages"
	"lemfi/simplebank/pkg/errorResponse"
	"lemfi/simplebank/pkg/requestHandler"
	"lemfi/simplebank/pkg/responseHandler"

	"github.com/gin-gonic/gin"
)

// GetExchangeRateController returns exchange rate for a currency pair with amount calculations
func (exchangeRateController *ExchangeRateController) GetExchangeRateController(c *gin.Context) {
	config.Logger.Info("Getting exchange rate for currency pair", "method", "POST", "endpoint", "/exchange-rates/calculate")

	var req requests.GetExchangeRateRequest

	err := requestHandler.ReadJSONGin(c, &req, exchangeRateValidation.GetExchangeRateValidationMessages)
	if err != nil {
		config.Logger.Error("Failed to read exchange rate request", "error", err.Error())
		errorResponse.BadRequestResponse(c, err)
		return
	}

	config.Logger.Info("Exchange rate request validated successfully",
		"from_currency", req.FromCurrency,
		"to_currency", req.ToCurrency,
		"amount", req.Amount.String(),
	)

	result, err := exchangeRateController.exchangeRateService.GetExchangeRate(c.Request.Context(), req)
	if err != nil {
		config.Logger.Error("Failed to get exchange rate", "error", err.Error(),
			"from_currency", req.FromCurrency,
			"to_currency", req.ToCurrency,
			"amount", req.Amount.String(),
		)
		if clientErr, isClient := core.IsClientError(err); isClient {
			errorResponse.BadRequestResponse(c, clientErr)
		} else {
			errorResponse.ServerErrorResponse(c, err)
		}
		return
	}

	config.Logger.Info("Exchange rate retrieved successfully",
		"from_currency", req.FromCurrency,
		"to_currency", req.ToCurrency,
		"rate", result.ExchangeRate.Rate.String(),
		"amount_to_send", result.AmountToSend.String(),
		"amount_to_receive", result.AmountToReceive.String(),
	)

	response := responseHandler.Envelope{
		"exchange_rate":     result.ExchangeRate,
		"amount_to_send":    result.AmountToSend,
		"amount_to_receive": result.AmountToReceive,
		"can_transact":      result.CanTransact,
		"message":           result.Message,
	}

	err = responseHandler.WriteJSON(c.Writer, http.StatusOK, response, nil)
	if err != nil {
		config.Logger.Error("Failed to write JSON response", "error", err.Error())
		errorResponse.ServerErrorResponse(c, err)
		return
	}

	config.Logger.Info("Exchange rate response written successfully",
		"from_currency", req.FromCurrency,
		"to_currency", req.ToCurrency,
		"rate", result.ExchangeRate.Rate.String(),
	)
}
