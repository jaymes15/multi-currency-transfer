package exchangeRates

import (
	services "lemfi/simplebank/internal/apps/exchangeRates/services"
)

type ExchangeRateController struct {
	exchangeRateService *services.ExchangeRateService
}

func NewExchangeRateController(exchangeRateService *services.ExchangeRateService) *ExchangeRateController {
	return &ExchangeRateController{
		exchangeRateService: exchangeRateService,
	}
}
