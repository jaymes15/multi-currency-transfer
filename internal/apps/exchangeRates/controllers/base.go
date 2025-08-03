package exchangeRates

import (
	services "lemfi/simplebank/internal/apps/exchangeRates/services"
)

type ExchangeRateController struct {
	exchangeRateService services.ExchangeRateServiceInterface
}

func NewExchangeRateController(exchangeRateService services.ExchangeRateServiceInterface) *ExchangeRateController {
	return &ExchangeRateController{
		exchangeRateService: exchangeRateService,
	}
}
